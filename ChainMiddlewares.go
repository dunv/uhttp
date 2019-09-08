package uhttp

import (
	"net/http"

	"github.com/dunv/uhttp/models"
)

// Chain chain multiple middlewares
// copied from: https://hackernoon.com/simple-http-middleware-with-go-79a4ad62889b
func Chain(mw ...models.Middleware) models.Middleware {
	return func(final http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			last := final
			for i := len(mw) - 1; i >= 0; i-- {
				last = mw[i](last)
			}
			last(w, r)
		}
	}
}
