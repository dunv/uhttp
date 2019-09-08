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
	"github.com/dunv/uhttp/models"
)

// ParseModel parses and adds a model from a requestbody if wanted
func ParseModel(handler models.Handler) models.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			var reflectModel reflect.Value
			doParsing := false
			if r.Method == http.MethodPost && handler.PostModel != nil {
				reflectModel = reflect.New(reflect.TypeOf(handler.PostModel))
				doParsing = true
			} else if r.Method == http.MethodGet && handler.GetModel != nil {
				reflectModel = reflect.New(reflect.TypeOf(handler.GetModel))
				doParsing = true
			} else if r.Method == http.MethodDelete && handler.DeleteModel != nil {
				reflectModel = reflect.New(reflect.TypeOf(handler.DeleteModel))
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
