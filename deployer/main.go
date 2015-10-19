package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
)

var (
	command string
	addr    string
)

func init() {
	flag.StringVar(&command, "cmd", "", "command to execute")
	flag.StringVar(&addr, "listen", ":5050", "listen address")
}

func main() {
	flag.Parse()
	if command == "" {
		log.Fatal("command is required")
	}

	mux := http.NewServeMux()
	mux.Handle("/execute", http.HandlerFunc(executeHandler))
	log.Printf("Listening on %s command %s\n", addr, command)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}

}

func executeHandler(w http.ResponseWriter, r *http.Request) {
	// Get expected phid
	phid := r.FormValue("phid")
	if phid == "" {
		log.Println("Error: phid param not defined")
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if _, err := os.Stat(command); os.IsNotExist(err) {
		log.Printf("Error: command nod defined: %s\n", command)
		http.Error(w, "Command doesn't exists", http.StatusBadRequest)
		return
	}

	go func() {
		cmd := exec.Command(command, phid)
		err := cmd.Run()
		if err != nil {
			log.Printf("Command %s %s error %s\n", command, phid, err)
			return
		}
		log.Printf("Command %s %s successfull executed", command, phid)
	}()

	fmt.Fprintf(w, "Processing... %s %s", command, phid)

}
