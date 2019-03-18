package uhttp

import (
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
)

// Middleware define type
type Middleware func(next http.HandlerFunc) http.HandlerFunc

// ContextKey so go does not throw an error
type ContextKey string

// Config vars
var mongoClient *mongo.Client
var disableCors bool
var bCryptSecret string
var authMiddleware *Middleware
var authUserResolver *func(*http.Request) string

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
	Pattern        string
	Handler        http.HandlerFunc
	RequiredParams Params
	OptionalParams Params
	Methods        []string
	DbRequired     bool
	AuthRequired   bool
	AuthMiddleware *Middleware
}

// SetConfig set config for all handlers
func SetConfig(_mongoClient *mongo.Client, _disableCors bool, _bCryptSecret string) {
	mongoClient = _mongoClient
	disableCors = _disableCors
	bCryptSecret = _bCryptSecret
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
	chain := Chain(SetCors(disableCors), AddBCryptSecret(bCryptSecret), SetJSONResponse, Enforce(handler.Methods), WithRequiredParams(handler.RequiredParams), WithOptionalParams(handler.OptionalParams))

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

	if handler.DbRequired {
		chain = Chain(chain, WithDB(mongoClient))
	}

	// Do logging here so we have all contexts available
	chain = Chain(chain, Logging(authUserResolver))

	http.Handle(pattern, chain(handler.Handler))
}
