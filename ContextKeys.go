package uhttp

import "net/http"

type ContextKey string

const (
	CtxKeyPostModel                 ContextKey = "uhttp.postModel"
	CtxKeyIsAutomaticCacheExecution ContextKey = "uhttp.isAutomaticCacheExecution"
	CtxKeyCache                     ContextKey = "uhttp.cache"
	CtxKeyGetParams                 ContextKey = "uhttp.getParams"
	CtxKeyResponseWriter            ContextKey = "uhttp.responseWriter"
)

func IsAutomaticCacheExecution(r *http.Request) bool {
	if val := r.Context().Value(CtxKeyIsAutomaticCacheExecution); val != nil {
		if isAutomaticCacheExecution, ok := val.(bool); ok {
			return isAutomaticCacheExecution
		}
	}
	return false
}
