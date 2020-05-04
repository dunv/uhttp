package uhttp

import (
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/dunv/uhelpers"
	"github.com/dunv/ulog"
)

type LoggingResponseWriter struct {
	http.ResponseWriter
	statusCode       int
	additionalOutput map[string]string
	headerWritten    bool
}

func newLoggingResponseWriter(w http.ResponseWriter) *LoggingResponseWriter {
	return &LoggingResponseWriter{
		w,
		http.StatusOK,
		map[string]string{},
		false,
	}
}

func (lrw *LoggingResponseWriter) AddLogOutput(key, value string) {
	lrw.additionalOutput[key] = value
}

func (lrw *LoggingResponseWriter) WriteHeader(code int) {
	if !lrw.headerWritten {
		lrw.statusCode = code
		lrw.headerWritten = true
		lrw.ResponseWriter.WriteHeader(code)
	}
}

func fmtDuration(d time.Duration) string {
	d = d.Round(time.Millisecond)
	milli := d / time.Millisecond
	return fmt.Sprintf("%dms", milli)
}

// Logging log time, method and path of an HTTP-Request
func addLoggingMiddleware(u *UHTTP) func(next http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {

			lrw := newLoggingResponseWriter(w)
			start := time.Now()

			next.ServeHTTP(lrw, r)

			logLineParams := map[string]string{}
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

			params, err := r.Context().Value(CtxKeyGetParams).(R).Printable()
			if err != nil {
				ulog.Errorf("error when trying to log %s", err)
				u.RenderError(w, r, errors.New("internal server error"))
				return
			}
			for key, value := range params {
				logLineParams[fmt.Sprintf("urlParam-%s", key)] = value
			}

			keys := uhelpers.StringKeysFromMap(logLineParams)
			sort.Strings(keys)

			for _, key := range keys {
				logString = fmt.Sprintf("%s [%s: %s]", logString, key, logLineParams[key])
			}

			if len(lrw.additionalOutput) != 0 {
				for key, value := range lrw.additionalOutput {
					logString = fmt.Sprintf("%s [%s: %s]", logString, key, value)
				}
			}

			u.opts.log.Info(logString)
		}
	}
}

func AddLogOutput(w interface{}, key, value string) error {
	writer, ok := w.(*LoggingResponseWriter)
	if !ok {
		return fmt.Errorf("passed in parameter was not of type LoggingResponseWriter (%T)", w)
	}
	writer.AddLogOutput(key, value)
	return nil
}
