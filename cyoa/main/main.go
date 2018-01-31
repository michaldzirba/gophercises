package main

import (
	"flag"
	"fmt"
	"github.com/michaldzirba/gophercises/cyoa"
	"net/http"
)

const (
	defaultaddress = ":8080"
)

var (
	address_ = defaultaddress
	log      = fmt.Println
)

func init() {
	addressPtr := flag.String("address", defaultaddress, "address/port on which the server is running")
	flag.Parse()

	address_ = *addressPtr
}

func main() {
	// handler that maps to /story_name/pagetitle
	// resolves a story_name from sories folder
	// this is parsed to a Story, maybe cashed
	// story is passed to a template engine to be randered, and returned to a browser
	log("running server at", address_)
	http.ListenAndServe(address_, cyoa.StoryHandler{})
}
