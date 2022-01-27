package uhttp

import (
	"fmt"
	"net/http"

	"github.com/dunv/ulog"
)

func AuthBasic(u *UHTTP, expectedUsername string, expectedHashedPasswordSha256 string) Middleware {
	return Middleware(func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			actualUsername, actualPlainPassword, ok := r.BasicAuth()
			actualHashedPassword := fmt.Sprintf("%x", u.sha256.Sum([]byte(actualPlainPassword)))

			if !ok || actualUsername != expectedUsername || actualHashedPassword != expectedHashedPasswordSha256 {
				u.RenderErrorWithStatusCode(w, r, http.StatusUnauthorized, fmt.Errorf("Unauthorized"), u.opts.logHandlerErrors)
				return
			}
			ulog.LogIfError(AddLogOutput(w, "authMethod", "basic"))
			ulog.LogIfError(AddLogOutput(w, "user", actualUsername))
			next.ServeHTTP(w, r)
		}
	})
}
