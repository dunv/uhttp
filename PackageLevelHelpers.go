package uhttp

import (
	"io"
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

// Returns a reader for an http.Request or http.Response body
// Will regard the "Content-Encoding" header
func ReaderHelper(header http.Header, body io.ReadCloser) (io.Reader, error) {
	return helpers.ReaderHelper(header, body)
}

// Parses a json http.Request body while taking gzip-encoding into account
// Will regard the "Content-Encoding" header
func ParseBody(r *http.Request, model interface{}) error {
	return helpers.ParseBody(r, model)
}
