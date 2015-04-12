package main

// go run noend.go
// testing

import (
	"log"
	"time"

)


func main() {
	i := 1
	for {
		log.Printf("Loop %d", i)
		time.Sleep(time.Second)
		i++
		if i == 10 {
			break
		}
	}

}