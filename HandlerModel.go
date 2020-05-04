package uhttp

import (
	"context"
	"net/http"
	"time"
)

// Handler configured
type Handler struct {
	Pattern        string
	PostHandler    func(u *UHTTP) http.HandlerFunc
	PostModel      interface{}
	GetHandler     func(u *UHTTP) http.HandlerFunc
	GetModel       interface{}
	DeleteHandler  func(u *UHTTP) http.HandlerFunc
	DeleteModel    interface{}
	RequiredGet    R
	OptionalGet    R
	AddMiddlewares []Middleware
	AddMiddleware  *Middleware
	PreProcess     func(ctx context.Context) error
	Timeout        *time.Duration
	TimeoutMessage *string
}

func (h Handler) WsReady(u *UHTTP) Middleware {
	chain := Chain(
		parseModelMiddleware(u, h.PostModel, h.GetModel, h.DeleteModel),
		getParamsMiddleware(u, h.OptionalGet, h.RequiredGet),
	)

	// Add contexts
	for key, value := range requestContext {
		chain = Chain(chain, WithContextMiddleware(u, key, value))
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
	return Chain(chain, PreProcessMiddleware(u, h.PreProcess))
}

func (h Handler) HandlerFunc(u *UHTTP) http.HandlerFunc {
	chain := Chain(
		corsMiddleware(u),
		jsonResponseMiddleware(u),
		parseModelMiddleware(u, h.PostModel, h.GetModel, h.DeleteModel),
		getParamsMiddleware(u, h.OptionalGet, h.RequiredGet),
		addLoggingMiddleware(u),
	)

	// Add contexts
	for key, value := range requestContext {
		chain = Chain(chain, WithContextMiddleware(u, key, value))
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
	chain = Chain(chain, PreProcessMiddleware(u, h.PreProcess))

	// Timeouts
	if h.Timeout != nil {
		msg := "timeout"
		if h.TimeoutMessage != nil {
			msg = *h.TimeoutMessage
		}
		return http.TimeoutHandler(SelectMethod(u, chain, h), *h.Timeout, msg).ServeHTTP
	}

	return SelectMethod(u, chain, h)
}
