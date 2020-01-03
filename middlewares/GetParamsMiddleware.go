package middlewares

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dunv/uhttp/contextkeys"
	"github.com/dunv/uhttp/helpers"
	"github.com/dunv/uhttp/params"
)

func GetParams(optionalGet params.R, requiredGet params.R) func(next http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			actualRaw := r.URL.Query() // map[string][]string
			actual := map[string]string{}
			for key, values := range actualRaw {
				if len(values) > 0 {
					if len(values) > 1 {
						helpers.RenderError(w, r, fmt.Errorf("param %s has more than one value (given multiple times), this is not supported", key))
						return
					}
					actual[key] = values[0]
				}
			}

			paramMap := params.R{}
			err := params.ValidateParams(requiredGet, actual, paramMap, true)
			if err != nil {
				helpers.RenderError(w, r, fmt.Errorf("%v", err))
				return
			}

			err = params.ValidateParams(optionalGet, actual, paramMap, false)
			if err != nil {
				helpers.RenderError(w, r, fmt.Errorf("%v", err))
			}

			ctx := context.WithValue(r.Context(), contextkeys.CtxKeyGetParams, paramMap)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	}
}
