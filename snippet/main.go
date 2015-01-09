package main

import (
	"code.google.com/p/goplan9/plan9/acme"
	"fmt"
	"os"
	"strconv"
	"log"

)

var Snippets = map[string]string {
	"co": "{%% %s %%}",
	"va": "{{ %s }}",
		
}

// Returns the {% selection %} in a acme window
func main() {

	fmt.Println("Runing...")

	if len(os.Args) != 2 {
		fmt.Println("Wrong params")
		os.Exit(1)
	}
	snippet := "Test"
	command := os.Args[1]
	if val, ok := Snippets[command]; ok == false {
    	fmt.Println("Snippet not found")
    	os.Exit(1)
	} else {
		snippet = val
	}
		
	id, _ := strconv.Atoi(os.Getenv("winid"))
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
	//fmt.Print(a, b)
	
	// get user selection
	var selection string
	if a == b {
		selection = ""
	} else {
		// @TODO Look for a more eficient way of getting selection
		body, err := wfile.ReadAll("body")
		if err != nil {
			log.Fatal(err)
		}
		selection = string(body[a:b])
	}
	
	// @TODO if seleciton is blank, put the cursor on the middle of the snippet.
	_, err = wfile.Write("data", []byte(fmt.Sprintf(snippet, selection)))
	if err != nil {
		log.Fatal(err)
	}
	
	

}
