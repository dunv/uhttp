package contextkeys

type ContextKey string

const (
	CtxKeyBCryptSecret ContextKey = "bCryptSecret"
	CtxKeyPostModel    ContextKey = "postModel"
	CtxKeyParams       ContextKey = "params"
	CtxKeyGetParams    ContextKey = "getParams"
)
