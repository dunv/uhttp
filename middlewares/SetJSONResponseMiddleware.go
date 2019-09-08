package middlewares 

import (
	"net/http"
)

// SetJSONResponse set response headers
func SetJSONResponse(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	}
}
