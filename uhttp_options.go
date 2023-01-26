package uhttp

import (
	"io"
	"log"
	"net/http"
	"time"

	"github.com/klauspost/compress/flate"
	"github.com/klauspost/compress/gzip"
)

type UhttpOption interface {
	apply(*uhttpOptions)
}

type uhttpOptions struct {
	cors               string
	log                Logger
	logEncodingError   func(template string, args ...interface{})
	logParseModelError func(template string, args ...interface{})
	logHandlerError    func(template string, args ...interface{})

	// Global middlewares
	globalMiddlewares []Middleware

	// If handler panics, send the information to the client
	sendPanicInfoToClient bool
	handleHandlerPanics   []func(r *http.Request, err error)

	// Http-Server options
	address           string
	serveMux          *http.ServeMux
	readTimeout       time.Duration
	readHeaderTimeout time.Duration
	writeTimeout      time.Duration
	idleTimeout       time.Duration

	// Encodings
	enableGzip              bool
	gzipCompressionLevel    int
	enableBrotli            bool
	brotliCompressionLevel  int
	enableDeflate           bool
	deflateCompressionLevel int

	// Logging
	silentStaticFileRegistration bool

	// TLS
	enableTLS      bool
	tlsErrorLogger *log.Logger
	tlsCertPath    *string
	tlsKeyPath     *string

	// Prometheus
	enableMetrics bool
	metricsSocket string
	metricsPath   string

	// Caching
	cacheTTLEnforcerInterval time.Duration

	// Granular logging
	logHandlerCalls                 bool
	logHandlerErrors                bool
	logHandlerRegistrations         bool
	logCacheRuns                    bool
	logCustomMiddlewareRegistration bool
	logStaticFileAccess             bool
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

func WithLogger(logger Logger) UhttpOption {
	return newFuncUhttpOption(func(o *uhttpOptions) {
		o.log = logger
	})
}

func WithGzipCompression(enable bool, level int) UhttpOption {
	return newFuncUhttpOption(func(o *uhttpOptions) {
		o.enableGzip = enable
		if level < gzip.HuffmanOnly {
			level = gzip.HuffmanOnly
		} else if level > gzip.BestCompression {
			level = gzip.BestCompression
		}
		o.gzipCompressionLevel = level
	})
}

func WithBrotliCompression(enable bool, level int) UhttpOption {
	return newFuncUhttpOption(func(o *uhttpOptions) {
		o.enableBrotli = enable
		if level < 0 {
			level = 0
		} else if level > 11 {
			level = 11
		}
		o.brotliCompressionLevel = level
	})
}

func WithDeflateCompression(enable bool, level int) UhttpOption {
	return newFuncUhttpOption(func(o *uhttpOptions) {
		o.enableDeflate = enable
		if level < flate.HuffmanOnly {
			level = flate.HuffmanOnly
		} else if level > flate.BestCompression {
			level = flate.BestCompression
		}
		o.deflateCompressionLevel = level
	})
}

// func WithEncodingErrorLogLevel(level zapcore.Level) UhttpOption {
// 	return newFuncUhttpOption(func(o *uhttpOptions) {
// 		o.logEncodingError = level
// 	})
// }

// func WithParseModelErrorLogLevel(level zapcore.Level) UhttpOption {
// 	return newFuncUhttpOption(func(o *uhttpOptions) {
// 		o.logParseModelError = level
// 	})
// }

// func WithHandlerErrorLogLevel(logErrors bool, level zapcore.Level) UhttpOption {
// 	return newFuncUhttpOption(func(o *uhttpOptions) {
// 		o.logHandlerErrors = logErrors
// 		o.logHandlerError = level
// 	})
// }

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

func WithGlobalMiddlewares(middlewares ...Middleware) UhttpOption {
	return newFuncUhttpOption(func(o *uhttpOptions) {
		o.globalMiddlewares = append(o.globalMiddlewares, middlewares...)
	})
}

func WithTLS(certPath string, keyPath string, tlsErrorLogger *log.Logger) UhttpOption {
	usedLogger := log.New(io.Discard, "", log.Lshortfile)
	if tlsErrorLogger != nil {
		usedLogger = tlsErrorLogger
	}
	return newFuncUhttpOption(func(o *uhttpOptions) {
		o.enableTLS = true
		o.tlsErrorLogger = usedLogger
		o.tlsCertPath = &certPath
		o.tlsKeyPath = &keyPath
	})
}

func WithMetrics(metricsSocket string, metricsPath string) UhttpOption {
	return newFuncUhttpOption(func(o *uhttpOptions) {
		o.enableMetrics = true
		o.metricsSocket = metricsSocket
		o.metricsPath = metricsPath
	})
}

func WithSendPanicInfoToClient(sendPanicInfoToClient bool) UhttpOption {
	return newFuncUhttpOption(func(o *uhttpOptions) {
		o.sendPanicInfoToClient = sendPanicInfoToClient
	})
}

func WithGranularLogging(logHandlerCalls bool, logHandlerRegistrations bool, logStaticFileAccess bool) UhttpOption {
	return newFuncUhttpOption(func(o *uhttpOptions) {
		o.logHandlerCalls = logHandlerCalls
		o.logHandlerRegistrations = logHandlerRegistrations
		o.logStaticFileAccess = logStaticFileAccess
	})
}

// WithSilentStaticFileRegistration disables logging for every static file registration, handy, if dealing with large file trees
func WithSilentStaticFileRegistration(makeFileRegistrationSilent bool) UhttpOption {
	return newFuncUhttpOption(func(o *uhttpOptions) {
		o.silentStaticFileRegistration = makeFileRegistrationSilent
	})
}

func WithCacheTTLEnforcerInterval(i time.Duration) UhttpOption {
	return newFuncUhttpOption(func(o *uhttpOptions) {
		o.cacheTTLEnforcerInterval = i
	})
}

func WithLogCustomMiddlewareRegistration() UhttpOption {
	return newFuncUhttpOption(func(o *uhttpOptions) {
		o.logCustomMiddlewareRegistration = true
	})
}

func WithLogCacheRuns(logCacheRuns bool) UhttpOption {
	return newFuncUhttpOption(func(o *uhttpOptions) {
		o.logCacheRuns = logCacheRuns
	})
}

func WithHandleHandlerPanics(fn func(r *http.Request, err error)) UhttpOption {
	return newFuncUhttpOption(func(o *uhttpOptions) {
		o.handleHandlerPanics = append(o.handleHandlerPanics, fn)
	})
}
