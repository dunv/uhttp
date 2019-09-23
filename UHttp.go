package uhttp

import (
	"net/http"

	"github.com/dunv/uhttp/logging"
	"github.com/dunv/ulog"
)

// TODO: make cors more configurable
// TODO: create uwebsocket lib
// TODO: add filters for logging (i.e. do not log everything, or only user etc)
// TODO: make statistics trackable
// TODO: add license stuff to the repos
// TODO: add readme to repos

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
func Handle(pattern string, handler Handler) {
	handlerFunc := handler.HandlerFunc()
	if handler.GetHandler != nil {
		logging.Logger.Infof("Registered http GET %s", pattern)
	} else if handler.PostHandler != nil {
		logging.Logger.Infof("Registered http POST %s", pattern)
	} else if handler.DeleteHandler != nil {
		logging.Logger.Infof("Registered http DELETE %s", pattern)
	}
	http.Handle(pattern, handlerFunc)
}
