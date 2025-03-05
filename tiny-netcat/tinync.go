// A small netcat implementation
// go build tinync.go
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
)


func main() {
	// Setting up flags
	listen := flag.Bool("l", false, "Listen mode")
	RPORT := flag.Int("p",8000, "-p [PORT]")
  RHOST := flag.Arg(0)
  COMMAND := flag.String("e","", "Command to be executed")
	flag.Parse()

	// Usage
	if (len(flag.Args()) < 1 && !(*listen)) || 
     (!(*listen) && *COMMAND == "") {
    fmt.Println("Usage: tinync -l [-p RPORT]")
    fmt.Println("       tinync [-p RPORT] [-e COMMAND] RHOST")
		return
	}


	// Functionality
	if *listen {
    // Default port at listening 
		reader := bufio.NewReader(os.Stdin)
		// Setting up listener
		listener, err := net.Listen("tcp", ":"+strconv.Itoa(*RPORT))
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Listening on port %d...\n", *RPORT)
		defer listener.Close()
		// Handle connections
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		handleConnServer(conn, reader)

	} else {
		// Setting up connection
		address := RHOST + ":" + strconv.Itoa(*RPORT)
		conn, err := net.Dial("tcp", address)
		if err != nil {
			fmt.Println("Error connecting to the host")
      fmt.Println(RHOST)
      fmt.Println(*RPORT)
			return
		}
		// Handle connection
		fmt.Printf("Connected to %s:%d\n", RHOST, *RPORT)
		handleConn(conn,*COMMAND)
	}
	// Error handling
	defer func() {
		if r := recover(); r != nil {
      fmt.Println("Usage: tinync -l [-p RPORT]")
      fmt.Println("       tinync [-p RPORT] [-e COMMAND] RHOST")
		}
	}()
} // main

/* Connection Handler (client) */
func handleConn(conn net.Conn, command string) {
	defer conn.Close()
	// cmd := exec.Command("/bin/sh", "-i")
  var cmd *exec.Cmd
  cparts := strings.Fields(command)
  if len(cparts) == 0 {
   cmd = exec.Command(command)   
  } else {
   cmd = exec.Command(cparts[0],cparts[1:]...)   
  }

	rp, wp := io.Pipe()
	cmd.Stdin = conn
	cmd.Stdout = wp
	cmd.Stderr = wp

	go io.Copy(conn, rp)

	if err := cmd.Run(); err != nil {
		fmt.Println(err)
		return
	}
}
func handleConnServer(conn net.Conn, reader *bufio.Reader) {
    defer conn.Close()
    fmt.Printf("Connected to %s\n", conn.RemoteAddr().String())

    // Read
    go func() {
        buf := make([]byte, 1024)
        for {
            n, err := conn.Read(buf)
            if err != nil {
                if err == io.EOF {
                    fmt.Println("Client disconnected")
                } else {
                    fmt.Println("Error reading from client:", err)
                }
                return
            }
            fmt.Print(string(buf[:n]))
        }
    }()

    // Write
    scanner := bufio.NewScanner(os.Stdin)
    for scanner.Scan() {
        input := scanner.Text()
        _, err := conn.Write([]byte(input + "\n"))
        if err != nil {
            fmt.Println("Error sending data to client:", err)
            return
        }
    }
    fmt.Println("Server input closed, waiting for client to disconnect...")
}
