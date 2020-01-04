package uhttp

import (
	"fmt"

	"github.com/dunv/uhelpers"
)

var additionalMiddlewares []Middleware = []Middleware{}

func AdditionalMiddlewares() []Middleware {
	return additionalMiddlewares
}

func AddMiddleware(mw Middleware) error {
	if !uhelpers.SliceContainsItem(additionalMiddlewares, mw) {
		additionalMiddlewares = append(additionalMiddlewares, mw)
		return nil
	}
	return fmt.Errorf("middleware already added (%+v)", mw)
}
