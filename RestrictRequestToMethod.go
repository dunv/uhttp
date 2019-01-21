package uhttp

import (
	"encoding/json"
	"net/http"
)

// Post Enforces a POST request
func Post(next http.HandlerFunc) http.HandlerFunc {
	return Enforce([]string{"POST"})(next)
}

// OptionsPost Enforces a POST or OPTIONS request
func OptionsPost(next http.HandlerFunc) http.HandlerFunc {
	return Enforce([]string{"POST", "OPTIONS"})(next)
}

// Get Enforces a GET request
func Get(next http.HandlerFunc) http.HandlerFunc {
	return Enforce([]string{"GET"})(next)
}

// Put Enforces a PUT request
func Put(next http.HandlerFunc) http.HandlerFunc {
	return Enforce([]string{"PUT"})(next)
}

func enforce(methods []string, next http.HandlerFunc) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if SliceContainsItem(methods, r.Method) {
			next.ServeHTTP(w, r)
		} else {
			js, _ := json.Marshal(Error{"Method not allowed"})
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write(js)
		}
	}
}

// Enforce generic method for enforcing method
func Enforce(methods []string) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if SliceContainsItem(methods, r.Method) {
				next.ServeHTTP(w, r)
			} else {
				js, _ := json.Marshal(Error{"Method not allowed"})
				w.WriteHeader(http.StatusMethodNotAllowed)
				w.Write(js)
			}
		}
	}
}
