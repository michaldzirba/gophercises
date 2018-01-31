package cyoa

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	start = "intro"
)

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
	text      string
	PageTitle string // points to Page.Title
}
type Page struct {
	Title   string
	Story   string
	options []Arc
}

type Story interface {
	get(pageTitle string) *Page
}

type GlobalStory struct{} // this is a story that is a list of all stories

type JsonStory struct {
	pages map[string]*Page
}

func (story *JsonStory) fromfile(jsonfile string) (Story, error) {
	bytes, err := ioutil.ReadFile(jsonfile)
	check(err)

	return story.create(bytes)
}
func check(e error) {
	if e != nil {
		panic(e)
	}
}

func (story *JsonStory) create(jsonbytes []byte) (Story, error) {
	// load json here
	err := json.Unmarshal(jsonbytes, story)
	check(err)

	// check if the initial page exists
	_, ok := story.pages[start]
	if !ok {
		return nil, errors.New("Story does not contain initial node")
	}

	return story, nil
}

func (story *GlobalStory) get(pageTitle string) *Page {
	// a story that has links to all the currently defined stories. 'landing page'
	return nil
}

func (story *JsonStory) get(pageTitle string) *Page {
	page, ok := story.pages[pageTitle]
	if !ok {
		page = story.pages[start]
	}

	return page
}

type StoryHandler struct {
}

func (handler StoryHandler) ServeHTTP(rw http.ResponseWriter, rq *http.Request) {
	storyname, pagetitle := split(rq.URL.Path)

	story := getstory(storyname)
	page := (*story).get(pagetitle)

	render(&rw, page)
}

func render(rw *http.ResponseWriter, page *Page) {

}

func getstory(storyname string) *Story {
	return nil
}

func split(path string) (string, string) {
	s := strings.Split(path, "/")
	if len(s) >= 2 {
		return s[1], s[2]
	}
	return "", ""
}
