package uhttp

import (
	"net/http"
	"path/filepath"
	"regexp"
)

// ServeSinglePageApp
func ServeSinglePageApp(path string, mainHTMLFile string) {
	http.HandleFunc("/", SinglePageAppHandler(path, mainHTMLFile))
}

// Exposed handler for testing
func SinglePageAppHandler(path string, mainHTMLFile string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Hint from the docs (when testing for index.html we get "301 Moved Permanently"):
		// https://golang.org/pkg/net/http/#FileServer
		// As a special case, the returned file server redirects any request ending in "/index.html" to the same path, without the final "index.html".

		extension, _ := regexp.MatchString("\\.+[a-zA-Z]+", r.URL.EscapedPath())
		if extension {
			// If the url contains an extension, use file server
			http.FileServer(http.Dir(path)).ServeHTTP(w, r)
		} else {
			// If not -> always serve index
			http.ServeFile(w, r, filepath.Join(path, mainHTMLFile))
		}
	})
}
