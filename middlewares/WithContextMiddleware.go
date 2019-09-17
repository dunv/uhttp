package middlewares

import (
	"context"
	"net/http"

	"github.com/dunv/uhttp/models"
)

// WithContext attaches any object to the context
func WithContext(key string, value interface{}) models.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			httpContext := context.WithValue(r.Context(), key, value)
			next.ServeHTTP(w, r.WithContext(httpContext))
		}
	}
}
