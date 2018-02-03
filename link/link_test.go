package link

import (
	"log"
	"strings"
	"testing"
	"unicode"
)

func TestProcess(t *testing.T) {
	for s, l := range tests {
		if i, err := Process(strings.NewReader(s)); err == nil {
			if !allequal(i, l) {
				t.Fail()
			}
		} else {
			t.Error(err)
		}
	}
}
func equal(actual, expected Link) bool {
	ok := actual.Link == expected.Link && actual.Text == expected.Text
	if !ok {
		log.Print("actual [", actual, "], expected [", expected, "]")
	}
	return ok
}

func allequal(actual, expected []Link) bool {
	if len(actual) != len(expected) {
		return false
	}

	for i, _ := range expected {
		if !equal(actual[i], expected[i]) {
			return false
		}
	}

	return true
}

func TestString(t *testing.T) {
	l := Link{"a", "b"}
	s := strip(l.String())

	if `{"Link":"a","Text":"b"}` != s {
		t.Fail()
	}
}

func strip(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, str)
}

var (
	tests = map[string][]Link{
		`<html>
  <body>
    <h1>Hello!</h1>
    <a href="/other-page">A link to another page</a>
  </body>
  </html>
  `: []Link{Link{"/other-page", "A link to another page"}},

		`<html>
	<head>
	  <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/font-awesome/4.7.0/css/font-awesome.min.css">
	</head>
	<body>
	  <h1>Social stuffs</h1>
	  <div>
	    <a href="https://www.twitter.com/joncalhoun">
	      Check me out on twitter
	      <i class="fa fa-twitter" aria-hidden="true"></i>
	    </a>
	    <a href="https://github.com/gophercises">
	      Gophercises is on <strong>Github</strong>!
	    </a>
	  </div>
	</body>
	</html>`: []Link{Link{"https://www.twitter.com/joncalhoun", "Check me out on twitter"}, Link{"https://github.com/gophercises", "Gophercises is on Github!"}},
	}
)
