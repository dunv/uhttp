package uhttp

import (
	"context"
	"net/http"
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
		parseModelMiddleware(u, h.opts.PostModel, h.opts.GetModel, h.opts.DeleteModel),
		getParamsMiddleware(u, h.opts),
		// Do not add logging here: a WS connection has more states which should be logged separately e.g. in the handler
	)

	// Add original responseWriter
	c = chain(c, withOriginalResponseWriter(u))

	// Add contexts
	for key, value := range u.requestContext {
		c = chain(c, withContextMiddleware(u, key, value))
	}

	// Add global middlewares
	for key := range u.opts.globalMiddlewares {
		c = chain(c, u.opts.globalMiddlewares[key])
	}

	// Add handler-specified middlewares
	for key := range h.opts.Middlewares {
		c = chain(c, h.opts.Middlewares[key])
	}

	// Add preProcess
	return chain(c, preProcessMiddleware(u, h.opts.PreProcess))
}

func (h Handler) HandlerFunc(u *UHTTP) http.HandlerFunc {

	c := chain(
		corsMiddleware(u),
		jsonResponseMiddleware(u),
		parseModelMiddleware(u, h.opts.PostModel, h.opts.GetModel, h.opts.DeleteModel),
		getParamsMiddleware(u, h.opts),
		addLoggingMiddleware(u),
	)

	// Add original responseWriter
	c = chain(c, withOriginalResponseWriter(u))

	// Add contexts
	for key, value := range u.requestContext {
		c = chain(c, withContextMiddleware(u, key, value))
	}

	// Add global middlewares
	for key := range u.opts.globalMiddlewares {
		c = chain(c, u.opts.globalMiddlewares[key])
	}

	// Add handler-specified middlewares
	for key := range h.opts.Middlewares {
		c = chain(c, h.opts.Middlewares[key])
	}

	// Add preProcess
	c = chain(c, preProcessMiddleware(u, h.opts.PreProcess))

	if h.opts.CacheEnable {
		c = chain(c, cacheMiddleware(u, h))
	}

	// Timeouts
	if h.opts.Timeout != 0 {
		return http.TimeoutHandler(c(selectMethodMiddleware(u, h.opts)), h.opts.Timeout, h.opts.TimeoutMessage).ServeHTTP
	}

	return c(selectMethodMiddleware(u, h.opts))

}
