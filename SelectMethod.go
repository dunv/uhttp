package uhttp

import (
	"fmt"
	"net/http"
)

func SelectMethod(u *UHTTP, chain Middleware, handler Handler) http.HandlerFunc {
	return chain(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && handler.GetHandler != nil {
			handler.GetHandler(u)(w, r)
			return
		} else if r.Method == http.MethodPost && handler.PostHandler != nil {
			handler.PostHandler(u)(w, r)
			return
		} else if r.Method == http.MethodDelete && handler.DeleteHandler != nil {
			handler.DeleteHandler(u)(w, r)
			return
		}

		w.WriteHeader(http.StatusMethodNotAllowed)
		u.RenderError(w, r, fmt.Errorf("method not allowed"))
	})
}
