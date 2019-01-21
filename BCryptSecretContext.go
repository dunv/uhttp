package uhttp

import (
	"context"
	"net/http"
)

// CtxKeyBCryptSecret for retrieving the secret in a handlers
const CtxKeyBCryptSecret = ContextKey("bCryptSecret")

// AddBCryptSecret make it available in handlers
func AddBCryptSecret(bCryptSecret string) func(next http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), CtxKeyBCryptSecret, bCryptSecret)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	}
}
