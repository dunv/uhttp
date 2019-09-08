package middlewares

import (
	"context"
	"net/http"

	"github.com/dunv/uhttp/contextkeys"
)

// AddBCryptSecret make it available in handlers
func AddBCryptSecret(bCryptSecret string) func(next http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), contextkeys.CtxKeyBCryptSecret, bCryptSecret)
			next.ServeHTTP(w, r.WithContext(ctx))

		}
	}
}
