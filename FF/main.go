// This command is a fast file searcher for acme..
// Just hit FF on a folder... Exec querys ( middle click a word )
// To get results.
// Index is stored in RAM.
// +todo Handle file exclusions as an argument.
//	On every folder you can declare a .exclude file, with string patterns for 
//  discarded files/folders.
// + Also implemented a history of querys just to save some keystrokes. ( Hit reload )
// The code is writed quick and dirty... some ugly practices but, works ;)
package main

// go install github.com/jordic/cmd/FF

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	//"time"

	"9fans.net/go/acme"
	"github.com/dgryski/go-fuzzstr"
)

var (
	list     = make([]string, 1024)
	pwd      string
	win      *acme.Win
	Prefixes = []string{".", "env", "cache", "bower", "node", "upload"}
	Suffixes = []string{".pyc", "env", "cache", ".jpg", ".png", ".gif"}
	Excludes = []string{}
	Queries  = []string{}
	Index    *fuzzstr.Index
)

type DocSort []fuzzstr.Posting

func (a DocSort) Len() int      { return len(a) }
func (a DocSort) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a DocSort) Less(i, j int) bool {
	if a[i].Pos == a[j].Pos {
		return len(list[a[i].Doc]) < len(list[a[j].Doc])
	}
	return a[i].Pos < a[j].Pos

}

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

	LoadExcludePatterns()
	PopulateFileList()

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

			q := string(e.Text)
			if false == InSlice(Queries, q) {
				Queries = append(Queries, q)
			}
			result := Index.Query(string(e.Text))
			if len(result) > 50 {
				result = result[:50]
			}
			WriteResults(result)
			continue
		}
		log.Println(e)
		win.WriteEvent(e)
	}

}

func WriteResults(result DocSort) {
	win.Addr(",")
	win.Write("data", nil)

	sort.Sort(result)
	var buff bytes.Buffer
	buff.WriteString("\n\n")

	for _, p := range result {
		buff.Write([]byte(list[p.Doc] + "\n"))
		log.Printf("total files %d %s", p.Pos, list[p.Doc])
	}
	win.Write("body", buff.Bytes())

	win.Ctl("clean")
	_ = win.Addr("#0")
	_ = win.Ctl("dot=addr\n")

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

		for i := range Excludes {
			if strings.HasPrefix(name, Excludes[i]) {
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
	win.Write("body", []byte(fmt.Sprintf("\n\n>>>Total Files %d. Exec your query\n", len(list))))
	if len(Queries) > 0 {
		l := strings.Join(Queries, "\n")
		win.Write("body", []byte(l))
	}

	win.Addr("#0,#0")
	win.Ctl("dot=addr\n")
	win.Ctl("clean")
}

// touch .exclude

func LoadExcludePatterns() {
	if _, err := os.Stat(".exclude"); err == nil {

		file, err := os.Open(".exclude")
		if err != nil {
			fmt.Println("Error loading file")
		}
		sc := bufio.NewScanner(file)
		for sc.Scan() {
			Excludes = append(Excludes, sc.Text())
		}

	}

}

func InSlice(slice []string, s string) bool {
	for _, k := range slice {
		if k == s {
			return true
		}
	}
	return false
}