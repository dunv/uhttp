package uhttp

import (
	"encoding/json"
	"net/http"

	"github.com/dunv/ulog"
)

func Render(w http.ResponseWriter, r *http.Request, model interface{}) {
	ulog.LogIfError(json.NewEncoder(w).Encode(model))
}

func RenderWithStatusCode(w http.ResponseWriter, r *http.Request, statusCode int, model interface{}) {
	w.WriteHeader(statusCode)
	ulog.LogIfError(json.NewEncoder(w).Encode(model))
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
	js, _ := json.Marshal(map[string]string{"msg": msg})
	w.WriteHeader(statusCode)
	ulog.LogIfErrorSecondArg(w.Write(js))
	if customLog != nil {
		customLog.Errorf("Msg in %s: %s", r.RequestURI, msg)
	} else {
		ulog.Errorf("Msg in %s: %s", r.RequestURI, msg)
	}
}

func renderErrorWithStatusCode(w http.ResponseWriter, r *http.Request, statusCode int, err error) {
	if err != nil {
		js, _ := json.Marshal(map[string]string{"error": err.Error()})
		w.WriteHeader(statusCode)
		ulog.LogIfErrorSecondArg(w.Write(js))
		if customLog != nil {
			customLog.Errorf("Error in %s: %s", r.RequestURI, err.Error())
		} else {
			ulog.Errorf("Error in %s: %s", r.RequestURI, err.Error())
		}
	} else {
		ulog.Panic("Error to be rendered is nil")
	}
}
