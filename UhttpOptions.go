package uhttp

import (
	"net/http"
	"time"

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

	// Http-Server options
	address           string
	serveMux          *http.ServeMux
	readTimeout       time.Duration
	readHeaderTimeout time.Duration
	writeTimeout      time.Duration
	idleTimeout       time.Duration
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

func WithAddress(address string) UhttpOption {
	return newFuncUhttpOption(func(o *uhttpOptions) {
		o.address = address
	})
}

func WithReadTimeout(readTimeout time.Duration) UhttpOption {
	return newFuncUhttpOption(func(o *uhttpOptions) {
		o.readTimeout = readTimeout
	})
}

func WithReadHeaderTimeout(readHeaderTimeout time.Duration) UhttpOption {
	return newFuncUhttpOption(func(o *uhttpOptions) {
		o.readHeaderTimeout = readHeaderTimeout
	})
}

func WithWriteTimeout(writeTimeout time.Duration) UhttpOption {
	return newFuncUhttpOption(func(o *uhttpOptions) {
		o.writeTimeout = writeTimeout
	})
}

func WithIdleTimeout(idleTimeout time.Duration) UhttpOption {
	return newFuncUhttpOption(func(o *uhttpOptions) {
		o.idleTimeout = idleTimeout
	})
}
