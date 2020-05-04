package uhttp

import "net/http"

type Middleware func(next http.HandlerFunc) http.HandlerFunc

// Chain chain multiple middlewares
// copied from: https://hackernoon.com/simple-http-middleware-with-go-79a4ad62889b
func Chain(mw ...Middleware) Middleware {
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
