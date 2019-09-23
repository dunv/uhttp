package uhttp

import (
	"github.com/dunv/uhelpers"
	"github.com/dunv/uhttp/logging"
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
	logging.Logger.Warnf("middleware already added", mw)
}
