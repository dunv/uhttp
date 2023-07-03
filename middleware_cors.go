package uhttp

import (
	"net/http"
)

// Set CORS response headers
func corsMiddleware(u *UHTTP) func(next http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if u.opts.cors == "" {
				next.ServeHTTP(w, r)
				return
			}

			if r.Method == "OPTIONS" {
				w.Header().Set("Access-Control-Allow-Origin", u.opts.cors)
				w.Header().Set("Access-Control-Allow-Methods", r.Header.Get("Access-Control-Request-Method"))
				w.Header().Set("Access-Control-Allow-Headers", r.Header.Get("Access-Control-Request-Headers"))
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				w.Header().Set("Access-Control-Allow-Max-Age", "86400")
				if _, err := w.Write([]byte{}); err != nil {
					u.Log().Errorf("%s", err)
				}
				return
			}

			w.Header().Set("Access-Control-Allow-Origin", "*")
			next.ServeHTTP(w, r)
		}
	}
}
