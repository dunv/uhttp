package uhttp

import (
	"context"
	"net/http"
	"time"
)

// Handler configured
type Handler struct {
	Pattern        string
	PostHandler    http.HandlerFunc
	PostModel      interface{}
	GetHandler     http.HandlerFunc
	GetModel       interface{}
	DeleteHandler  http.HandlerFunc
	DeleteModel    interface{}
	RequiredGet    R
	OptionalGet    R
	AddMiddlewares []Middleware
	AddMiddleware  *Middleware
	PreProcess     func(ctx context.Context) error
	Timeout        *time.Duration
	TimeoutMessage *string
}

func (h Handler) WsReady(opts *uhttpOptions) Middleware {
	chain := Chain(
		ParseModelMiddleware(opts, h.PostModel, h.GetModel, h.DeleteModel),
		getParamsMiddleware(h.OptionalGet, h.RequiredGet),
	)

	// Add contexts
	for key, value := range requestContext {
		chain = Chain(chain, WithContextMiddleware(key, value))
	}

	// Add global middlewares
	for key := range additionalMiddlewares {
		chain = Chain(chain, additionalMiddlewares[key])
	}

	// Add handler-specified middlewares
	for key := range h.AddMiddlewares {
		chain = Chain(chain, h.AddMiddlewares[key])
	}

	// Add handler-specified middleware
	if h.AddMiddleware != nil {
		chain = Chain(chain, *h.AddMiddleware)
	}

	// Add preProcess
	return Chain(chain, PreProcessMiddleware(h.PreProcess))
}

func (h Handler) HandlerFunc(opts *uhttpOptions) http.HandlerFunc {
	chain := Chain(
		SetCorsMiddleware(&opts.cors),
		SetJSONResponseMiddleware,
		ParseModelMiddleware(opts, h.PostModel, h.GetModel, h.DeleteModel),
		getParamsMiddleware(h.OptionalGet, h.RequiredGet),
		AddLoggingMiddleware,
	)

	// Add contexts
	for key, value := range requestContext {
		chain = Chain(chain, WithContextMiddleware(key, value))
	}

	// Add global middlewares
	for key := range additionalMiddlewares {
		chain = Chain(chain, additionalMiddlewares[key])
	}

	// Add handler-specified middlewares
	for key := range h.AddMiddlewares {
		chain = Chain(chain, h.AddMiddlewares[key])
	}

	// Add handler-specified middleware
	if h.AddMiddleware != nil {
		chain = Chain(chain, *h.AddMiddleware)
	}

	// Add preProcess
	chain = Chain(chain, PreProcessMiddleware(h.PreProcess))

	// Timeouts
	if h.Timeout != nil {
		msg := "timeout"
		if h.TimeoutMessage != nil {
			msg = *h.TimeoutMessage
		}
		return http.TimeoutHandler(SelectMethod(chain, h), *h.Timeout, msg).ServeHTTP
	}

	return SelectMethod(chain, h)
}
