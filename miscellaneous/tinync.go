// simple netcat implementation with -e option available
package main

import (
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
	port := flag.Int("p", 8000, "Port to listen/connect")
	command := flag.String("e", "", "Command to be executed (server side only)")
	flag.Parse()

	// Usage
	if len(os.Args) < 2 {
		fmt.Println("Usage:")
		fmt.Println("  Server mode: tinync -l -p [PORT] -e [COMMAND]")
		fmt.Println("  Client mode: tinync -p [PORT] HOST")
		return
	}

	if *listen {
		// Server mode
		startServer(*port, *command)
	} else {
		// Client mode
		if flag.NArg() < 1 {
			fmt.Println("Error: RHOST argument required in client mode")
			return
		}
		host := flag.Arg(0)
		startClient(host, *port)
	}
}

func startServer(port int, command string) {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer listener.Close()
	fmt.Printf("Listening on port %d...\n", port)

	conn, err := listener.Accept()
	if err != nil {
		fmt.Println("Error accepting connection:", err)
		return
	}
	defer conn.Close()
	fmt.Printf("Connection from %s\n", conn.RemoteAddr())

	if command != "" {
		// Execute command and handle connection
		handleCommandConnection(conn, command)
	} else {
		handleInteractiveServer(conn)
	}
}

func startClient(host string, port int) {
	address := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	defer conn.Close()
	fmt.Printf("Connected to %s\n", address)

	// Read from connection and print to stdout
	go func() {
		io.Copy(os.Stdout, conn)
		fmt.Println("\nConnection closed by server")
		os.Exit(0)
	}()

	// Read from stdin and send to connection
	io.Copy(conn, os.Stdin)
}

func handleCommandConnection(conn net.Conn, command string) {
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return
	}

	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Stdout = conn
	cmd.Stderr = conn
	cmd.Stdin = conn

	if err := cmd.Run(); err != nil {
		fmt.Println("Error executing command:", err)
		return
	}
}

func handleInteractiveServer(conn net.Conn) {
	// Read from connection and print to stdout
	go func() {
		io.Copy(os.Stdout, conn)
		fmt.Println("\nClient disconnected")
		os.Exit(0)
	}()

	// Read from stdin and send to connection
	io.Copy(conn, os.Stdin)
}
