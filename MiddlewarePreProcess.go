package uhttp

import (
	"context"
	"net/http"
)

func PreProcessMiddleware(u *UHTTP, preProcess func(ctx context.Context) error) func(next http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if preProcess != nil {
				err := (preProcess)(r.Context())
				if err != nil {
					u.RenderError(w, r, err)
					return
				}
			}
			next.ServeHTTP(w, r)
		}
	}
}
