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
func withParams(params Params, required bool, customLog *CustomLogger) Middleware {
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
						RenderMessageWithStatusCode(w, r, 400, fmt.Sprintf("Param %s is required", paramName), customLog)
						return
					}
					if paramRequirement.AllValues {
						var tmp *string
						paramMap.(map[string]interface{})[paramName] = tmp
					} else if paramRequirement.Enum != nil {
						var tmp *string
						paramMap.(map[string]interface{})[paramName] = tmp
					} else if paramRequirement.Date {
						var tmp *time.Time
						paramMap.(map[string]interface{})[paramName] = tmp
					} else if paramRequirement.Int {
						var tmp *int
						paramMap.(map[string]interface{})[paramName] = tmp
					}
				} else {
					paramValue := keys[0]
					if !paramRequirement.AllValues {
						validated := false

						if paramRequirement.Enum != nil {
							for _, enumValue := range paramRequirement.Enum {
								if enumValue == paramValue {
									validated = true
									if required {
										paramMap.(map[string]interface{})[paramName] = enumValue
									} else {
										paramMap.(map[string]interface{})[paramName] = &enumValue
									}
								}
							}
						} else if paramRequirement.Date {
							var err error
							timeValue, err := time.Parse(time.RFC3339, paramValue)
							if err != nil {
								RenderMessageWithStatusCode(w, r, 400, fmt.Sprintf("Param %s has to be a date (%s), error %s", paramName, paramValue, err), customLog)
								return
							}
							if required {
								paramMap.(map[string]interface{})[paramName] = timeValue
							} else {
								paramMap.(map[string]interface{})[paramName] = &timeValue
							}
							validated = true
						} else if paramRequirement.Int {
							var err error
							intValue, err := strconv.ParseInt(paramValue, 10, 64)
							if err != nil {
								RenderMessageWithStatusCode(w, r, 400, fmt.Sprintf("Param %s has to be an integer (%s), error %s", paramName, paramValue, err), customLog)
								return
							}
							if required {
								paramMap.(map[string]interface{})[paramName] = intValue
							} else {
								paramMap.(map[string]interface{})[paramName] = &intValue
							}
							validated = true
						}

						if !validated {
							RenderMessageWithStatusCode(w, r, 400, fmt.Sprintf("Param %s can only assume these values %s", paramName, paramRequirement.Enum), customLog)
							return
						}
					} else {
						if required {
							paramMap.(map[string]interface{})[paramName] = paramValue
						} else {
							paramMap.(map[string]interface{})[paramName] = &paramValue
						}
					}
				}
			}
			ctx := context.WithValue(r.Context(), CtxKeyParams, paramMap)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	}
}

// WithOptionalParams parses and adds optional params to request
func WithOptionalParams(params Params, customLog *CustomLogger) Middleware {
	return withParams(params, false, customLog)
}

// WithRequiredParams parses and adds required params to request
func WithRequiredParams(params Params, customLog *CustomLogger) Middleware {
	return withParams(params, true, customLog)
}
