// Go-tcp-port-scanner

package main

import (
	"flag"
	"fmt"
	"net"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

func worker(ports chan int, results chan int, dir string) {
	for p := range ports {
		address := fmt.Sprintf("%s:%d", dir, p)
		conn, err := net.Dial("tcp", address)
		if err != nil {
			results <- 0
			continue
		}
		conn.Close()
		results <- p
	}
}

var portFlag = flag.String("p","nil", "port number, port range")

  func main() {

  primero := 1
  ultimo := 1024

	flag.Parse()

	if flag.NArg() == 0 {
		fmt.Println("Usage: tcpscan [-p PORT] hosts")
		return
	}

  if *portFlag != "nil" {
    // a single port
    match,_ := regexp.MatchString("^\\d+$",*portFlag)
    if match {
     primero,_ = strconv.Atoi(*portFlag) 
     ultimo = primero
    // a port range
    } else if match,_ := regexp.MatchString("^\\d+-\\d+$",*portFlag); match {
     s := strings.Split(*portFlag,"-")
     izq,_ := strconv.Atoi(s[0])
     der,_ := strconv.Atoi(s[1])
     primero, ultimo = izq, der 
     if izq > der { 
      primero, ultimo = ultimo, primero
     }
    // incorrect format
    } else {
     fmt.Println("Incorrect port format, please specify a single port or a range. Example: -p 20-80") 
     return
    }
  }


		for i := 0; i < flag.NArg(); i++ {

			ports := make(chan int, 100)
			results := make(chan int)
			var openPorts []int
			dir := flag.Arg(i)

			for i := 0; i < cap(ports); i++ {
				go worker(ports, results, dir)
			}

			go func() {
				for i := primero; i <= ultimo; i++ {
					ports <- i
				}
			}()

			for i := 0; i < ultimo - primero + 1; i++ {
				port := <-results
				if port != 0 {
					openPorts = append(openPorts, port)
				}
			}

			close(ports)
			close(results)
			sort.Ints(openPorts)
			for _, port := range openPorts {
				fmt.Printf("%s:%d/tcp open\n", dir, port)
			}
  }
}
