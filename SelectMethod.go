package uhttp

import (
	"fmt"
	"net/http"
)

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
		RenderError(w, r, fmt.Errorf("Method not allowed"))
	})
}
