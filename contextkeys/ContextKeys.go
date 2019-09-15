package contextkeys

type ContextKey string

const (
	// TODO: move into uauth, provide generic mechanism for adding context
	CtxKeyBCryptSecret ContextKey = "uauth__bCryptSecret"
	CtxKeyPostModel    ContextKey = "uhttp__postModel"
	CtxKeyGetParams    ContextKey = "uhttp__getParams"
)
