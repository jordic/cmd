package main

import (
	"code.google.com/p/goplan9/plan9/acme"
	"fmt"
	"os"
	"strconv"
	"log"

)

// Returns the selection in a acme window
func main() {

	id, _ := strconv.Atoi(os.Getenv("winid"))

	//fmt.Println("winid", id)

	wfile, err := acme.Open(id, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Frist read is buggeous
	_, _, _ = wfile.ReadAddr()

	err = wfile.Ctl("addr=dot\n")
	if err != nil {
		log.Fatal(err)
	}

	a, b, _ := wfile.ReadAddr()
	fmt.Print(a, b)
	
	/* If not a selection return ""
	if a == b {
		fmt.Print(string(a))
	}
	
	body, err := wfile.ReadAll("body")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(body[a:b]))*/
	
	

}
