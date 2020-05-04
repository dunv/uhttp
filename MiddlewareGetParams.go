package uhttp

import (
	"context"
	"fmt"
	"net/http"
)

func getParamsMiddleware(optionalGet R, requiredGet R) func(next http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			actualRaw := r.URL.Query() // map[string][]string
			actual := map[string]string{}
			for key, values := range actualRaw {
				if len(values) > 0 {
					if len(values) > 1 {
						RenderError(w, r, fmt.Errorf("param %s has more than one value (given multiple times), this is not supported", key))
						return
					}
					actual[key] = values[0]
				}
			}

			paramMap := R{}
			err := ValidateParams(requiredGet, actual, paramMap, true)
			if err != nil {
				RenderError(w, r, fmt.Errorf("%v", err))
				return
			}

			err = ValidateParams(optionalGet, actual, paramMap, false)
			if err != nil {
				RenderError(w, r, fmt.Errorf("%v", err))
			}

			ctx := context.WithValue(r.Context(), CtxKeyGetParams, paramMap)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	}
}
