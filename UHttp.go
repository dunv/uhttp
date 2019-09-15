package uhttp

import (
	"net/http"

	"github.com/dunv/uhttp/contextkeys"
	"github.com/dunv/uhttp/middlewares"
	"github.com/dunv/uhttp/models"
	"github.com/dunv/ulog"
)

// TODO: change setters into a single config object (all pointers for nilchecking)
// TODO: make cors more configurable
// TODO: create uwebsocket lib
// TODO: add filters for logging (i.e. do not log everything, or only user etc)
// TODO: make statistics trackable
// TODO: add license stuff to the repos
// TODO: add readme to repos
// TODO: write tests?!
// TODO: move all mongo-specific things into umongo -> ALL libs should not have to rely on mongo

// Config vars
var disableCors bool
var bCryptSecret string
var authMiddleware *models.Middleware
var authUserResolver *func(*http.Request) string
var additionalContext map[contextkeys.ContextKey]interface{}
var customLog ulog.ULogger

// SetConfig set config for all handlers
func SetConfig(_additionalContext map[contextkeys.ContextKey]interface{}, _disableCors bool, _bCryptSecret string, _customLog ulog.ULogger) {
	additionalContext = _additionalContext
	disableCors = _disableCors
	bCryptSecret = _bCryptSecret
	customLog = _customLog

	// Make expected output (which is only for info, not for debugging) more readable
	ulog.AddSkipFunctions(
		"github.com/dunv/uhttp.RenderError",
		"github.com/dunv/uhttp.RenderErrorWithStatusCode",
		"github.com/dunv/uhttp.renderErrorWithStatusCode",
		"github.com/dunv/uhttp.RenderMessage",
		"github.com/dunv/uhttp.RenderMessageWithStatusCode",
		"github.com/dunv/uhttp.renderMessageWithStatusCode",
	)
	ulog.AddReplaceFunction("github.com/dunv/uhttp/middlewares.AddLogging.func1.1", "uhttp.Logging")
	ulog.AddReplaceFunction("github.com/dunv/uhttp.Handle", "uhttp.Handle")
}

// SetAuthMiddleware <-
func SetAuthMiddleware(mw models.Middleware) {
	authMiddleware = &mw
}

// SetAuthUserResolver <-
func SetAuthUserResolver(resolver *func(*http.Request) string) {
	authUserResolver = resolver
}

// Handle configuration
func Handle(pattern string, handler models.Handler) {
	chain := Chain(
		middlewares.SetCors(disableCors),
		middlewares.AddBCryptSecret(bCryptSecret),
		middlewares.SetJSONResponse,
		middlewares.ParseModel(handler),
		middlewares.GetParams(handler),
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

	for key, value := range additionalContext {
		chain = Chain(chain, middlewares.WithContext(key, value))
	}

	chain = Chain(chain, middlewares.PreProcess(handler))

	// Do logging here so we have all contexts available
	chain = Chain(chain, middlewares.AddLogging(authUserResolver))

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
