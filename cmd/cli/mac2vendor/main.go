package main

import (
	"flag"
	"fmt"
	"os"

	m2v "github.com/n3integration/mac2vendor"
)

func main() {
	update := flag.Bool("update", false, "whether or not to rebuild the mapping table")
	mac := flag.String("mac", "", "the mac address to resolve")
	quiet := flag.Bool("quiet", false, "quiet mode")
	flag.Parse()

	if *mac == "" && !*update {
		fmt.Println("Usage: ", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	} else if *update {
		err := load(Dat)
		if err != nil {
			fmt.Println("error: ", err)
			os.Exit(1)
		}
	} else {
		vnd, err := m2v.Lookup(*mac)
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
}
