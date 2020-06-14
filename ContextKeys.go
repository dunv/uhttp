package uhttp

type ContextKey string

const (
	CtxKeyPostModel ContextKey = "uhttp.postModel"
	CtxKeyGetParams ContextKey = "uhttp.getParams"
)
