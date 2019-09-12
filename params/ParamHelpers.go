package params

import (
	"net/http"
	"time"

	"github.com/dunv/uhttp/contextkeys"
	"github.com/dunv/ulog"
)

// A long list of helpers with identical code... I am missing generics here
// I just do not want this error-checking-code in all my handlers

func GetAsString(key string, r *http.Request) *string {
	// Publish an error only in the logs, if the key is not present in the request context
	// it obviously points to a bug in the code not an error on the user's side
	paramMap, ok := r.Context().Value(contextkeys.CtxKeyGetParams).(R)
	if !ok {
		ulog.Error("contextkeys.ContextKeyGetParams is not present in the request's context. please check the handler's definition")
		return nil
	}

	stringValue, ok := paramMap[key].(string)
	if !ok {
		return nil
	}

	return &stringValue
}

func GetAsBool(key string, r *http.Request) *bool {
	paramMap, ok := r.Context().Value(contextkeys.CtxKeyGetParams).(R)
	if !ok {
		ulog.Error("contextkeys.ContextKeyGetParams is not present in the request's context. please check the handler's definition")
		return nil
	}

	value, ok := paramMap[key].(bool)
	if !ok {
		return nil
	}

	return &value
}

func GetAsInt(key string, r *http.Request) *int {
	paramMap, ok := r.Context().Value(contextkeys.CtxKeyGetParams).(R)
	if !ok {
		ulog.Error("contextkeys.ContextKeyGetParams is not present in the request's context. please check the handler's definition")
		return nil
	}

	value, ok := paramMap[key].(int)
	if !ok {
		return nil
	}

	return &value
}

func GetAsInt32(key string, r *http.Request) *int32 {
	paramMap, ok := r.Context().Value(contextkeys.CtxKeyGetParams).(R)
	if !ok {
		ulog.Error("contextkeys.ContextKeyGetParams is not present in the request's context. please check the handler's definition")
		return nil
	}

	value, ok := paramMap[key].(int32)
	if !ok {
		return nil
	}

	return &value
}

func GetAsInt64(key string, r *http.Request) *int64 {
	paramMap, ok := r.Context().Value(contextkeys.CtxKeyGetParams).(R)
	if !ok {
		ulog.Error("contextkeys.ContextKeyGetParams is not present in the request's context. please check the handler's definition")
		return nil
	}

	value, ok := paramMap[key].(int64)
	if !ok {
		return nil
	}

	return &value
}

func GetAsFloat32(key string, r *http.Request) *float32 {
	paramMap, ok := r.Context().Value(contextkeys.CtxKeyGetParams).(R)
	if !ok {
		ulog.Error("contextkeys.ContextKeyGetParams is not present in the request's context. please check the handler's definition")
		return nil
	}

	value, ok := paramMap[key].(float32)
	if !ok {
		return nil
	}

	return &value
}

func GetAsTime(key string, r *http.Request) *time.Time {
	paramMap, ok := r.Context().Value(contextkeys.CtxKeyGetParams).(R)
	if !ok {
		ulog.Error("contextkeys.ContextKeyGetParams is not present in the request's context. please check the handler's definition")
		return nil
	}

	value, ok := paramMap[key].(time.Time)
	if !ok {
		return nil
	}

	return &value
}
