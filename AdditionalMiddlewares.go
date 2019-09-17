package uhttp

import (
	"github.com/dunv/uhelpers"
	"github.com/dunv/uhttp/logging"
	"github.com/dunv/uhttp/models"
)

var additionalMiddlewares []models.Middleware = []models.Middleware{}

func AdditionalMiddlewares() []models.Middleware {
	return additionalMiddlewares
}

func AddMiddleware(mw models.Middleware) {
	if !uhelpers.SliceContainsItem(additionalMiddlewares, mw) {
		additionalMiddlewares = append(additionalMiddlewares, mw)
		return
	}
	logging.Logger.Warnf("middleware already added", mw)
}
