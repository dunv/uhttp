package uhttp

import (
	"encoding/json"
	"net/http"
)

// SelectMethod <-
func SelectMethod(chain Middleware, handler Handler) http.HandlerFunc {
	return chain(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && handler.GetHandler != nil {
			handler.GetHandler(w, r)
			return
		} else if r.Method == http.MethodPost && handler.PostHandler != nil {
			handler.PostHandler(w, r)
			return
		} else if r.Method == http.MethodDelete && handler.DeleteHandler != nil {
			handler.DeleteHandler(w, r)
			return
		}

		w.WriteHeader(http.StatusMethodNotAllowed)
		CheckAndLogError(json.NewEncoder(w).Encode(Error{"Method not allowed"}))
	})
}
