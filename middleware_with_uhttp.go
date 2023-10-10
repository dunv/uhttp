package uhttp

import (
	"context"
	"net/http"
)

// attaches uhttp to the context
func withUHTTP(u *UHTTP) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			httpContext := context.WithValue(r.Context(), CtxKeyUHTTP, u)
			next.ServeHTTP(w, r.WithContext(httpContext))
		}
	}
}

func GetUHTTP(r *http.Request) *UHTTP {
	if val := r.Context().Value(CtxKeyUHTTP); val != nil {
		if uhttp, ok := val.(*UHTTP); ok {
			return uhttp
		}
	}

	panic("Could not get UHTTP from context")
}
