package uhttp

import (
	"net/http"
)

func PreProcess(handler Handler) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if handler.PreProcess != nil {
				err := (handler.PreProcess)(r.Context())
				if err != nil {
					RenderError(w, r, err, customLog)
					return
				}
			}
			next.ServeHTTP(w, r)
		}
	}
}
