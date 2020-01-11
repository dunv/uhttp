package uhttp

import (
	"fmt"
	"net/http"

	"github.com/dunv/ulog"
)

// TODO: make cors more configurable
// TODO: add filters for logging (i.e. do not log everything, or only user etc)
// TODO: make statistics trackable
// TODO: write more tests
// TODO: add a "server-default" with timeouts

func init() {
	// Make expected output (which is only for info, not for debugging) more readable
	ulog.AddSkipFunctions(
		"github.com/dunv/uhttp.RenderError",
		"github.com/dunv/uhttp/helpers.RenderError",
		"github.com/dunv/uhttp.RenderErrorWithStatusCode",
		"github.com/dunv/uhttp/helpers.RenderErrorWithStatusCode",
		"github.com/dunv/uhttp.RenderMessage",
		"github.com/dunv/uhttp/helpers.RenderMessage",
		"github.com/dunv/uhttp.RenderMessageWithStatusCode",
		"github.com/dunv/uhttp/helpers.RenderMessageWithStatusCode",
		"github.com/dunv/uhttp.renderMessageWithStatusCode",
		"github.com/dunv/uhttp/helpers.renderMessageWithStatusCode",
		"github.com/dunv/uhttp/helpers.renderErrorWithStatusCode",
		"github.com/dunv/uhttp.renderErrorWithStatusCode",
	)
	ulog.AddReplaceFunction("github.com/dunv/uhttp.AddLoggingMiddleware.func1", "uhttp.Logging")
	ulog.AddReplaceFunction("github.com/dunv/uhttp.Handle", "uhttp.Handle")
}

// Handle configuration
func Handle(pattern string, handler Handler) {
	handlerFunc := handler.HandlerFunc()
	if handler.GetHandler != nil {
		Logger.Infof("Registered http GET %s", pattern)
	} else if handler.PostHandler != nil {
		Logger.Infof("Registered http POST %s", pattern)
	} else if handler.DeleteHandler != nil {
		Logger.Infof("Registered http DELETE %s", pattern)
	}
	http.Handle(pattern, handlerFunc)
}

func ParsedModel(r *http.Request) interface{} {
	parsedModel := r.Context().Value(CtxKeyPostModel)
	if parsedModel != nil {
		return parsedModel
	}
	Logger.Error("Using parsedModel in a request without parsedModel")
	return nil
}

func AddLogOutput(w interface{}, key, value string) error {
	writer, ok := w.(*LoggingResponseWriter)
	if !ok {
		return fmt.Errorf("passed in parameter was not of type LoggingResponseWriter (%T)", w)
	}
	writer.AddLogOutput(key, value)
	return nil
}
