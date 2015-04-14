// Process watcher and reloader for acme
// Inspired in 9fans.net/go/acme/Watch
// adapted to work starting servers.
// Currently works starting a django runsrever, 
// Run env/bin/python
// Or Run inside a go package, for runing it
// @todo remove builded files
package main

import (
	"bufio"
	"flag"

	"io"
	"io/ioutil"
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
var lastcmd *exec.Cmd

var (
	sc = flag.String("sc", "runserver_plus", "Python Mode runserver")
)


const (
	ModeGo = iota
	ModePython
)

var Mode = ModeGo

func main() {
	flag.Parse()
	args = flag.Args()

	if len(args) >= 1 && strings.Contains(args[0], "python") {
		Mode = ModePython
		
			args = append(args, "manage.py")
			args = append(args, *sc)
		
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
	if Mode == ModeGo {
		watcher, err = fsnotify.NewWatcher()
		if err != nil {
			log.Fatal(err)
		}
		defer watcher.Close()
		go fswatch()
	}

	// @todo recursive watch subdirs
	if Mode == ModeGo {
		for _, ff := range loadWatchFileList() {
			watcher.Add(ff)
		}
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
				if lastcmd != nil {
					lastcmd.Process.Kill()
					lastcmd.Wait()
				}
				win.Ctl("delete")
			}
		}
		win.WriteEvent(e)
	}
	os.Exit(0)
}

func runit() {
	//var lastcmd *exec.Cmd
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
			var cmd *exec.Cmd
			if Mode == ModeGo {
				// frist we need to build the exec
				win.Addr(",")
				win.Write("data", nil)
				win.Ctl("clean")
				win.Write("body", []byte("Building exec..."))
				tmpfile := tempfile()
				cm1 := exec.Command("go", "build", "-o", tmpfile)
				if out, err := cm1.CombinedOutput(); err != nil {
					win.Fprintf("body", "%s: %s\n", strings.Join(args, " "), err)
					win.Write("body", out)
					continue
				}

				win.Write("body", []byte("Running exec..."))
				cmd = exec.Command(tmpfile, args...)
				cmd.Dir = pwd

			} else {
				cmd = exec.Command(args[0], args[1:]...)
			}

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

// Django dev server, forks itself when runing.
// After killing, the main process, still the auto_reloader is up,
// and must be killed to handle next reload.
func KillForkedProcess() {
	plist, _ := ps.Processes()
	for _, p := range plist {
		//log.Printf("%s %d %d", p.Executable(), p.Pid(), p.PPid())
		if strings.Contains(p.Executable(), "python") {
			err := syscall.Kill(p.Pid(), syscall.SIGKILL)
			if err != nil {
				log.Print("Error killing pyhon process", err)
			}
		}
	}
}

func loadWatchFileList() []string {
	pwd, _ = os.Getwd()
	f, err := os.Open(pwd)
	defer f.Close()
	if err != nil {
		panic(err)
	}
	files, _ := f.Readdir(-1)
	watchlist := make([]string, 0)

	for _, k := range files {
		name := k.Name()
		if strings.HasSuffix(name, ".go") {
			if !strings.Contains(name, "test") {
				watchlist = append(watchlist, name)
			}
		}
	}
	return watchlist
}

func tempfile() string {
	f, _ := ioutil.TempFile("", "run-")
	f.Close()
	os.Remove(f.Name())
	return f.Name()
}