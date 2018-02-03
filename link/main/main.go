package main

import (
	"fmt"
	"github.com/michaldzirba/gophercises/link"
	"os"
	"strings"
)

var (
	log = fmt.Println
)

func main() {
	for _, arg := range os.Args[1:] { // first would be the program name

		if !strings.HasPrefix(arg, "-") && isfile(arg) {
			file, err := os.Open(arg)
			check(err)
			print(link.Process(file))
		}
	}
}

func isfile(name string) bool {
	if _, err := os.Stat(name); err != nil {
		return false
	}
	return true
}

func print(links []link.Link, err error) {
	for _, l := range links {
		log(l)
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
