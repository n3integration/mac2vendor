package main

import "fmt"
import m2v "github.com/n3integration/mac2vendor"

func main() {
	mac2vnd, err := m2v.Load(m2v.MAC2VND)
	if err != nil {
		fmt.Println("error: ", err)
	}

	vnd, err := mac2vnd.Lookup("84:38:35:70:aa:52")
	if err != nil {
		fmt.Println("lookup error: ", err)
	} else {
		fmt.Println("found ==> ", vnd)
	}
}
