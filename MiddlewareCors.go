package uhttp

import (
	"net/http"

	"github.com/dunv/ulog"
)

// Set CORS response headers
func corsMiddleware(u *UHTTP) func(next http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if u.opts.cors == "" {
				next.ServeHTTP(w, r)
			} else {
				if r.Method == "OPTIONS" {
					w.Header().Set("Access-Control-Allow-Origin", u.opts.cors)
					w.Header().Set("Access-Control-Allow-Methods", r.Header.Get("Access-Control-Request-Method"))
					w.Header().Set("Access-Control-Allow-Headers", r.Header.Get("Access-Control-Request-Headers"))
					w.Header().Set("Access-Control-Allow-Credentials", "true")
					w.Header().Set("Access-Control-Allow-Max-Age", "86400")
					ulog.LogIfErrorSecondArg(w.Write([]byte{}))
				} else {
					w.Header().Set("Access-Control-Allow-Origin", "*")
					next.ServeHTTP(w, r)
				}
			}
		}
	}
}
