package urlshort

import (
	out "fmt"
	"github.com/boltdb/bolt"
	"gopkg.in/yaml.v2"
	"net/http"
)

var log = out.Println

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	paths, err := parseyml(yml)
	if err != nil {
		return nil, err
	}

	return MapHandler(paths, fallback), nil
}

func MapHandler(paths map[string]string, fallback http.Handler) http.HandlerFunc {
	get := func(p string) (string, bool) {
		path, ok := paths[p]
		return path, ok
	}

	return http.HandlerFunc(makefunction(get, fallback))
}

func BoltHandler(fallback http.Handler) http.HandlerFunc {
	return http.HandlerFunc(makefunction(getFromBolt, fallback))
}

func getFromBolt(path string) (url string, ok bool) {

	return "", false
}

func makefunction(get func(string) (string, bool), fallback http.Handler) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		path, ok := get(r.URL.Path)
		if ok {
			http.Redirect(w, r, path, http.StatusSeeOther)
		} else {
			fallback.ServeHTTP(w, r)
		}
	}
}
func parseyml(yml []byte) (map[string]string, error) {
	var paths []YamlMapping // mapping to an array without super type
	err := yaml.Unmarshal(yml, &paths)
	if err != nil {
		return nil, err
	}

	pathmap := make(map[string]string)
	for _, v := range paths {
		pathmap[v.Path] = v.URL
	}
	return pathmap, nil
}

type YamlMapping struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
}
