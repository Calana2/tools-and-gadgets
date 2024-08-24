package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func handle(src net.Conn) {
 defer src.Close()
 RHOST := os.Args[1]
 RPORT := os.Args[2]

 dst,err := net.Dial("tcp",RHOST+":"+RPORT)

 defer dst.Close()

 if err != nil {
  log.Fatalln("Unable to connect to our unreachable host")
 }
 

 // io.Copy as goroutine because it can blocking
 // Copy our source's output to the destination
 go func() { 
   if _,err = io.Copy(dst,src); err != nil {
    log.Println(err)
   }
 }()

 // Copy our destination's output back to our source
 if _,err = io.Copy(src,dst); err != nil {
  log.Println(err)
 }

}


func main() {

 if len(os.Args) != 3 {
  fmt.Println("Usage: http-proxy RHOST RPORT")
  return
 }

 listener,err := net.Listen("tcp",":80")
 if err != nil {
  log.Fatalln("Unable to bind to port, execute this program as administrator")
 }

 log.Println("Listening on port 80...")
 
 for {
  conn,err := listener.Accept();
  if err != nil {
   log.Fatalln("Unable to accept connection")
  }
  go handle(conn)
 }
}
