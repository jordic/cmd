package main

import (
	"bytes"
	"code.google.com/p/goplan9/plan9/acme"
	"fmt"
	"os"
	"strconv"
	"log"

)

var Snippets = map[string]string {
	"co": "{%% %s %%}",
	"va": "{{ %s }}",
	"t": "<%[1]s></%[1]s>",
	"div": "<div>\n%s\n</div>",
	"p": "<p>%s</p>",
	"span": "<span>%s</span>",
	"/*": "/* %s */",
	"//": "// %s",
	"trans": "{% trans '%s' %}",
	"ha": "['%s']",
	"inc": "{% include '%s' %}",
}

// Small snippet manager for acme... 
// For the moment, snippets are compiled inside the bin
// go install github.com/jordic/cmd/snippet/
// for calling snippets b3 on snippet co, if u have 
// selection, this will be available on the snippet.
// as ex: "snippet t" with cursor selected will replace cursor 
// for <cursor><cursor/>
// If u want to contribute, or share your snippets, just fork it
//

func main() {

	if len(os.Args) != 2 {
		fmt.Println("Wrong params. Usage")
		// Show list of snippets suitable for coping in command ba
		var cmds bytes.Buffer
		for k, _ := range Snippets {
        	cmds.WriteString(fmt.Sprintf("'sn %s', ", k))
    	}
    	fmt.Println(cmds.String())
		os.Exit(2)
	}
	var snippet string
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
	// Read current cursor position
	a, b, _ := wfile.ReadAddr()
	
	//fmt.Print(a, b)
	// get user selection
	var selection string
	if a == b {
		selection = ""
	} else {
		data, err := wfile.ReadAll("data")
		if err != nil {
			log.Fatal(err)
		}
		selection = string(data[0:(b-a)])
		//fmt.Println(string(data[0:(b-a)]))
		
		// restore address after read.
		err = wfile.Addr("#%d,#%d", a, b)
		if err != nil {
			log.Fatal(err)
		}
	}
	
	result := fmt.Sprintf(snippet, selection)
	_, err = wfile.Write("data", []byte(result))
	if err != nil {
		log.Fatal(err)
	}
	// Try to put cursor on middle snippet 
	// if empty selection
	if selection == "" {
		c := a + len(result)/2
		err = wfile.Addr("#%d,#%d", c, c)
		if err != nil {
			log.Fatal(err)
		}
		_ = wfile.Ctl("dot=addr\n")
		
	}
	

}