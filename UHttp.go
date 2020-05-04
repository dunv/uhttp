package uhttp

import (
	"fmt"
	"net/http"

	"github.com/dunv/uhelpers"
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
		"github.com/dunv/uhttp.renderErrorWithStatusCode",
		"github.com/dunv/uhttp/helpers.renderErrorWithStatusCode",
		"github.com/dunv/uhttp.renderErrorWithStatusCode",
		"github.com/dunv/uhttp/helpers.renderErrorWithStatusCode",
		"github.com/dunv/uhttp.rawRenderErrorWithStatusCode",
		"github.com/dunv/uhttp/helpers.rawRenderErrorWithStatusCode",
	)
	ulog.AddReplaceFunction("github.com/dunv/uhttp.AddLoggingMiddleware.func1", "uhttp.Logging")
	ulog.AddReplaceFunction("github.com/dunv/uhttp.Handle", "uhttp.Handle")
}

type UHTTP struct {
	opts           *uhttpOptions
	requestContext map[string]interface{}
}

func NewUHTTP(opts ...UhttpOption) *UHTTP {
	mergedOpts := &uhttpOptions{
		cors:                    "*",
		log:                     ulog.NewUlog(),
		gzipCompressionLevel:    4,
		encodingErrorLogLevel:   ulog.LEVEL_ERROR,
		parseModelErrorLogLevel: ulog.LEVEL_ERROR,
		serveMux:                http.NewServeMux(),
	}
	for _, opt := range opts {
		opt.apply(mergedOpts)
	}
	return &UHTTP{
		opts:           mergedOpts,
		requestContext: map[string]interface{}{},
	}
}

func (u *UHTTP) AddContext(key string, value interface{}) error {
	keys := uhelpers.StringKeysFromMap(u.requestContext)
	if !uhelpers.SliceContainsItem(keys, key) {
		u.requestContext[key] = value
		return nil
	}
	return fmt.Errorf("contextKey %s already exists", key)
}

// Handle configuration
func (u *UHTTP) Handle(pattern string, handler Handler) {
	handlerFunc := handler.HandlerFunc(u)
	if handler.GetHandler != nil {
		u.opts.log.Infof("Registered http GET %s", pattern)
	} else if handler.PostHandler != nil {
		u.opts.log.Infof("Registered http POST %s", pattern)
	} else if handler.DeleteHandler != nil {
		u.opts.log.Infof("Registered http DELETE %s", pattern)
	}
	u.opts.serveMux.Handle(pattern, handlerFunc)
}
