package uhttp

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// HelperMethod for rendering a JSON model
func (u *UHTTP) Render(w http.ResponseWriter, r *http.Request, model interface{}) {
	u.rawRenderWithStatusCode(w, r, http.StatusOK, model)
}

// HelperMethod for rendering a JSON model with statusCode in the response
func (u *UHTTP) RenderWithStatusCode(w http.ResponseWriter, r *http.Request, statusCode int, model interface{}) {
	u.rawRenderWithStatusCode(w, r, statusCode, model)
}

// HelperMethod for rendering an error as JSON while automatically setting a 400 statusCode
func (u *UHTTP) RenderError(w http.ResponseWriter, r *http.Request, err error) {
	u.RenderErrorWithStatusCode(w, r, http.StatusBadRequest, err, true)
}

// HelperMethod for rendering an error as JSON with defining a custom statusCode
func (u *UHTTP) RenderErrorWithStatusCode(w http.ResponseWriter, r *http.Request, statusCode int, err error, logOut bool) {
	if err != nil {
		u.rawRenderWithStatusCode(w, r, statusCode, NewHttpErrorResponse(err))
		if logOut {
			u.opts.log.Errorf("renderError [path: %s] %s", r.RequestURI, err.Error())
		}
	} else {
		u.opts.log.Panic("Error to be rendered is nil")
	}
}

// Internal helperMethod with is used for ALL rendering throughout uhttp
// Takes care of encoding responses
func (u *UHTTP) rawRenderWithStatusCode(w http.ResponseWriter, r *http.Request, statusCode int, model interface{}) {
	var writer io.Writer

	// The go-http-client implementation decodes gzip out-of-the-box, but only if it gets 200 OK
	// For now: use the same behavior here
	if statusCode == http.StatusOK {
		switch r.Header.Get("Accept-Encoding") {
		case "gzip":
			w.Header().Add("Content-Encoding", "gzip")
			var err error
			writer, err = gzip.NewWriterLevel(w, u.opts.gzipCompressionLevel)
			if err != nil {
				u.opts.log.Panic(fmt.Errorf("could not initialize gzip writer (%s)", err))
			}
		default:
			writer = w
		}
	} else {
		writer = w
	}

	w.WriteHeader(statusCode)

	err := json.NewEncoder(writer).Encode(model)
	if err != nil {
		// TODO: find a way of doing this per handler!
		u.opts.log.LogWithLevelf(u.opts.encodingErrorLogLevel, "err encoding http response (%s)", err)
		return
	}

	switch typedWriter := writer.(type) {
	case *gzip.Writer:
		err = typedWriter.Close()
		if err != nil {
			// TODO: find a way of doing this per handler!
			u.opts.log.LogWithLevelf(u.opts.encodingErrorLogLevel, "err closing gzip writer (%s)", err)
		}
	}

	switch responseWriter := w.(type) {
	case *cachingResponseWriter:
		err = responseWriter.Close()
		if err != nil {
			// TODO: find a way of doing this per handler!
			u.opts.log.LogWithLevelf(u.opts.encodingErrorLogLevel, "err closing cachingResponseWriter (%s)", err)
		}
	}
}
