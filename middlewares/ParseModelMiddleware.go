package middlewares

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"

	"github.com/dunv/uhttp/contextkeys"
	"github.com/dunv/uhttp/helpers"
)

// ParseModel parses and adds a model from a requestbody if wanted
func ParseModel(postModel interface{}, getModel interface{}, deleteModel interface{}) func(next http.HandlerFunc) http.HandlerFunc {
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
				// TODO: maybe find more efficient way of restoring body

				// save body
				var bodyBytes []byte
				if r.Body != nil {
					bodyBytes, _ = ioutil.ReadAll(r.Body)
					r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
				}

				modelInterface := reflectModel.Interface()
				err := json.NewDecoder(r.Body).Decode(modelInterface)
				defer r.Body.Close()
				if err != nil {
					helpers.RenderMessageWithStatusCode(w, r, 400, fmt.Sprintf("Could not decode request body (%s)", err))
					return
				}

				// restore body
				if r.Body != nil {
					r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
				}

				ctx := context.WithValue(r.Context(), contextkeys.CtxKeyPostModel, modelInterface)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			next.ServeHTTP(w, r)
		}
	}
}
