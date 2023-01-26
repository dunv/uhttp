package uhttp

import (
	"context"
	"net/http"
	"reflect"
	"runtime"
	"strings"
)

func NewHandler(opts ...HandlerOption) Handler {
	mergedOpts := &handlerOptions{}
	withDefaults().apply(mergedOpts)
	for _, opt := range opts {
		opt.apply(mergedOpts)
	}
	return Handler{opts: *mergedOpts}
}

type Handler struct {
	opts handlerOptions
}

type HandlerFunc func(r *http.Request, returnCode *int) interface{}

type HandlerFuncWithModel func(r *http.Request, model interface{}, returnCode *int) interface{}

type PreProcessFunc func(ctx context.Context) error

func (h Handler) WsReady(u *UHTTP) Middleware {
	c := chain(
		parseModelMiddleware(u, h.opts.postModel, h.opts.getModel, h.opts.deleteModel),
		getParamsMiddleware(u, h.opts),
		// Do not add logging here: a WS connection has more states which should be logged separately e.g. in the handler
	)

	// Add original responseWriter
	c = chain(c, withOriginalResponseWriter(u))

	// Add contexts
	for key, value := range u.requestContext {
		c = chain(c, WithContextMiddleware(key, value))
	}

	// Add global middlewares
	for key := range u.opts.globalMiddlewares {
		c = chain(c, u.opts.globalMiddlewares[key])
	}

	// Add handler-specified middlewares
	for key := range h.opts.middlewares {
		c = chain(c, h.opts.middlewares[key])
	}

	// Add preProcess
	return chain(c, preProcessMiddleware(u, h.opts.preProcess))
}

func (h Handler) HandlerFunc(u *UHTTP) http.HandlerFunc {
	return h.handlerFuncExcludeMiddlewareByName(u, nil)
}

func (h Handler) handlerFuncExcludeMiddlewareByName(u *UHTTP, exclude *string) http.HandlerFunc {
	// Outer middlewares
	c := chain(
		corsMiddleware(u),
		jsonResponseMiddleware(u),
		addLoggingMiddleware(u, &h, false),
	)

	// Add original responseWriter
	c = chain(c, withOriginalResponseWriter(u))

	// Add contexts
	for key, value := range u.requestContext {
		c = chain(c, WithContextMiddleware(key, value))
	}

	// Add global middlewares
	for key := range u.opts.globalMiddlewares {
		f := u.opts.globalMiddlewares[key]
		fName := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
		if u.opts.logCustomMiddlewareRegistration {
			u.opts.log.Infof("Registering custom-middleware for handler %s: %s", h.opts.handlerPattern, fName)
		}

		if exclude != nil {
			if strings.Contains(fName, *exclude) {
				continue
			}
		}
		c = chain(c, u.opts.globalMiddlewares[key])
	}

	// Add parsers
	c = chain(c, parseModelMiddleware(u, h.opts.postModel, h.opts.getModel, h.opts.deleteModel))
	c = chain(c, getParamsMiddleware(u, h.opts))

	// Add handler-specified middlewares
	for key := range h.opts.middlewares {
		f := h.opts.middlewares[key]

		// convenience feature: middleware can be nil (makes it easier to define handlers sometimes)
		if f == nil {
			continue
		}

		fName := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
		if u.opts.logCustomMiddlewareRegistration {
			u.opts.log.Infof("Registering custom-middleware for handler %s: %s", h.opts.handlerPattern, fName)
		}

		if exclude != nil {
			if strings.Contains(fName, *exclude) {
				continue
			}
		}
		c = chain(c, h.opts.middlewares[key])
	}

	// Add preProcess
	c = chain(c, preProcessMiddleware(u, h.opts.preProcess))

	if h.opts.cacheEnable {
		c = chain(c, cacheMiddleware(u, h))
	}

	// Timeouts
	if h.opts.timeout != 0 {
		return http.TimeoutHandler(c(selectMethodMiddleware(u, h.opts)), h.opts.timeout, h.opts.timeoutMessage).ServeHTTP
	}

	return c(selectMethodMiddleware(u, h.opts))

}
