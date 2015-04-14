package main

// go install github.com/jordic/cmd/FF

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"fmt"
	"sort"

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


type DocSort []fuzzstr.Posting
func (a DocSort) Len() int           { return len(a) }
func (a DocSort) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a DocSort) Less(i, j int) bool { return a[i].Pos < a[j].Pos }

type NullWriter int
func (NullWriter) Write([]byte) (int, error) { return 0, nil }


func main() {

	pwd, _ = os.Getwd()
	var err error
	win, err = acme.New()
	if err != nil {
		log.Fatal(err)
	}
	win.Name(pwd + "/+Search")
	win.Ctl("clean")
	win.Fprintf("tag", "Reload ")
	
	log.SetOutput(new(NullWriter))
	
	PopulateFileList()
	
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
			
			if string(e.Text) == "Reload" {
				win.Addr(",")
				win.Write("data", nil)
				win.Ctl("clean")
				list = make([]string, 1024)
				PopulateFileList()
				continue
			}
			
			result := Index.Query( string(e.Text) )
			
			WriteResults(result)		
			continue		
		}
		win.WriteEvent(e)
	}
	os.Exit(0)

}


func WriteResults(result DocSort) {
	win.Addr(",")
	win.Write("data", nil)
	
	
	sort.Sort(result)
	
	for _, p := range result {
		win.Write("body", []byte(list[p.Doc] + "\n"))
		log.Printf("total files %d %s", p.Pos, list[p.Doc])
	}
	
	win.Ctl("clean")
	
}

func PopulateFileList() {

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

		log.Println(path)
		list = append(list, strings.TrimPrefix(path, pwd+"/"))
		
		return nil

	})

	Index = fuzzstr.NewIndex(list)
	win.Write("body", []byte(fmt.Sprintf("Total Files %d. Exec your query", len(list))))

}