package uhttp

type ContextKey string

const (
	CtxKeyPostModel      ContextKey = "uhttp.postModel"
	CtxKeyCache          ContextKey = "uhttp.cache"
	CtxKeyGetParams      ContextKey = "uhttp.getParams"
	CtxKeyResponseWriter ContextKey = "uhttp.responseWriter"
)
