package uhttp

import (
	"net/http"

	"github.com/dunv/uhttp/helpers"
)

// I haven't found a better way of exposing these methods in the package directly yet...

func Render(w http.ResponseWriter, r *http.Request, model interface{}) {
	helpers.Render(w, r, model)
}

func RenderWithStatusCode(w http.ResponseWriter, r *http.Request, statusCode int, model interface{}) {
	helpers.RenderWithStatusCode(w, r, statusCode, model)
}
func RenderError(w http.ResponseWriter, r *http.Request, err error) {
	helpers.RenderError(w, r, err)
}
func RenderErrorWithStatusCode(w http.ResponseWriter, r *http.Request, statusCode int, err error) {
	helpers.RenderErrorWithStatusCode(w, r, statusCode, err)
}

func RenderMessage(w http.ResponseWriter, r *http.Request, msg string) {
	helpers.RenderMessage(w, r, msg)
}

func RenderMessageWithStatusCode(w http.ResponseWriter, r *http.Request, statusCode int, msg string) {
	helpers.RenderMessageWithStatusCode(w, r, statusCode, msg)
}
