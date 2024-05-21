package uhttp

import (
	"context"
	"net/http"
)

// WithContext attaches the original responseWriter to the context
func withOriginalResponseWriter(_ *UHTTP) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			httpContext := context.WithValue(r.Context(), CtxKeyResponseWriter, w)
			next.ServeHTTP(w, r.WithContext(httpContext))
		}
	}
}
