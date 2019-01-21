package uhttp

import (
	"net/http"
	"path/filepath"
	"regexp"
)

// ServeSinglePageApp <-
func ServeSinglePageApp(path string, mainHTMLFile string) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		extension, _ := regexp.MatchString("\\.+[a-zA-Z]+", r.URL.EscapedPath())
		// If the url contains an extension, use file server
		if extension {
			http.FileServer(http.Dir(path)).ServeHTTP(w, r)
		} else {
			http.ServeFile(w, r, filepath.Join(path, mainHTMLFile))
		}
	})
}
