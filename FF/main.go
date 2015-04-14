package main

// go install github.com/jordic/cmd/FF

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"fmt"

	"9fans.net/go/acme"
	"github.com/dgryski/go-fuzzstr"
	"time"
)

var (
	list     = make([]string, 1024)
	pwd      string
	win      *acme.Win
	Prefixes = []string{".", "env", "cache", "bower", "node"}
	Suffixes = []string{".pyc", "env", "cache"}
	Index    *fuzzstr.Index
)

func main() {

	pwd, _ = os.Getwd()
	var err error
	win, err = acme.New()
	if err != nil {
		log.Fatal(err)
	}
	win.Name(pwd + "/+Search")
	win.Ctl("clean")
	
	
	// Populate file list
	filepath.Walk(pwd, func(path string, info os.FileInfo, err error) error {

		name := strings.ToLower(info.Name())
		for i := range Prefixes {
			if strings.HasPrefix(name, Prefixes[i]) {
				if info.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
		}

		for i := range Suffixes {
			if strings.HasSuffix(name, Suffixes[i]) {
				if info.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
		}

		//log.Println(path)
		list = append(list, strings.TrimPrefix(path, pwd+"/"))
		//win.Write("body", []byte(strings.TrimPrefix(path, pwd+"/")+"\n"))
		return nil

	})
	
	
	Index = fuzzstr.NewIndex(list)

	
	win.Write("body", []byte(fmt.Sprintf("Total Files %d", len(list))))
	
	go events()
	
	for {
		time.Sleep(100 * time.Millisecond)
	}

}

func events() {
	for e := range win.EventChan() {
		switch e.C2 {
		case 'x', 'X': // execute
			if string(e.Text) == "Del" {
				win.Ctl("delete")
			}
			log.Println("Query: ", string(e.Text))
			result := Index.Query( string(e.Text) )
			
			WriteResults(result)		
			continue		
		}
		win.WriteEvent(e)
	}
	os.Exit(0)

}


func WriteResults(result []fuzzstr.Posting) {
	win.Addr(",")
	win.Write("data", nil)
	win.Ctl("clean")
	log.Println("total files", len(result))
	for _, p := range result {
		win.Write("body", []byte(list[p.Doc] + "\n"))
	}
	
}