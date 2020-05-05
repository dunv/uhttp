package uhttp

import (
	"context"
	"fmt"
	"net/http"
)

func NewHandler(opts ...HandlerOption) *handlerOptions {
	mergedOpts := &handlerOptions{}
	for _, opt := range opts {
		opt.apply(mergedOpts)
	}
	return mergedOpts
}

type HandlerFunc func(r *http.Request, returnCode *int) interface{}

type HandlerFuncWithModel func(r *http.Request, model interface{}, returnCode *int) interface{}

type PreProcessFunc func(ctx context.Context) error

func (h handlerOptions) WsReady(u *UHTTP) Middleware {
	c := chain(
		parseModelMiddleware(u, h.PostModel, h.GetModel, h.DeleteModel),
		getParamsMiddleware(u, h.OptionalGet, h.RequiredGet),
	)

	// Add contexts
	for key, value := range u.requestContext {
		c = chain(c, withContextMiddleware(u, key, value))
	}

	// Add global middlewares
	for key := range u.opts.globalMiddlewares {
		c = chain(c, u.opts.globalMiddlewares[key])
	}

	// Add handler-specified middlewares
	for key := range h.Middlewares {
		c = chain(c, h.Middlewares[key])
	}

	// Add preProcess
	return chain(c, preProcessMiddleware(u, h.PreProcess))
}

func (h handlerOptions) HandlerFunc(u *UHTTP) http.HandlerFunc {
	c := chain(
		corsMiddleware(u),
		jsonResponseMiddleware(u),
		parseModelMiddleware(u, h.PostModel, h.GetModel, h.DeleteModel),
		getParamsMiddleware(u, h.OptionalGet, h.RequiredGet),
		addLoggingMiddleware(u),
	)

	// Add contexts
	for key, value := range u.requestContext {
		c = chain(c, withContextMiddleware(u, key, value))
	}

	// Add global middlewares
	for key := range u.opts.globalMiddlewares {
		c = chain(c, u.opts.globalMiddlewares[key])
	}

	// Add handler-specified middlewares
	for key := range h.Middlewares {
		c = chain(c, h.Middlewares[key])
	}

	// Add preProcess
	c = chain(c, preProcessMiddleware(u, h.PreProcess))

	// Timeouts
	if h.Timeout != 0 {
		return http.TimeoutHandler(SelectMethod(u, c, h), h.Timeout, h.TimeoutMessage).ServeHTTP
	}

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()

	return SelectMethod(u, c, h)
}
