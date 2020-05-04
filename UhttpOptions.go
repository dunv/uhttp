package uhttp

import (
	"net/http"

	"github.com/dunv/ulog"
)

type UhttpOption interface {
	apply(*uhttpOptions)
}

type uhttpOptions struct {
	cors                    string
	log                     ulog.ULogger
	gzipCompressionLevel    int
	encodingErrorLogLevel   ulog.LogLevel
	parseModelErrorLogLevel ulog.LogLevel
	serveMux                *http.ServeMux
}

type funcUhttpOption struct {
	f func(*uhttpOptions)
}

func (fdo *funcUhttpOption) apply(do *uhttpOptions) {
	fdo.f(do)
}

func newFuncUhttpOption(f func(*uhttpOptions)) *funcUhttpOption {
	return &funcUhttpOption{f: f}
}

func WithCORS(cors string) UhttpOption {
	return newFuncUhttpOption(func(o *uhttpOptions) {
		o.cors = cors
	})
}

func WithLogger(logger ulog.ULogger) UhttpOption {
	return newFuncUhttpOption(func(o *uhttpOptions) {
		o.log = logger
	})
}

func WithGzipCompressionLevel(level int) UhttpOption {
	return newFuncUhttpOption(func(o *uhttpOptions) {
		o.gzipCompressionLevel = level
	})
}

func WithEncodingErrorLogLevel(level ulog.LogLevel) UhttpOption {
	return newFuncUhttpOption(func(o *uhttpOptions) {
		o.encodingErrorLogLevel = level
	})
}

func WithParseModelErrorLogLevel(level ulog.LogLevel) UhttpOption {
	return newFuncUhttpOption(func(o *uhttpOptions) {
		o.parseModelErrorLogLevel = level
	})
}

func WithServeMux(serveMux *http.ServeMux) UhttpOption {
	return newFuncUhttpOption(func(o *uhttpOptions) {
		o.serveMux = serveMux
	})
}
