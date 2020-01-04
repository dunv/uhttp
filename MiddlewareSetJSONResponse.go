package uhttp 

import (
	"net/http"
)

// SetJSONResponse set response headers
func SetJSONResponseMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	}
}
