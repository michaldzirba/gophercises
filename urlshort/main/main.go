package main

import (
	"flag"
	"fmt"
	urlshort "github.com/michaldzirba/gophercises/urlshort"
	"io/ioutil"
	"net/http"
	"os"
)

var yamlfile_ string

func init() {
	yamlfilePtr := flag.String("yaml", "", "yaml file for the definition ")
	flag.Parse()

	yamlfile_ = *yamlfilePtr
}

func main() {
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	handler := urlshort.MapHandler(pathsToUrls, mux)
	if yamlfile_ != `` {
		var err error
		handler, err = urlshort.YAMLHandler(ReadYamlFile(yamlfile_), handler)
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", handler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}

func checkerror(err error) {
	if err != nil {
		panic(err)
	}
}

func ReadYamlFile(yamlfile string) []byte {
	file, err := os.Open(yamlfile)
	defer file.Close()
	checkerror(err)
	b, err := ioutil.ReadFile(yamlfile)
	checkerror(err)
	return b
}
