package middlewares

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dunv/uhttp/models"
	"github.com/dunv/ulog"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func fmtDuration(d time.Duration) string {
	d = d.Round(time.Millisecond)
	milli := d / time.Millisecond
	return fmt.Sprintf("%dms", milli)
}

// Logging log time, method and path of an HTTP-Request
func AddLogging(resolver *func(*http.Request) string) models.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			lrw := newLoggingResponseWriter(w)
			start := time.Now()

			next.ServeHTTP(lrw, r)

			elapsed := time.Since(start)
			var user string
			realIP := r.Header.Get("X-Real-IP") // nginx-proxy adds this header
			if realIP == "" {
				realIP = r.RemoteAddr
			}

			if resolver != nil {
				resolverFunc := *resolver
				user = resolverFunc(r)
			}

			// Do this after "all other middleware went through". That way we can catch the correct statusCode
			ulog.Infof("Uhttp [from: %s] [user: %s] [time: %s] [status: %d] [method: %s] [uri: %s]", realIP, user, fmtDuration(elapsed), lrw.statusCode, r.Method, r.RequestURI)
		}
	}
}
