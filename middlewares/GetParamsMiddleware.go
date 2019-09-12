package middlewares

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dunv/uhttp/contextkeys"
	"github.com/dunv/uhttp/helpers"
	"github.com/dunv/uhttp/models"
	"github.com/dunv/uhttp/params"
)

func GetParams(h models.Handler) models.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// Discard values if they occur more than once ("my" design decision here)
			actualRaw := r.URL.Query() // map[string][]string
			actual := map[string]string{}
			for key, values := range actualRaw {
				if len(values) > 0 {
					actual[key] = values[0]
				}
			}

			paramMap := params.R{}
			err := params.ValidateParams(h.RequiredGet, actual, paramMap, true)
			if err != nil {
				helpers.RenderError(w, r, fmt.Errorf("%v", err))
			}

			err = params.ValidateParams(h.OptionalGet, actual, paramMap, false)
			if err != nil {
				helpers.RenderError(w, r, fmt.Errorf("%v", err))
			}

			ctx := context.WithValue(r.Context(), contextkeys.CtxKeyParams, paramMap)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	}
}
