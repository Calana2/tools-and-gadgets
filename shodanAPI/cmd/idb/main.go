package main

import (
	"log"
  "fmt"
	"os"
  "net"
	"shodanAPI/internetdb"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalln("Usage: idbsearch searchterm")
	}

  domain := os.Args[1]
  IP,err := net.LookupHost(domain) // ipv4
    if err != nil {
        fmt.Println("Error resolving domain:", err)
        return
          }

  query :=  internetdb.New()
  host,err := query.IpLookup(IP[0])
  if err != nil {
   log.Fatalln(err)
  }

// Output

 if(host == nil && err == nil) {
  return
 }

  fmt.Printf("IP: %s",host.Ip)

  fmt.Printf("\n\nOpen ports: ")
  for _,p := range host.Ports {
   fmt.Printf("%d  ",p) 
  }

  fmt.Printf("\n\nHostnames: ")
  for _,h := range host.Hostnames {
   fmt.Printf("%s  ",h) 
  }

  fmt.Printf("\n\nVulnerabilities: ")
  for _,v := range host.Vulnerabilities {
   fmt.Printf("%s  ",v) 
  }
  fmt.Printf("\n\n")
}


