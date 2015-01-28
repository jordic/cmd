package main

// filter command for acme windows
// filters CURRENT window line by line matching (selection or param)
// go install github.com/jordic/cmd/wf/

import (
	"bufio"
	"bytes"
	"code.google.com/p/goplan9/plan9/acme"
	"github.com/jordic/fuzzyfs"
	"log"
	"os"
	"regexp"
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
	var pattern string
	if len(os.Args) == 1 {
		// look for a selection
		pattern, _, _ = win.GetSelection()
		if pattern == "" {
			log.Fatal("Ask for something¿?")
		}
	} else {
		pattern = os.Args[1]
	}

	t, err := win.ReadAll("body")
	if err != nil {
		log.Fatal("Can't read window contents")
	}

	win.Addr(",")
	win.Write("data", nil)

	scanner := bufio.NewScanner(bytes.NewBuffer(t))
	re := regexp.MustCompile(pattern)
	output := []string{}
	for scanner.Scan() {
		line := scanner.Bytes()
		if re.Match(line) {
			//win.Write("body", []byte(fmt.Sprintf("%s\n", line)))
			output = append(output, string(line))
		}
	}
	win.Ctl("dirty")
	win.Write("body", []byte(strings.Join(output, "\n")))

	err = win.Addr("#%d", 0)
	if err != nil {
		log.Fatal(err)
	}
	_ = win.Ctl("dot=addr\n")
	_ = win.Ctl("show")
}
