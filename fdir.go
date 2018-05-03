package main

// Small file search for acme, using fuzzysearch.
// go build fdir.go

import (
	//	"bytes"
	"code.google.com/p/goplan9/plan9/acme"
	"fmt"
	"github.com/jordic/fuzzyfs"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var llista *fuzzyfs.DirList

type Win struct {
	*acme.Win
}

func (w *Win) GetSelection() (string, int, int) {
	// First read is buggeous
	_, _, _ = w.ReadAddr()
	err := w.Ctl("addr=dot\n")
	if err != nil {
		log.Fatal(err)
	}
	var selection string
	// Check cursor position
	a, b, _ := w.ReadAddr()
	if a == b {
		selection = ""
	} else {
		data, err := w.ReadAll("data")
		if err != nil {
			log.Fatal(err)
		}
		selection = string(data[0:(b - a)])
		// restore selection
		err = w.Addr("#%d,#%d", a, b)
		if err != nil {
			log.Fatal(err)
		}
	}
	return selection, a, b
}

func main() {
	// current windo id
	id, _ := strconv.Atoi(os.Getenv("winid"))
	aw, err := acme.Open(id, nil)
	if err != nil {
		log.Fatal(err)
	}
	win := Win{
		aw,
	}
	var query string
	if len(os.Args) == 1 {
		// look for a selection
		query, _, _ = win.GetSelection()
		if query == "" {
			log.Fatal("Ask for something¿?")
		}
	} else {
		query = os.Args[1]
	}

	

	//fmt.Printf("Query %s\n", query)
	path, err := filepath.Abs(".")
	path = path + "/"

	llista = fuzzyfs.NewDirList()
	llista.MaxDepth = 5
	llista.PathSelect = fuzzyfs.DirsAndSymlinksAsDirs
	err = llista.Populate(path, nil)
	if err != nil {
		log.Fatal("Can't populate file list")
	}
	res := llista.Query(query, 3)

	var w2 *acme.Win
	t, err := win.ReadAll("tag")
	if strings.Contains(string(t), "+Search") {
		w2 = aw
		// Clear current window
		w2.Addr(",")
		w2.Write("data", nil)
	} else {
		w2, err = acme.New()
	}

	
	w2.Name(path +"+Search")

	for _, r := range res {
		//fmt.Println( r.Path )
		w2.Write("body", []byte("./"+r.Path+"\n"))
	}
	w2.Ctl("clean")
	w2.Ctl("cleartag")
	w2.Ctl("noscroll")
	w2.Fprintf("tag", fmt.Sprintf(" fsearch %s", query))
	err = w2.Addr("#%d", 0)
	if err != nil {
		log.Fatal(err)
	}
	_ = w2.Ctl("dot=addr\n")
	_ = w2.Ctl("show")
}

// AllFiles match all files..
func allFiles(path string, info os.FileInfo) bool {
	if strings.HasPrefix(info.Name(), ".") {
		return false
	}
	if strings.HasSuffix(info.Name(), ".pyc") {
		return false
	}
	return true
}
