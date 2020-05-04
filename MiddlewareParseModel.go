package uhttp

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
)

// ParseModel parses and adds a model from a requestbody if wanted
func parseModelMiddleware(u *UHTTP, postModel interface{}, getModel interface{}, deleteModel interface{}) func(next http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			var reflectModel reflect.Value
			doParsing := false
			if r.Method == http.MethodPost && postModel != nil {
				reflectModel = reflect.New(reflect.TypeOf(postModel))
				doParsing = true
			} else if r.Method == http.MethodGet && getModel != nil {
				reflectModel = reflect.New(reflect.TypeOf(getModel))
				doParsing = true
			} else if r.Method == http.MethodDelete && deleteModel != nil {
				reflectModel = reflect.New(reflect.TypeOf(deleteModel))
				doParsing = true
			}

			if doParsing {
				// Save body
				var bodyBytes []byte
				if r.Body != nil {
					var err error
					bodyBytes, err = ioutil.ReadAll(r.Body)
					defer r.Body.Close()
					if err != nil {
						u.renderErrorWithStatusCode(w, r, http.StatusInternalServerError, fmt.Errorf("Could not decode request body (%s)", err), false)
						Logger.LogWithLevelf(u.opts.parseModelErrorLogLevel, "parseModelError [path: %s] Could not decode request body %s", r.RequestURI, err.Error())
						return
					}
					r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
				}

				// Parse body
				modelInterface := reflectModel.Interface()
				err := GzipDecodeRequestBody(r, modelInterface)
				if err != nil {
					u.renderErrorWithStatusCode(w, r, http.StatusBadRequest, fmt.Errorf("Could not decode request body (%s)", err), false)
					Logger.LogWithLevelf(u.opts.parseModelErrorLogLevel, "parseModelError [path: %s] Could not decode request body %s", r.RequestURI, err.Error())
					return
				}

				// Restore body
				if r.Body != nil {
					r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
				}

				ctx := context.WithValue(r.Context(), CtxKeyPostModel, modelInterface)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			next.ServeHTTP(w, r)
		}
	}
}

func ParsedModel(r *http.Request) interface{} {
	parsedModel := r.Context().Value(CtxKeyPostModel)
	if parsedModel != nil {
		return parsedModel
	}
	Logger.Error("Using parsedModel in a request without parsedModel")
	return nil
}
