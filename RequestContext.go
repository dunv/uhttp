package uhttp

import (
	"github.com/dunv/uhelpers"
	"github.com/dunv/ulog"
)

type ContextMap map[string]interface{}

var requestContext ContextMap = ContextMap{}

func RequestContext() ContextMap {
	return requestContext
}

func AddContext(key string, value interface{}) {
	keys := uhelpers.StringKeysFromMap(requestContext)
	if !uhelpers.SliceContainsItem(keys, key) {
		requestContext[key] = value
	} else {
		ulog.Warnf("contextKey %s already exists", key)
	}
}
