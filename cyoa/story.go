package cyoa

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
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
	log             = fmt.Println
)

func init() {

	storyproviders_[jsonProvider] = fromfile
}

/**
"intro": {
  "title": "The Little Blue Gopher",
  "story": [
    "Once upon a time, long long ago, there was a little blue gopher. Our little blue friend wanted to go on an adventure, but he wasn't sure where to go. Will you go on an adventure with him?",
    "One of his friends once recommended going to New York to make friends at this mysterious thing called \"GothamGo\". It is supposed to be a big event with free swag and if there is one thing gophers love it is free trinkets. Unfortunately, the gopher once heard a campfire story about some bad fellas named the Sticky Bandits who also live in New York. In the stories these guys would rob toy stores and terrorize young boys, and it sounded pretty scary.",
    "On the other hand, he has always heard great things about Denver. Great ski slopes, a bad hockey team with cheap tickets, and he even heard they have a conference exclusively for gophers like himself. Maybe Denver would be a safer place to visit."
  ],
  "options": [
    {
      "text": "That story about the Sticky Bandits isn't real, it is from Home Alone 2! Let's head to New York.",
      "arc": "new-york"
    },
    {
      "text": "Gee, those bandits sound pretty real to me. Let's play it safe and try our luck in Denver.",
      "arc": "denver"
    }
  ]
}
*/
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
		log("could not find story by name: " + filename)
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

	log(story)

	_, ok := story[start]
	if !ok {
		return nil, errors.New("Story does not contain initial node")
	}

	return &story, nil
}

func (story GlobalStory) get(pageTitle string) *Page {
	// a story that has links to all the currently defined stories. 'landing page'
	return nil
}

func (story JsonStory) get(pageTitle string) *Page {
	page, ok := story[pageTitle]
	if !ok {
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

	f := storyproviders_[handler.Datasource]
	if f != nil {
		story, err := f(storyname)
		if err == nil {
			page := story.get(pagetitle)
			if page != nil {
				render(rw, page)
			}
		}
	}
	handleerror(rw, rq)
}

func render(rw http.ResponseWriter, page *Page) {

	rw.Write([]byte("title: " + page.Title + "\n"))

	for _, s := range page.Story {
		rw.Write([]byte("story: " + s + "\n"))
	}
	for _, a := range page.Options {
		rw.Write([]byte("title: " + a.Arc + "\n"))
		rw.Write([]byte("title: " + a.Text + "\n"))
	}

}

func handleerror(rw http.ResponseWriter, rq *http.Request) {
	rw.Write([]byte("error"))
}

func getstory(storyname string) Story {
	return nil
}

func split(path string) (string, string) {
	s := strings.Split(path, "/")

	log(s, len(s))

	if len(s) >= 3 {
		return s[1], s[2]
	}
	return "", ""
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
