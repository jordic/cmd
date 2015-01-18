package main

// go build pipe_django.go
// cp pipe_django /Users/jordi/bin

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// Filters django manage.py runserver, changing paths to 
// media editable files.. and hiding images..
// launch from acme: 
// ../env/bin/python manage.py runserver -v 0  |[2] pipe_django
// redirecting stderr to stdin
// in bash... ./manage.py runserver 2>&1 | pipe_django

func main() {
	media_dir := os.Getenv("MEDIA_DIR")
	
	//static_dir := os.Getenv("STATIC_DIR")
	scanner := bufio.NewScanner(os.Stdin)
	re := regexp.MustCompile("^.*\"GET")
	re2 := regexp.MustCompile(" HTTP/1.1\".*$")
	is_image := regexp.MustCompile(".*(.png|.jpg|.gif|.jpeg)$")
	
	for scanner.Scan() {
		text:= scanner.Text()
		out := re.ReplaceAllString(text, "")
		out = re2.ReplaceAllString(out, "")
		if(media_dir != "") {
			out = strings.Replace(out, "/media", media_dir, -1)
		}
		if is_image.MatchString(out) != true {		
			fmt.Println(out)
		}
	}
	
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "error:",err)
		os.Exit(1)
	}
	
}