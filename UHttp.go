package uhttp

import (
	"net/http"

	"github.com/dunv/uhttp/logging"
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

func init() {
	// Make expected output (which is only for info, not for debugging) more readable
	ulog.AddSkipFunctions(
		"github.com/dunv/uhttp.RenderError",
		"github.com/dunv/uhttp.RenderErrorWithStatusCode",
		"github.com/dunv/uhttp.renderErrorWithStatusCode",
		"github.com/dunv/uhttp.RenderMessage",
		"github.com/dunv/uhttp.RenderMessageWithStatusCode",
		"github.com/dunv/uhttp.renderMessageWithStatusCode",
	)
	ulog.AddReplaceFunction("github.com/dunv/uhttp/middlewares.AddLogging.func1", "uhttp.Logging")
	ulog.AddReplaceFunction("github.com/dunv/uhttp.Handle", "uhttp.Handle")
}

// Handle configuration
func Handle(pattern string, handler models.Handler) {
	chain := Chain(
		middlewares.SetCors(*config.DisableCORS),
		middlewares.SetJSONResponse,
		middlewares.ParseModel(handler),
		middlewares.GetParams(handler),
	)

	// Add contexts
	for key, value := range requestContext {
		chain = Chain(chain, middlewares.WithContext(key, value))
	}

	// Add global middlewares
	for key := range additionalMiddlewares {
		chain = Chain(chain, additionalMiddlewares[key])
	}

	// Add handler-specified middlewares
	for key := range handler.AddMiddlewares {
		chain = Chain(chain, handler.AddMiddlewares[key])
	}

	// Add handler-specified middleware
	if handler.AddMiddleware != nil {
		chain = Chain(chain, *handler.AddMiddleware)
	}

	// Add preProcess
	chain = Chain(chain, middlewares.PreProcess(handler))

	// Do logging here so we have all contexts available
	chain = Chain(chain, middlewares.AddLogging)

	if handler.GetHandler != nil {
		logging.Logger.Infof("Registered http GET %s", pattern)
	} else if handler.PostHandler != nil {
		logging.Logger.Infof("Registered http POST %s", pattern)
	} else if handler.DeleteHandler != nil {
		logging.Logger.Infof("Registered http DELETE %s", pattern)
	}

	http.Handle(pattern, SelectMethod(chain, handler))
}
