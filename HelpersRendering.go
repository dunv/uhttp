package uhttp

import (
	"compress/gzip"
	"encoding/json"
	"io"
	"net/http"
)

func (u *UHTTP) Render(w http.ResponseWriter, r *http.Request, model interface{}) {
	u.rawRenderWithStatusCode(w, r, http.StatusOK, model)
}

func (u *UHTTP) RenderWithStatusCode(w http.ResponseWriter, r *http.Request, statusCode int, model interface{}) {
	u.rawRenderWithStatusCode(w, r, statusCode, model)
}

func (u *UHTTP) RenderError(w http.ResponseWriter, r *http.Request, err error) {
	u.renderErrorWithStatusCode(w, r, http.StatusBadRequest, err, true)
}

func (u *UHTTP) RenderErrorWithStatusCode(w http.ResponseWriter, r *http.Request, statusCode int, err error) {
	u.renderErrorWithStatusCode(w, r, statusCode, err, true)
}

func (u *UHTTP) RenderMessage(w http.ResponseWriter, r *http.Request, msg string) {
	u.renderMessageWithStatusCode(w, r, http.StatusOK, msg, true)
}

func (u *UHTTP) RenderMessageWithStatusCode(w http.ResponseWriter, r *http.Request, statusCode int, msg string) {
	u.renderMessageWithStatusCode(w, r, statusCode, msg, true)
}

func (u *UHTTP) renderMessageWithStatusCode(w http.ResponseWriter, r *http.Request, statusCode int, msg string, logOut bool) {
	u.rawRenderWithStatusCode(w, r, statusCode, map[string]string{"msg": msg})
	if logOut {
		u.opts.log.Infof("renderMessage [path: %s] %s", r.RequestURI, msg)
	}
}

func (u *UHTTP) renderErrorWithStatusCode(w http.ResponseWriter, r *http.Request, statusCode int, err error, logOut bool) {
	if err != nil {
		u.rawRenderWithStatusCode(w, r, statusCode, map[string]string{"error": err.Error()})
		if logOut {
			u.opts.log.Errorf("renderError [path: %s] %s", r.RequestURI, err.Error())
		}
	} else {
		u.opts.log.Panic("Error to be rendered is nil")
	}
}

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
				u.opts.log.Panicf("could not initialize gzip writer (%s)", err)
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
}
