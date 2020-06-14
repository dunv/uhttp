package uhttp

import (
	"context"
	"net/http"
	"testing"
)

func testMiddleware(key ContextKey, value interface{}) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			httpContext := context.WithValue(r.Context(), key, value)
			next.ServeHTTP(w, r.WithContext(httpContext))
		}
	}
}

func TestAdditionalMiddlewareGlobally(t *testing.T) {
	ctxKey := ContextKey("manuallyAddedGlobally")
	middleware := testMiddleware(ctxKey, map[string]string{"manuallyAdded": "manuallyAdded"})
	u := NewUHTTP(WithGlobalMiddlewares([]Middleware{middleware}))

	handler := NewHandler(WithGet(func(r *http.Request, ret *int) interface{} {
		return r.Context().Value(ctxKey)
	}))

	expectedResponseBody := []byte(`{"manuallyAdded":"manuallyAdded"}`)

	ExecuteHandler(handler, http.MethodGet, http.StatusOK, nil, expectedResponseBody, u, t)
}

func TestAdditionalMiddlewareHandlerSingle(t *testing.T) {
	u := NewUHTTP()
	ctxKey := ContextKey("manuallyAddedSingleHandler")
	middleware := testMiddleware(ctxKey, map[string]string{"manuallyAdded": "manuallyAdded"})
	handler := NewHandler(
		WithGet(func(r *http.Request, ret *int) interface{} {
			return r.Context().Value(ctxKey)
		}),
		WithMiddlewares([]Middleware{middleware}),
	)

	expectedResponseBody := []byte(`{"manuallyAdded":"manuallyAdded"}`)

	ExecuteHandler(handler, http.MethodGet, http.StatusOK, nil, expectedResponseBody, u, t)
}

func TestAdditionalMiddlewareHandlerMultiple(t *testing.T) {
	u := NewUHTTP()

	ctxKey := ContextKey("manuallyAddedMultipleHandler")
	middleware := withContextMiddleware(u, ctxKey, map[string]string{"manuallyAdded": "manuallyAdded"})

	handler := NewHandler(
		WithGet(func(r *http.Request, ret *int) interface{} {
			return r.Context().Value(ctxKey)
		}),
		WithMiddlewares([]Middleware{middleware}),
	)

	expectedResponseBody := []byte(`{"manuallyAdded":"manuallyAdded"}`)

	ExecuteHandler(handler, http.MethodGet, http.StatusOK, nil, expectedResponseBody, u, t)
}
