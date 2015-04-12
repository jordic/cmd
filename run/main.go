// Process watcher and reloader for acme
// Inspired in 9fans.net/go/acme/Watch
package main

import (
	"bufio"
	"flag"
	
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"9fans.net/go/acme"
	fsnotify "gopkg.in/fsnotify.v1"
)

// go run main.go go run test/noend.go 





var win *acme.Win
var args []string
var needrun = make(chan int, 1)
var watcher *fsnotify.Watcher
var pwd string

func main() {
	flag.Parse()
	args = flag.Args()

	if len(args) < 1 {
		log.Println("Usage:\nrun command params")
		os.Exit(2)
	}

	var err error
	win, err = acme.New()
	if err != nil {
		log.Fatal(err)
	}
	
	//log.SetOutput(new(NullWriter))
	
	pwd, _ = os.Getwd()
	win.Name(pwd + "/+Run")
	win.Ctl("clean")
	win.Fprintf("tag", "Reload ")

	// Filesystem watcher
	watcher, err = fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()
	// @todo recursive watch subdirs
	watcher.Add(pwd)

	needrun <- 1
	go runit()
	go fswatch()
	go events()

	for {
		time.Sleep(100 * time.Millisecond)
	}

}

func fswatch() {
	for {
		select {
		case _ = <-watcher.Events:
			needrun <- 1
		case err := <-watcher.Errors:
			log.Println("error:", err)
		}
	}
}

func events() {
	for e := range win.EventChan() {
		switch e.C2 {
		case 'x', 'X': // execute
			if string(e.Text) == "Reload" {

				needrun <- 1

			}
			if string(e.Text) == "Del" {
				win.Ctl("delete")
			}
		}
		win.WriteEvent(e)
	}
	os.Exit(0)
}

func runit() {
	var lastcmd *exec.Cmd
	var q = make(chan int, 2)
	//log.Println("needrun out")

	for {
		select {
		case _ = <-needrun:
			//log.Println("needrunit")

			if lastcmd != nil {
				lastcmd.Process.Kill()
				lastcmd.Wait()
				q <- 1
				q <- 1
				//win.Write("body", []byte("killed"))
			}
			lastcmd = nil
			cmd := exec.Command(args[0], args[1:]...)

			stdout, _ := cmd.StdoutPipe()
			stderr, _ := cmd.StderrPipe()
			bfout := bufio.NewReader(stdout)
			bferr := bufio.NewReader(stderr)

			gr := func(buf io.Reader) {
				for {
					b := make([]byte, 8)
					_, err := buf.Read(b)
					if err != nil {
						break
					}
					win.Write("body", b)
					win.Ctl("clean")
					select {
					case _ = <-q:
						//log.Println("Cleaning goroutine")
						break
					default:
					}
				}
			}

			win.Addr(",")
			win.Write("data", nil)
			win.Ctl("clean")

			if err := cmd.Start(); err != nil {
				win.Fprintf("body", "%s: %s\n", strings.Join(args, " "), err)
				return
			}

			lastcmd = cmd
			go gr(bfout)
			go gr(bferr)

			/*func(){
				cmd.Wait()
				win.Write("body", []byte("Process died, Needs reload"))
			}()*/
			
			
		default:
		}
		time.Sleep(100 * time.Millisecond)
	}

}