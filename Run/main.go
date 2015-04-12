// Process watcher and reloader for acme
// Inspired in 9fans.net/go/acme/Watch
// adapted to work with servers that had been started..
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

	"syscall"

	"9fans.net/go/acme"
	"github.com/mitchellh/go-ps"
	fsnotify "gopkg.in/fsnotify.v1"
)

// go run main.go go run test/noend.go
// go install github.com/jordic/cmd/Run

var win *acme.Win
var args []string
var needrun = make(chan int, 1)
var watcher *fsnotify.Watcher
var pwd string

const (
	ModeGo = iota
	ModePython
	ModeUnknown
)

var Mode = ModeUnknown

func main() {
	flag.Parse()
	args = flag.Args()

	if len(args) < 1 {
		loadFileList()
	}

	if len(args)>2 && strings.HasSuffix(args[2], ".go") {
		Mode = ModeGo
	} else if strings.Contains(args[0], "python") {
		Mode = ModePython
		if len(args) == 1 {
			args = append(args, "manage.py")
			args = append(args, "runserver_plus")
		}
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
	if Mode == ModeGo || Mode == ModeUnknown {
		watcher, err = fsnotify.NewWatcher()
		if err != nil {
			log.Fatal(err)
		}
		defer watcher.Close()
		go fswatch()
	}

	// @todo recursive watch subdirs
	if Mode == ModeGo {
		for _, ff := range args[2:] {
			watcher.Add(ff)
		}
	} else if Mode == ModeUnknown {
		watcher.Add(pwd)
	}

	needrun <- 1
	go runit()
	go events()

	for {
		time.Sleep(100 * time.Millisecond)
	}

}

func fswatch() {
	watched := false
	for {
		timer := time.NewTimer(1 * time.Second)
		select {
		case _ = <-watcher.Events:
			watched = true
		case err := <-watcher.Errors:
			log.Println("error: ", err)
		case <-timer.C:
			if watched == true {
				needrun <- 1
				watched = false
			}
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
				if Mode == ModePython {
					KillForkedProcess()
				}
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
	//log.Println("needru out")

	for {
		select {
		case _ = <-needrun:
			//log.Println("neeruni")
			if lastcmd != nil {
				_, err := os.FindProcess(lastcmd.Process.Pid)
				if err != nil {
					log.Print("unabletofindprocess")
				} else {
					err = lastcmd.Process.Kill()
					if Mode == ModePython {
						KillForkedProcess()
					}
					if err == nil {
						lastcmd.Wait()
					}
				}
				q <- 1
				q <- 1
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

		default:
		}
		time.Sleep(100 * time.Millisecond)
	}

}

func KillForkedProcess() {
	plist, _ := ps.Processes()
	for _, p := range plist {
		//log.Println(p.Executable())
		if strings.Contains(p.Executable(), "python") {
			err := syscall.Kill(p.Pid(), syscall.SIGKILL)
			if err != nil {
				log.Print("Error killing pyhon process", err)
			}
		}
	}
}

func loadFileList() {
	pwd, _ = os.Getwd()
	f, err := os.Open(pwd)
	defer f.Close()
	if err != nil {
		panic(err)
	}
	files, _ := f.Readdir(-1)
	args = []string{"go", "run"}
	for _, k := range files {
		name := k.Name()
		if strings.HasSuffix(name, ".go") {
			if !strings.Contains(name, "test") {
				args = append(args, name)
			}
		}
	}

}