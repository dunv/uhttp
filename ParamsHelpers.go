package uhttp 

import (
	"context"
	"net/http"
	"time"
)

// A long list of helpers with identical code... I am missing generics here
// I just do not want this error-checking-code in all my handlers

func GetAsString(key string, r *http.Request) *string {
	return GetAsStringFromContext(key, r.Context())
}

func GetAsStringFromContext(key string, ctx context.Context) *string {
	// Publish an error only in the logs, if the key is not present in the request context
	// it obviously points to a bug in the code not an error on the user's side
	paramMap, ok := ctx.Value(CtxKeyGetParams).(R)
	if !ok {
		Logger.Errorf("ContextKeyGetParams is not present in the request's context. please check the handler's definition")
		return nil
	}

	stringValue, ok := paramMap[key].(string)
	if !ok {
		return nil
	}

	return &stringValue
}

func GetAsBool(key string, r *http.Request) *bool {
	return GetAsBoolFromContext(key, r.Context())
}

func GetAsBoolFromContext(key string, ctx context.Context) *bool {
	paramMap, ok := ctx.Value(CtxKeyGetParams).(R)
	if !ok {
		Logger.Errorf("ContextKeyGetParams is not present in the request's context. please check the handler's definition")
		return nil
	}

	value, ok := paramMap[key].(bool)
	if !ok {
		return nil
	}

	return &value
}

func GetAsInt(key string, r *http.Request) *int {
	return GetAsIntFromContext(key, r.Context())
}

func GetAsIntFromContext(key string, ctx context.Context) *int {
	paramMap, ok := ctx.Value(CtxKeyGetParams).(R)
	if !ok {
		Logger.Errorf("ContextKeyGetParams is not present in the request's context. please check the handler's definition")
		return nil
	}

	value, ok := paramMap[key].(int)
	if !ok {
		return nil
	}

	return &value
}

func GetAsInt32(key string, r *http.Request) *int32 {
	return GetAsInt32FromContext(key, r.Context())
}

func GetAsInt32FromContext(key string, ctx context.Context) *int32 {
	paramMap, ok := ctx.Value(CtxKeyGetParams).(R)
	if !ok {
		Logger.Errorf("ContextKeyGetParams is not present in the request's context. please check the handler's definition")
		return nil
	}

	value, ok := paramMap[key].(int32)
	if !ok {
		return nil
	}

	return &value
}

func GetAsInt64(key string, r *http.Request) *int64 {
	return GetAsInt64FromContext(key, r.Context())
}

func GetAsInt64FromContext(key string, ctx context.Context) *int64 {
	paramMap, ok := ctx.Value(CtxKeyGetParams).(R)
	if !ok {
		Logger.Errorf("ContextKeyGetParams is not present in the request's context. please check the handler's definition")
		return nil
	}

	value, ok := paramMap[key].(int64)
	if !ok {
		return nil
	}

	return &value
}

func GetAsFloat32(key string, r *http.Request) *float32 {
	return GetAsFloat32FromContext(key, r.Context())
}

func GetAsFloat32FromContext(key string, ctx context.Context) *float32 {
	paramMap, ok := ctx.Value(CtxKeyGetParams).(R)
	if !ok {
		Logger.Errorf("ContextKeyGetParams is not present in the request's context. please check the handler's definition")
		return nil
	}

	value, ok := paramMap[key].(float32)
	if !ok {
		return nil
	}

	return &value
}

func GetAsFloat64(key string, r *http.Request) *float64 {
	return GetAsFloat64FromContext(key, r.Context())
}

func GetAsFloat64FromContext(key string, ctx context.Context) *float64 {
	paramMap, ok := ctx.Value(CtxKeyGetParams).(R)
	if !ok {
		Logger.Errorf("ContextKeyGetParams is not present in the request's context. please check the handler's definition")
		return nil
	}

	value, ok := paramMap[key].(float64)
	if !ok {
		return nil
	}

	return &value
}

func GetAsTime(key string, r *http.Request) *time.Time {
	return GetAsTimeFromContext(key, r.Context())
}

func GetAsTimeFromContext(key string, ctx context.Context) *time.Time {
	paramMap, ok := ctx.Value(CtxKeyGetParams).(R)
	if !ok {
		Logger.Errorf("ContextKeyGetParams is not present in the request's context. please check the handler's definition")
		return nil
	}

	value, ok := paramMap[key].(time.Time)
	if !ok {
		return nil
	}

	return &value
}
