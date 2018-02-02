package cyoa

import (
	"encoding/json"
	"errors"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	start        = "intro"
	jsonProvider = "json"
)

var (
	storyproviders_ = make(map[string]func(string) (Story, error))
	pagetemplate_   *template.Template
)

func init() {
	storyproviders_[jsonProvider] = fromfile
	pagetemplate_ = template.Must(template.ParseFiles("templates/html.template"))
}

type Arc struct {
	Text string `json:"text"`
	Arc  string `json:"arc"` // points to Page.Title
}
type Page struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []Arc    `json:"options"`
}

type Story interface {
	get(pageTitle string) *Page
}

type GlobalStory struct{} // this is a story that is a list of all stories

type JsonStory map[string]Page

func getstoryfile(name string) string {
	return "./stories/" + name + ".json"
}

func fromfile(jsonfile string) (Story, error) {
	filename := getstoryfile(jsonfile)
	if !fileexists(filename) {
		log.Println("could not find story by name: " + filename)
		panic("no story")
	}

	bytes, err := ioutil.ReadFile(filename)
	check(err)

	return createJsonStory(bytes)
}

func fileexists(filename string) bool {
	if _, err := os.Stat("./" + filename); err == nil {
		return true
	}
	return false
}

func createJsonStory(jsonbytes []byte) (Story, error) {
	// load json here
	story := JsonStory{}
	err := json.Unmarshal(jsonbytes, &story)
	check(err)

	// check if the initial page exists

	_, ok := story[start]
	if !ok {
		return nil, errors.New("Story does not contain initial node")
	}

	return &story, nil
}

func (story GlobalStory) get(page string) *Page {
	log.Fatal(" no global story, so far")
	// a story that has links to all the currently defined stories. 'landing page'
	return nil
}

func (story JsonStory) get(pageid string) *Page {
	page, ok := story[pageid]
	if !ok {
		log.Fatal("no page by name [", pageid, "] looking for [", start, "]")
		for x, _ := range story {
			log.Println("available: ", x, "]")
		}
		page = story[start]
	}
	return &page
}

type StoryHandler struct {
	Datasource string
}

func (handler StoryHandler) ServeHTTP(rw http.ResponseWriter, rq *http.Request) {
	storyname, pagetitle := split(rq.URL.Path)
	if storyname == "" {
		return
	}

	f, ok := storyproviders_[handler.Datasource]
	if !ok {
		handleerror(rw, errors.New("no provider"))
	} else {
		story, err := f(storyname)
		if err != nil {
			handleerror(rw, err)
		} else {
			page := story.get(pagetitle)
			if page != nil {
				err := pagetemplate_.Execute(rw, page)
				if err != nil {
					handleerror(rw, err)
				}
			}
		}
	}
}

func handleerror(rw http.ResponseWriter, err error) {
	if err != nil {
		log.Fatal(err)
		rw.Write([]byte(err.Error()))
	}
}

func getstory(storyname string) Story {
	return nil
}

func split(path string) (string, string) {
	s := strings.Split(path, "/")

	if len(s) >= 3 {
		return s[1], s[2]
	}
	return "", ""
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}
