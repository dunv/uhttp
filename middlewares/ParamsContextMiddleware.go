package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/dunv/uhttp/contextkeys"
	"github.com/dunv/uhttp/helpers"
	"github.com/dunv/uhttp/models"
)

// WithParams parses and adds params to request
func withParams(params models.Params, required bool) models.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			paramMap := r.Context().Value(contextkeys.CtxKeyParams)
			if paramMap == nil {
				paramMap = map[string]interface{}{}
			}

			for paramName, paramRequirement := range params.ParamMap {
				keys, ok := r.URL.Query()[paramName]
				if !ok || len(keys) < 1 {
					if required {
						helpers.RenderMessageWithStatusCode(w, r, 400, fmt.Sprintf("Param %s is required", paramName))
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
					} else if paramRequirement.ShortDate {
						var tmp *time.Time
						paramMap.(map[string]interface{})[paramName] = tmp
					} else if paramRequirement.Int {
						var tmp *int64
						paramMap.(map[string]interface{})[paramName] = tmp
					} else if paramRequirement.Float {
						var tmp *float64
						paramMap.(map[string]interface{})[paramName] = tmp
					} else if paramRequirement.Bool {
						var tmp *bool
						paramMap.(map[string]interface{})[paramName] = tmp
					}

				} else {
					paramValue := keys[0]
					if !paramRequirement.AllValues {
						validated := false

						if paramRequirement.Enum != nil {
							for index, enumValue := range paramRequirement.Enum {
								if enumValue == paramValue {
									validated = true
									if required {
										paramMap.(map[string]interface{})[paramName] = paramRequirement.Enum[index]
									} else {
										paramMap.(map[string]interface{})[paramName] = &paramRequirement.Enum[index]
									}
								}
							}
						} else if paramRequirement.Date {
							var err error
							timeValue, err := time.Parse(time.RFC3339, paramValue)
							if err != nil {
								helpers.RenderMessageWithStatusCode(w, r, 400, fmt.Sprintf("Param %s has to be a date (%s), error %s", paramName, paramValue, err))
								return
							}
							if required {
								paramMap.(map[string]interface{})[paramName] = timeValue
							} else {
								paramMap.(map[string]interface{})[paramName] = &timeValue
							}
							validated = true
						} else if paramRequirement.ShortDate {
							var err error
							timeValue, err := time.Parse("2006-01-02", paramValue)
							if err != nil {
								helpers.RenderMessageWithStatusCode(w, r, 400, fmt.Sprintf("Param %s has to be a date (%s), error %s", paramName, paramValue, err))
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
								helpers.RenderMessageWithStatusCode(w, r, 400, fmt.Sprintf("Param %s has to be an integer (%s), error %s", paramName, paramValue, err))
								return
							}
							if required {
								paramMap.(map[string]interface{})[paramName] = intValue
							} else {
								paramMap.(map[string]interface{})[paramName] = &intValue
							}
							validated = true
						} else if paramRequirement.Float {
							var err error
							floatValue, err := strconv.ParseFloat(paramValue, 64)
							if err != nil {
								helpers.RenderMessageWithStatusCode(w, r, 400, fmt.Sprintf("Param %s has to be a float (%s), error %s", paramName, paramValue, err))
								return
							}
							if required {
								paramMap.(map[string]interface{})[paramName] = floatValue
							} else {
								paramMap.(map[string]interface{})[paramName] = &floatValue
							}
							validated = true
						} else if paramRequirement.Bool {
							var err error
							boolValue, err := strconv.ParseBool(paramValue)
							if err != nil {
								helpers.RenderMessageWithStatusCode(w, r, 400, fmt.Sprintf("Param %s has to be a bool (%s), error %s", paramName, paramValue, err))
								return
							}
							if required {
								paramMap.(map[string]interface{})[paramName] = boolValue
							} else {
								paramMap.(map[string]interface{})[paramName] = &boolValue
							}
							validated = true
						}

						if !validated {
							helpers.RenderMessageWithStatusCode(w, r, 400, fmt.Sprintf("Param %s can only assume these values %s", paramName, paramRequirement.Enum))
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
			ctx := context.WithValue(r.Context(), contextkeys.CtxKeyParams, paramMap)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	}
}

// WithOptionalParams parses and adds optional params to request
func WithOptionalParams(params models.Params) models.Middleware {
	return withParams(params, false)
}

// WithRequiredParams parses and adds required params to request
func WithRequiredParams(params models.Params) models.Middleware {
	return withParams(params, true)
}