package uhttp

import (
	"context"
	"net/http"

	"github.com/dunv/ulog"

	"go.mongodb.org/mongo-driver/mongo"
)

// TODO: change setters into a single config object (all pointers for nilchecking)
// TODO: make cors more configurable
// TODO: create uwebsocket lib
// TODO: move mongo out of this lib -> add hooks for "preEveryRequest"
// TODO: add filters for logging (i.e. do not log everything, or only user etc)
// TODO: make statistics trackable
// TODO: add license stuff to the repos
// TODO: add readme to repos
// TODO: write tests?!
// TODO: move all mongo-specific things into umongo -> ALL libs should not have to rely on mongo

// Middleware define type
type Middleware func(next http.HandlerFunc) http.HandlerFunc

// ContextKey define type
type ContextKey string

// Config vars
var mongoClients map[ContextKey]*mongo.Client
var disableCors bool
var bCryptSecret string
var authMiddleware *Middleware
var authUserResolver *func(*http.Request) string
var additionalContext map[ContextKey]interface{}
var customLog ulog.ULogger

// Chain chain multiple middlewares
// copied from: https://hackernoon.com/simple-http-middleware-with-go-79a4ad62889b
func Chain(mw ...Middleware) Middleware {
	return func(final http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			last := final
			for i := len(mw) - 1; i >= 0; i-- {
				last = mw[i](last)
			}
			last(w, r)
		}
	}
}

// Handler configured
type Handler struct {
	Pattern                   string
	PostHandler               http.HandlerFunc
	PostModel                 interface{}
	GetHandler                http.HandlerFunc
	GetModel                  interface{}
	DeleteHandler             http.HandlerFunc
	DeleteModel               interface{}
	RequiredParams            Params
	OptionalParams            Params
	DbRequired                []ContextKey
	AdditionalContextRequired []ContextKey
	AuthRequired              bool
	AuthMiddleware            *Middleware
	PreProcess                func(ctx context.Context) error
}

// SetConfig set config for all handlers
func SetConfig(_mongoClients map[ContextKey]*mongo.Client, _additionalContext map[ContextKey]interface{}, _disableCors bool, _bCryptSecret string, _customLog ulog.ULogger) {
	mongoClients = _mongoClients
	additionalContext = _additionalContext
	disableCors = _disableCors
	bCryptSecret = _bCryptSecret
	customLog = _customLog

	ulog.AddSkipFunctions(
		"github.com/dunv/uhttp.RenderError",
		"github.com/dunv/uhttp.RenderErrorWithStatusCode",
		"github.com/dunv/uhttp.renderErrorWithStatusCode",
		"github.com/dunv/uhttp.RenderMessage",
		"github.com/dunv/uhttp.RenderMessageWithStatusCode",
		"github.com/dunv/uhttp.renderMessageWithStatusCode",
	)
}

// SetAuthMiddleware <-
func SetAuthMiddleware(mw Middleware) {
	authMiddleware = &mw
}

// SetAuthUserResolver <-
func SetAuthUserResolver(resolver *func(*http.Request) string) {
	authUserResolver = resolver
}

// Handle configuration
func Handle(pattern string, handler Handler) {
	chain := Chain(
		SetCors(disableCors),
		AddBCryptSecret(bCryptSecret),
		SetJSONResponse,
		WithRequiredParams(handler.RequiredParams),
		WithOptionalParams(handler.OptionalParams),
		ParseModel(handler),
	)

	if handler.AuthRequired {
		if handler.AuthMiddleware != nil {
			// Use custom auth for this one handler if one is provided
			chain = Chain(chain, *handler.AuthMiddleware)
		} else {
			// If not custom auth is provided: try to use default and fail if there is no default
			if authMiddleware == nil {
				panic("Tried to use auth without setting auth-middleware first")
			} else {
				chain = Chain(chain, *authMiddleware)
			}
		}
	}

	for _, dbName := range handler.DbRequired {
		chain = Chain(chain, WithDB(dbName, mongoClients[dbName]))
	}

	for _, additionalContextKey := range handler.AdditionalContextRequired {
		if value, ok := additionalContext[additionalContextKey]; ok {
			chain = Chain(chain, WithContext(additionalContextKey, value))
		} else {
			ulog.Panicf("Tried to use context %s without configuring it first", string(additionalContextKey))
		}
	}

	chain = Chain(chain, PreProcess(handler))

	// Do logging here so we have all contexts available
	chain = Chain(chain, Logging(authUserResolver))

	if handler.GetHandler != nil {
		if customLog != nil {
			customLog.Infof("Registered http GET %s", pattern)
		} else {
			ulog.Infof("Registered http GET %s", pattern)
		}
	} else if handler.PostHandler != nil {
		if customLog != nil {
			customLog.Infof("Registered http POST %s", pattern)
		} else {
			ulog.Infof("Registered http POST %s", pattern)
		}

	} else if handler.DeleteHandler != nil {
		if customLog != nil {
			customLog.Infof("Registered http DELETE %s", pattern)
		} else {
			ulog.Infof("Registered http DELETE %s", pattern)
		}

	}

	http.Handle(pattern, SelectMethod(chain, handler))
}
