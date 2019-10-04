package uhttp

import "net/http"

func RealIP(r *http.Request) string {
	realIP := r.Header.Get("X-Real-IP") // nginx-proxy adds this header
	if realIP == "" {
		realIP = r.RemoteAddr
	}
	return realIP
}
