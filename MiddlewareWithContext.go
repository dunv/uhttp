package uhttp

import (
	"context"
	"net/http"
)

// WithContext attaches any object to the context
func withContextMiddleware(u *UHTTP, key ContextKey, value interface{}) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			httpContext := context.WithValue(r.Context(), key, value)
			next.ServeHTTP(w, r.WithContext(httpContext))
		}
	}
}
