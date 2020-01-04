package uhttp

import (
	"github.com/dunv/uhelpers"
)

var additionalMiddlewares []Middleware = []Middleware{}

func AdditionalMiddlewares() []Middleware {
	return additionalMiddlewares
}

func AddMiddleware(mw Middleware) {
	if !uhelpers.SliceContainsItem(additionalMiddlewares, mw) {
		additionalMiddlewares = append(additionalMiddlewares, mw)
		return
	}
	Logger.Warnf("middleware already added", mw)
}
