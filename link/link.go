package link

import (
	"bytes"
	"encoding/json"
	"golang.org/x/net/html"
	"io"
	"log"
	"strings"
)

const (
	empty = ""
	tab   = "  "
)

var (
	print = log.Println
)

type Link struct {
	Link string
	Text string
}

func (link Link) String() string {
	buffer := new(bytes.Buffer)
	encoder := json.NewEncoder(buffer)
	encoder.SetIndent(empty, tab)

	encoder.Encode(link)
	return buffer.String()
}
func Process(reader io.Reader) ([]Link, error) {
	links := make([]Link, 0, 8)
	z := html.NewTokenizer(reader)

	var link Link
	var in bool = false
loop:
	for {
		tt := z.Next()

		if tt == html.ErrorToken {
			break loop
		}
		if tt == html.EndTagToken {
			t := z.Token()
			if t.Data == "a" {
				links = append(links, link)
				in = false
			}
		}

		if tt == html.StartTagToken {
			t := z.Token()
			if t.Data == "a" {
				link = Link{}
				in = true
				for _, a := range t.Attr {
					if a.Key == "href" {
						link.Link = a.Val
						break
					}
				}
			}
		}

		if in && tt == html.TextToken {
			t := z.Token()
			link.Text = link.Text + strings.TrimSpace(t.Data)
		}
	}

	return links, nil

}
