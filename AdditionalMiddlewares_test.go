package uhttp

import (
	"context"
	"net/http"
	"testing"
)

func testMiddleware(key string, value interface{}) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			httpContext := context.WithValue(r.Context(), key, value)
			next.ServeHTTP(w, r.WithContext(httpContext))
		}
	}
}

func TestAdditionalMiddlewareGlobally(t *testing.T) {
	middleware := testMiddleware("manuallyAddedGlobally", map[string]string{"manuallyAdded": "manuallyAdded"})
	u := NewUHTTP(WithGlobalMiddlewares([]Middleware{middleware}))

	handler := NewHandler(WithGet(func(r *http.Request, ret *int) interface{} {
		return r.Context().Value("manuallyAddedGlobally")
	}))

	expectedResponseBody := []byte(`{"manuallyAdded":"manuallyAdded"}`)

	ExecuteHandler(handler, http.MethodGet, http.StatusOK, nil, expectedResponseBody, u, t)
}

func TestAdditionalMiddlewareHandlerSingle(t *testing.T) {
	u := NewUHTTP()
	middleware := testMiddleware("manuallyAddedSingleHandler", map[string]string{"manuallyAdded": "manuallyAdded"})
	handler := NewHandler(
		WithGet(func(r *http.Request, ret *int) interface{} {
			return r.Context().Value("manuallyAddedSingleHandler")
		}),
		WithMiddlewares([]Middleware{middleware}),
	)

	expectedResponseBody := []byte(`{"manuallyAdded":"manuallyAdded"}`)

	ExecuteHandler(handler, http.MethodGet, http.StatusOK, nil, expectedResponseBody, u, t)
}

func TestAdditionalMiddlewareHandlerMultiple(t *testing.T) {
	u := NewUHTTP()

	middleware := withContextMiddleware(u, "manuallyAddedMultipleHandler", map[string]string{"manuallyAdded": "manuallyAdded"})

	handler := NewHandler(
		WithGet(func(r *http.Request, ret *int) interface{} {
			return r.Context().Value("manuallyAddedMultipleHandler")
		}),
		WithMiddlewares([]Middleware{middleware}),
	)

	expectedResponseBody := []byte(`{"manuallyAdded":"manuallyAdded"}`)

	ExecuteHandler(handler, http.MethodGet, http.StatusOK, nil, expectedResponseBody, u, t)
}
