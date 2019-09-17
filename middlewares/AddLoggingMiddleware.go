package middlewares

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/dunv/uhelpers"
	"github.com/dunv/uhttp/helpers"
	"github.com/dunv/uhttp/logging"
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
func AddLogging(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lrw := newLoggingResponseWriter(w)
		start := time.Now()

		next.ServeHTTP(lrw, r)

		logLineParams := helpers.LogLineParams(r)
		logLineParams["duration"] = fmtDuration(time.Since(start))
		realIP := r.Header.Get("X-Real-IP") // nginx-proxy adds this header
		if realIP == "" {
			realIP = r.RemoteAddr
		}
		logLineParams["from"] = realIP
		logLineParams["status"] = strconv.Itoa(lrw.statusCode)
		logLineParams["method"] = r.Method
		logLineParams["uri"] = r.URL.EscapedPath()
		logString := "Uhttp"

		keys := uhelpers.StringKeysFromMap(logLineParams)
		sort.Strings(keys)
		for _, key := range keys {
			logString = fmt.Sprintf("%s [%s: %s]", logString, key, logLineParams[key])
		}
		logging.Logger.Info(logString)
	}
}
