package main

import (
	"flag"
	"github.com/michaldzirba/gophercises/cyoa"
	"log"
	"net/http"
)

const (
	defaultaddress    = ":8080"
	defaultDataSource = "json"
)

var (
	address_    = defaultaddress
	datasource_ = defaultDataSource
)

func init() {
	addressPtr := flag.String("address", defaultaddress, "address/port on which the server is running")
	datasourcePtr := flag.String("datasource", defaultDataSource, "describes the datasource from which the stories will be read")
	flag.Parse()

	address_ = *addressPtr
	datasource_ = *datasourcePtr
}

func main() {
	// handler that maps to /story_name/pagetitle
	// resolves a story_name from sories folder
	// this is parsed to a Story, maybe cashed
	// story is passed to a template engine to be randered, and returned to a browser
	log.Println("running server at", address_)
	log.Fatal(http.ListenAndServe(address_, cyoa.StoryHandler{datasource_}))
}
