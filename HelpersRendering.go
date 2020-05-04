package uhttp

import (
	"compress/gzip"
	"encoding/json"
	"io"
	"net/http"

	"github.com/dunv/ulog"
)

func Render(w http.ResponseWriter, r *http.Request, model interface{}) {
	rawRenderWithStatusCode(w, r, http.StatusOK, model)
}

func RenderWithStatusCode(w http.ResponseWriter, r *http.Request, statusCode int, model interface{}) {
	rawRenderWithStatusCode(w, r, statusCode, model)
}

func RenderError(w http.ResponseWriter, r *http.Request, err error) {
	renderErrorWithStatusCode(w, r, http.StatusBadRequest, err, true)
}

func RenderErrorWithStatusCode(w http.ResponseWriter, r *http.Request, statusCode int, err error) {
	renderErrorWithStatusCode(w, r, statusCode, err, true)
}

func RenderMessage(w http.ResponseWriter, r *http.Request, msg string) {
	renderMessageWithStatusCode(w, r, http.StatusOK, msg, true)
}

func RenderMessageWithStatusCode(w http.ResponseWriter, r *http.Request, statusCode int, msg string) {
	renderMessageWithStatusCode(w, r, statusCode, msg, true)
}

func renderMessageWithStatusCode(w http.ResponseWriter, r *http.Request, statusCode int, msg string, logOut bool) {
	rawRenderWithStatusCode(w, r, statusCode, map[string]string{"msg": msg})
	if logOut {
		Logger.Infof("renderMessage [path: %s] %s", r.RequestURI, msg)
	}
}

func renderErrorWithStatusCode(w http.ResponseWriter, r *http.Request, statusCode int, err error, logOut bool) {
	if err != nil {
		rawRenderWithStatusCode(w, r, statusCode, map[string]string{"error": err.Error()})
		if logOut {
			Logger.Errorf("renderError [path: %s] %s", r.RequestURI, err.Error())
		}
	} else {
		Logger.Panic("Error to be rendered is nil")
	}
}

func rawRenderWithStatusCode(w http.ResponseWriter, r *http.Request, statusCode int, model interface{}) {
	var writer io.Writer

	// The go-http-client implementation decodes gzip out-of-the-box, but only if it gets 200 OK
	// For now: use the same behavior here
	if statusCode == http.StatusOK {
		switch r.Header.Get("Accept-Encoding") {
		case "gzip":
			w.Header().Add("Content-Encoding", "gzip")
			var err error
			// TODO: find a way of doing this per handler!
			writer, err = gzip.NewWriterLevel(w, gzip.DefaultCompression)
			if err != nil {
				Logger.Panicf("could not initialize gzip writer (%s)", err)
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
		Logger.LogWithLevelf(ulog.LEVEL_ERROR, "err encoding http response (%s)", err)
		return
	}

	switch typedWriter := writer.(type) {
	case *gzip.Writer:
		err = typedWriter.Close()
		if err != nil {
			// TODO: find a way of doing this per handler!
			Logger.LogWithLevelf(ulog.LEVEL_ERROR, "err closing gzip writer (%s)", err)
		}
	}
}
