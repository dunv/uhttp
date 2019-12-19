package helpers

import (
	"compress/gzip"
	"encoding/json"
	"io"
	"net/http"

	"github.com/dunv/uhttp/logging"
)

func Render(w http.ResponseWriter, r *http.Request, model interface{}) {
	rawRenderWithStatusCode(w, r, http.StatusOK, model)
}

func RenderWithStatusCode(w http.ResponseWriter, r *http.Request, statusCode int, model interface{}) {
	rawRenderWithStatusCode(w, r, statusCode, model)
}

func RenderError(w http.ResponseWriter, r *http.Request, err error) {
	renderErrorWithStatusCode(w, r, http.StatusBadRequest, err)
}

func RenderErrorWithStatusCode(w http.ResponseWriter, r *http.Request, statusCode int, err error) {
	renderErrorWithStatusCode(w, r, statusCode, err)
}

func RenderMessage(w http.ResponseWriter, r *http.Request, msg string) {
	renderMessageWithStatusCode(w, r, http.StatusOK, msg)
}

func RenderMessageWithStatusCode(w http.ResponseWriter, r *http.Request, statusCode int, msg string) {
	renderMessageWithStatusCode(w, r, statusCode, msg)
}

func renderMessageWithStatusCode(w http.ResponseWriter, r *http.Request, statusCode int, msg string) {
	rawRenderWithStatusCode(w, r, statusCode, map[string]string{"msg": msg})
	logging.Logger.Infof("renderMessage [path: %s] %s", r.RequestURI, msg)
}

func renderErrorWithStatusCode(w http.ResponseWriter, r *http.Request, statusCode int, err error) {
	if err != nil {
		rawRenderWithStatusCode(w, r, statusCode, map[string]string{"error": err.Error()})
		logging.Logger.Errorf("renderError [path: %s] %s", r.RequestURI, err.Error())
	} else {
		logging.Logger.Panic("Error to be rendered is nil")
	}
}

func rawRenderWithStatusCode(w http.ResponseWriter, r *http.Request, statusCode int, model interface{}) {
	var writer io.Writer
	switch r.Header.Get("Accept-Encoding") {
	case "gzip":
		w.Header().Add("Content-Encoding", "gzip")
		var err error
		// TODO: find a way to use config for compressionLevel and enabling and disabling
		writer, err = gzip.NewWriterLevel(w, 5)
		if err != nil {
			logging.Logger.Panicf("could not initialize gzip writer (%s)", err)
		}
	default:
		writer = w
	}

	w.WriteHeader(statusCode)

	err := json.NewEncoder(writer).Encode(model)
	if err != nil {
		logging.Logger.Error("err encoding http response (%s)", err)
		return
	}

	switch typedWriter := writer.(type) {
	case *gzip.Writer:
		err = typedWriter.Close()
		if err != nil {
			logging.Logger.Error("err closing gzip writer (%s)", err)
		}
	}
}
