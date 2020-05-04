package uhttp

import (
	"net/http"
)

// set response headers
func jsonResponseMiddleware(u *UHTTP) func(next http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			next.ServeHTTP(w, r)
		}
	}
}
