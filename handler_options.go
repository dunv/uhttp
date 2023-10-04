package uhttp

import (
	"context"
	"log"
	"time"
)

type HandlerOption interface {
	apply(*handlerOptions)
}

type handlerOptions struct {
	get          HandlerFunc
	getWithModel HandlerFuncWithModel
	getModel     interface{}

	post          HandlerFunc
	postWithModel HandlerFuncWithModel
	postModel     interface{}

	delete          HandlerFunc
	deleteWithModel HandlerFuncWithModel
	deleteModel     interface{}

	requiredGet    R
	optionalGet    R
	middlewares    []Middleware
	preProcess     func(ctx context.Context) error
	timeout        time.Duration
	timeoutMessage string

	cacheEnable                         bool
	cacheFailedRequests                 bool
	cachePersistEncodings               bool
	cacheAutomaticUpdatesInterval       time.Duration
	cacheAutomaticUpdatesSkipMiddleware *string
	cacheAutomaticUpdatesParameters     []map[string]string
	cacheMaxAge                         time.Duration

	debugRawRequestBody func([]byte)

	loggingDisable bool

	// Read-only
	cacheBypassHeader string

	handlerPattern string
}

type funcHandlerOption struct {
	f func(*handlerOptions)
}

func (fdo *funcHandlerOption) apply(do *handlerOptions) {
	fdo.f(do)
}

func newFuncHandlerOption(f func(*handlerOptions)) *funcHandlerOption {
	return &funcHandlerOption{f: f}
}

func withDefaults() HandlerOption {
	return newFuncHandlerOption(func(o *handlerOptions) {
		o.cacheBypassHeader = "X-UHTTP-BYPASS-CACHE"
		o.debugRawRequestBody = func([]byte) {}
	})
}

// Func to be called when the request is invoked with `GET`
func WithGet(h HandlerFunc) HandlerOption {
	return newFuncHandlerOption(func(o *handlerOptions) {
		if o.getWithModel != nil {
			log.Println("ERROR cannot use WithGetModel in conjunction with WithGet. WithGet will supercede this assignment")
		}

		o.get = h
	})
}

// Func to be called when the request is invoked with `GET`
// and a request-body should be parsed into a model
func WithGetModel(m interface{}, h HandlerFuncWithModel) HandlerOption {
	return newFuncHandlerOption(func(o *handlerOptions) {
		if o.get != nil {
			log.Println("ERROR cannot use WithGetModel in conjunction with WithGet. WithGet will supercede this assignment")
		}

		o.getModel = m
		o.getWithModel = h
	})
}

// Func to be called when the request is invoked with `POST`
func WithPost(h HandlerFunc) HandlerOption {
	return newFuncHandlerOption(func(o *handlerOptions) {
		if o.postWithModel != nil {
			log.Println("ERROR cannot use WithPostModel in conjunction with WithPost. WithPost will supercede this assignment")
		}

		o.post = h
	})
}

// Func to be called when the request is invoked with `POST`
// and a request-body should be parsed into a model
func WithPostModel(m interface{}, h HandlerFuncWithModel) HandlerOption {
	return newFuncHandlerOption(func(o *handlerOptions) {
		if o.post != nil {
			log.Println("ERROR cannot use WithPostModel in conjunction with WithPost. WithPost will supercede this assignment")
		}

		o.postModel = m
		o.postWithModel = h
	})
}

// Func to be called when the request is invoked with `DELETE`
func WithDelete(h HandlerFunc) HandlerOption {
	return newFuncHandlerOption(func(o *handlerOptions) {
		if o.deleteWithModel != nil {
			log.Println("ERROR cannot use WithDeleteModel in conjunction with WithDelete. WithDelete will supercede this assignment")
		}

		o.delete = h
	})
}

// Func to be called when the request is invoked with `DELETE`
// and a request-body should be parsed into a model
func WithDeleteModel(m interface{}, h HandlerFuncWithModel) HandlerOption {
	return newFuncHandlerOption(func(o *handlerOptions) {
		if o.delete != nil {
			log.Println("ERROR cannot use WithDeleteModel in conjunction with WithDelete. WithDelete will supercede this assignment")
		}

		o.deleteModel = m
		o.deleteWithModel = h
	})
}

// Add required query-parameters which will be parsed and validated
// The framework will make sure they are present
func WithRequiredGet(r R) HandlerOption {
	return newFuncHandlerOption(func(o *handlerOptions) {
		o.requiredGet = r
	})
}

// Add optional query-parameters which will be parsed and validated
func WithOptionalGet(r R) HandlerOption {
	return newFuncHandlerOption(func(o *handlerOptions) {
		o.optionalGet = r
	})
}

// Add additional middlewares
func WithMiddlewares(m ...Middleware) HandlerOption {
	return newFuncHandlerOption(func(o *handlerOptions) {
		o.middlewares = m
	})
}

// Execute a function before the handler is invoked
func WithPreProcess(p PreProcessFunc) HandlerOption {
	return newFuncHandlerOption(func(o *handlerOptions) {
		o.preProcess = p
	})
}

// Execute the handler with a timeout (wrapped in an original golang `http.TimeoutHandler`)
func WithTimeout(timeout time.Duration, timeoutMessage string) HandlerOption {
	return newFuncHandlerOption(func(o *handlerOptions) {
		o.timeout = timeout
		o.timeoutMessage = timeoutMessage
	})
}

// Cache handler invocations with a maxAge
func WithCache(maxAge time.Duration) HandlerOption {
	return newFuncHandlerOption(func(o *handlerOptions) {
		o.cacheEnable = true
		o.cacheMaxAge = maxAge
	})
}

// Also cache failed requests
func WithCacheFailedRequests() HandlerOption {
	return newFuncHandlerOption(func(o *handlerOptions) {
		o.cacheFailedRequests = true
	})
}

// Call handler in the background discarding the response (only useful if cache is enabled)
func WithAutomaticCacheUpdates(interval time.Duration, skipMiddleware *string, parameters []map[string]string) HandlerOption {
	return newFuncHandlerOption(func(o *handlerOptions) {
		o.cacheAutomaticUpdatesInterval = interval
		o.cacheAutomaticUpdatesSkipMiddleware = skipMiddleware
		if parameters == nil {
			o.cacheAutomaticUpdatesParameters = []map[string]string{{}}
		} else {
			o.cacheAutomaticUpdatesParameters = parameters
		}
	})
}

// When creating the cache, not only keep the response model in the cache
// but also create all enabled compressed versions of it
// this will take load of the server if many calls hit the cache, but comes with a heavy memory penalty
func WithCachePersistEncodings() HandlerOption {
	return newFuncHandlerOption(func(o *handlerOptions) {
		o.cachePersistEncodings = true
	})
}

// Disable access-log for this handler
func WithDisableAccessLogging() HandlerOption {
	return newFuncHandlerOption(func(o *handlerOptions) {
		o.loggingDisable = true
	})
}

// Register callback for raw request-body (for debugging)
func WithDebugRawRequestBody(fn func([]byte)) HandlerOption {
	return newFuncHandlerOption(func(o *handlerOptions) {
		o.debugRawRequestBody = fn
	})
}
