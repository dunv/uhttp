package uhttp

import (
	"context"
	"fmt"
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
	c := chain(
		parseModelMiddleware(u, h.PostModel, h.GetModel, h.DeleteModel),
		getParamsMiddleware(u, h.OptionalGet, h.RequiredGet),
	)

	// Add contexts
	for key, value := range u.requestContext {
		c = chain(c, withContextMiddleware(u, key, value))
	}

	// Add global middlewares
	for key := range additionalMiddlewares {
		c = chain(c, additionalMiddlewares[key])
	}

	// Add handler-specified middlewares
	for key := range h.AddMiddlewares {
		c = chain(c, h.AddMiddlewares[key])
	}

	// Add handler-specified middleware
	if h.AddMiddleware != nil {
		c = chain(c, *h.AddMiddleware)
	}

	// Add preProcess
	return chain(c, preProcessMiddleware(u, h.PreProcess))
}

func (h Handler) HandlerFunc(u *UHTTP) http.HandlerFunc {
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
	for key := range additionalMiddlewares {
		c = chain(c, additionalMiddlewares[key])
	}

	// Add handler-specified middlewares
	for key := range h.AddMiddlewares {
		c = chain(c, h.AddMiddlewares[key])
	}

	// Add handler-specified middleware
	if h.AddMiddleware != nil {
		c = chain(c, *h.AddMiddleware)
	}

	// Add preProcess
	c = chain(c, preProcessMiddleware(u, h.PreProcess))

	// Timeouts
	if h.Timeout != nil {
		msg := "timeout"
		if h.TimeoutMessage != nil {
			msg = *h.TimeoutMessage
		}
		return http.TimeoutHandler(SelectMethod(u, c, h), *h.Timeout, msg).ServeHTTP
	}

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()

	return SelectMethod(u, c, h)
}
