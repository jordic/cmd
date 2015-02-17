package main



import (
	"bytes"
	"code.google.com/p/goplan9/plan9/acme"
	"fmt"
	"log"
	"os"
	"strconv"
)


func main() {
	
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
	// Read current cursor position
	q0, q1, _ := wfile.ReadAddr()

	if q0 == q1 {
		data, err := wfile.ReadAll("data")
		if err != nil {
			log.Fatal(err)
		}
		
	}


}