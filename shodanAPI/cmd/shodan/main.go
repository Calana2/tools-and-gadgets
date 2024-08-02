package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"shodanAPI/shodan"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalln("Usage: shodan searchterm")
	}

	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Error loading .env file")
	}

	apiKey := os.Getenv("SHODAN_API_KEY")
	s := shodan.New(apiKey)

	info, err := s.APIInfo()
	if err != nil {
		log.Panicln(err)
	}

	fmt.Printf("Query Credits: %d\nScan Credits: %d\n", info.QueryCredits, info.ScanCredits)
 fmt.Printf("Plan: %s\nHTTPS: %t\n", info.Plan, info.HTTPS)
	fmt.Printf("Telnet: %t\nUnlocked: %t\n\n", info.Telnet, info.Unlocked)

	hostSearch, err := s.HostSearch(os.Args[1])
	if err != nil {
		log.Panicln(err)
	}

	if len(hostSearch.Matches) == 0 {
		fmt.Println("No matches")
		return
	}

	for _, host := range hostSearch.Matches {
		fmt.Printf("%18s%8d\n", host.IPString, host.Port)
	}

}
