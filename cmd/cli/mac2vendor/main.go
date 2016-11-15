package main

import (
	"flag"
	"fmt"
	"os"

	m2v "github.com/n3integration/mac2vendor"
)

func main() {
	mac := flag.String("mac", "", "the mac address to resolve")
	quiet := flag.Bool("quiet", false, "quiet mode")

	flag.Parse()

	if *mac == "" {
		fmt.Println("Usage: ", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	mac2vnd, err := m2v.Load(m2v.Dat)
	if err != nil {
		fmt.Println("error: ", err)
		os.Exit(1)
	}

	vnd, err := mac2vnd.Lookup(*mac)
	if err != nil {
		fmt.Println("lookup error: ", err)
		os.Exit(1)
	}

	if *quiet {
		fmt.Println(vnd)
	} else {
		fmt.Printf("   MAC: %s\n", *mac)
		fmt.Printf("Vendor: %s\n", vnd)
	}
}
