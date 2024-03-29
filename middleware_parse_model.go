package uhttp

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"reflect"
)

// ParseModel parses and adds a model from a requestbody if wanted
func parseModelMiddleware(u *UHTTP, handlerOpts handlerOptions, postModel interface{}, getModel interface{}, deleteModel interface{}) func(next http.HandlerFunc) http.HandlerFunc {
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
					bodyBytes, err = io.ReadAll(r.Body)
					defer r.Body.Close()
					if err != nil {
						u.RenderErrorWithStatusCode(w, r, http.StatusInternalServerError, fmt.Errorf("Could not decode request body (%s)", err), false)
						u.opts.logParseModelError("parseModelError [path: %s] Could not decode request body %s", r.RequestURI, err.Error())
						// execute callback for rawRequestBody also in case of error
						handlerOpts.debugRawRequestBody(bodyBytes)
						return
					}

					// execute callback for rawRequestBody
					handlerOpts.debugRawRequestBody(bodyBytes)

					// restore body for further use
					r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
				}

				// Parse body
				modelInterface := reflectModel.Interface()
				err := decodeRequestBody(r, modelInterface)
				if err != nil {
					u.RenderErrorWithStatusCode(w, r, http.StatusBadRequest, fmt.Errorf("Could not decode request body (%s)", err), false)
					u.opts.logParseModelError("parseModelError [path: %s] Could not decode request body %s", r.RequestURI, err.Error())
					return
				}

				// Restore body
				if r.Body != nil {
					r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
				}

				ctx := context.WithValue(r.Context(), CtxKeyPostModel, modelInterface)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			next.ServeHTTP(w, r)
		}
	}
}

func parsedModel(r *http.Request) interface{} {
	parsedModel := r.Context().Value(CtxKeyPostModel)
	return parsedModel
}
