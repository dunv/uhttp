package uhttp

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// Params <-
type Params struct {
	ParamMap map[string]ParamRequirement
}

// ParamRequirement <-
type ParamRequirement struct {
	AllValues bool
	Date      bool
	Enum      []string
	Int       bool
}

// CtxKeyParams is the context key to retrieve the params
const CtxKeyParams = ContextKey("params")

// WithParams parses and adds params to request
func withParams(params Params, required bool) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			paramMap := r.Context().Value(CtxKeyParams)
			if paramMap == nil {
				paramMap = map[string]interface{}{}
			}

			for paramName, paramRequirement := range params.ParamMap {
				keys, ok := r.URL.Query()[paramName]
				if !ok || len(keys) < 1 {
					if required {
						RenderMessageWithStatusCode(w, r, 400, fmt.Sprintf("Param %s is required", paramName))
						return
					}
				} else {
					paramValue := keys[0]
					var parsedValue interface{}
					if !paramRequirement.AllValues {
						validated := false

						if paramRequirement.Enum != nil {
							for _, enumValue := range paramRequirement.Enum {
								if enumValue == paramValue {
									validated = true
									parsedValue = enumValue
								}
							}
						} else if paramRequirement.Date {
							var err error
							parsedValue, err = time.Parse(time.RFC3339, paramValue)

							if err != nil {
								RenderMessageWithStatusCode(w, r, 400, fmt.Sprintf("Param %s has to be a date (%s), error %s", paramName, paramValue, err))
								return
							}

							validated = true
						} else if paramRequirement.Int {
							var err error
							parsedValue, err = strconv.ParseInt(paramValue, 10, 64)
							if err != nil {
								RenderMessageWithStatusCode(w, r, 400, fmt.Sprintf("Param %s has to be an integer (%s), error %s", paramName, paramValue, err))
								return
							}
							validated = true
						}

						if !validated {
							RenderMessageWithStatusCode(w, r, 400, fmt.Sprintf("Param %s can only assume these values %s", paramName, paramRequirement.Enum))
							return
						}
					} else {
						parsedValue = paramValue
					}
					paramMap.(map[string]interface{})[paramName] = parsedValue
				}
			}
			ctx := context.WithValue(r.Context(), CtxKeyParams, paramMap)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	}
}

// WithOptionalParams parses and adds optional params to request
func WithOptionalParams(params Params) Middleware {
	return withParams(params, false)
}

// WithRequiredParams parses and adds required params to request
func WithRequiredParams(params Params) Middleware {
	return withParams(params, true)
}
