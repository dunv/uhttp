package uhttp

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"net/http"
	"path"
	"sort"
	"strconv"
	"time"

	"github.com/dunv/uhelpers"
	"github.com/dunv/ulog"
)

const NO_LOG_MAGIC_URL_FORCE_CACHE = "UHTTP_NO_LOG_FORCE_CACHE"

type LoggingResponseWriter struct {
	u                *UHTTP
	w                http.ResponseWriter
	statusCode       int
	additionalOutput map[string]string
	wroteHeader      bool
}

func newLoggingResponseWriter(w http.ResponseWriter, u *UHTTP) *LoggingResponseWriter {
	return &LoggingResponseWriter{
		u:                u,
		w:                w,
		additionalOutput: map[string]string{},
	}
}

func (lrw *LoggingResponseWriter) AddLogOutput(key, value string) {
	lrw.additionalOutput[key] = value
}

// Delegate Header() to underlying responseWriter
func (lrw *LoggingResponseWriter) Header() http.Header {
	return lrw.w.Header()
}

// Delegate Write() to underlying responseWriter
func (lrw *LoggingResponseWriter) Write(data []byte) (int, error) {
	// the default implementation in net/http/server.go (line 1577 in go 1.17.2) writes the response-header as
	// soon as write is called, if there are no headers written yet
	if !lrw.wroteHeader {
		lrw.WriteHeader(http.StatusOK)
	}
	return lrw.w.Write(data)
}

// Delegate WriteHeader() to underlying responseWriter AND save code
func (lrw *LoggingResponseWriter) WriteHeader(code int) {
	// only write headers once
	if lrw.wroteHeader {
		// copied straight out of the standard-library net/http/server.go
		caller := relevantCaller()
		lrw.u.opts.log.Warnf("superfluous response.WriteHeader call from %s (%s:%d). could happen if the responseWriter is used in a uhttp.Handler AND the function returns something non-nil", caller.Function, path.Base(caller.File), caller.Line)
		return
	}

	lrw.statusCode = code
	lrw.wroteHeader = true
	lrw.w.WriteHeader(code)
}

// Delegate Hijack() to underlying responseWriter
func (lrw *LoggingResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h, ok := lrw.w.(http.Hijacker)
	if !ok {
		return nil, nil, errors.New("hijack not supported")
	}
	return h.Hijack()
}

// Logging log time, method and path of an HTTP-Request
func addLoggingMiddleware(u *UHTTP) func(next http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if r.URL.String() == NO_LOG_MAGIC_URL_FORCE_CACHE {
				next.ServeHTTP(w, r)
				return
			}

			lrw := newLoggingResponseWriter(w, u)
			start := time.Now()

			next.ServeHTTP(lrw, r)

			logLineParams := map[string]string{}
			duration := time.Since(start)
			logLineParams["duration"] = uhelpers.FmtDuration(duration)
			realIP := r.Header.Get("X-Real-IP") // nginx-proxy adds this header
			if realIP == "" {
				realIP = r.RemoteAddr
			}
			logLineParams["from"] = realIP
			logLineParams["status"] = strconv.Itoa(lrw.statusCode)
			logLineParams["method"] = r.Method
			logLineParams["uri"] = r.URL.EscapedPath()
			logString := "Uhttp"

			if paramsRaw, ok := r.Context().Value(CtxKeyGetParams).(R); ok {
				params, err := paramsRaw.Printable()
				if err != nil {
					ulog.Errorf("error when trying to log %s", err)
					u.RenderError(w, r, errors.New("internal server error"))
					return
				}
				for key, value := range params {
					logLineParams[fmt.Sprintf("urlParam-%s", key)] = value
				}
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

			if u.opts.enableMetrics {
				u.opts.log.LogIfError(HandleMetrics(u.metrics, r.Method, lrw.statusCode, r.URL.EscapedPath(), duration))
			}

			if u.opts.logHandlerCalls {
				u.opts.log.Info(logString)
			}
		}
	}
}

func AddLogOutput(w interface{}, key, value string) error {
	writer, ok := w.(*LoggingResponseWriter)
	if !ok {
		// If we cannot add information (this is the case when we are using websockets)
		// just ignore this call
		return nil
	}
	writer.AddLogOutput(key, value)
	return nil
}
