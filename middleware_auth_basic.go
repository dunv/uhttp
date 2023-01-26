package uhttp

import (
	"crypto/sha256"
	"fmt"
	"net/http"
)

func AuthBasic(u *UHTTP, expectedUsername string, expectedHashedPasswordSha256 string) Middleware {
	return Middleware(func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			actualUsername, actualPlainPassword, ok := r.BasicAuth()
			actualHashedPassword := fmt.Sprintf("%x", sha256.Sum256([]byte(actualPlainPassword)))

			if !ok || actualUsername != expectedUsername || actualHashedPassword != expectedHashedPasswordSha256 {
				u.RenderErrorWithStatusCode(w, r, http.StatusUnauthorized, fmt.Errorf("Unauthorized"), u.opts.logHandlerErrors)
				return
			}
			if err := AddLogOutput(w, "authMethod", "basic"); err != nil {
				u.Log().Errorf("%s", err)
			}
			if err := AddLogOutput(w, "user", actualUsername); err != nil {
				u.Log().Errorf("%s", err)
			}
			next.ServeHTTP(w, r)
		}
	})
}
