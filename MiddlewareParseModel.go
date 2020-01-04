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
func ParseModelMiddleware(postModel interface{}, getModel interface{}, deleteModel interface{}) func(next http.HandlerFunc) http.HandlerFunc {
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
						RenderErrorWithStatusCode(w, r, http.StatusInternalServerError, fmt.Errorf("Could not decode request body (%s)", err))
						return
					}
					r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
				}

				// Parse body
				modelInterface := reflectModel.Interface()
				err := ParseBody(r, modelInterface)
				if err != nil {
					RenderErrorWithStatusCode(w, r, http.StatusBadRequest, fmt.Errorf("Could not decode request body (%s)", err))
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