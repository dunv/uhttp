package uhttp_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/dunv/uhttp"
)

func testMiddleware(key uhttp.ContextKey, value interface{}) uhttp.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			httpContext := context.WithValue(r.Context(), key, value)
			next.ServeHTTP(w, r.WithContext(httpContext))
		}
	}
}

func TestAdditionalMiddlewareGlobally(t *testing.T) {
	ctxKey := uhttp.ContextKey("manuallyAddedGlobally")
	middleware := testMiddleware(ctxKey, map[string]string{"manuallyAdded": "manuallyAdded"})
	u := uhttp.NewUHTTP(uhttp.WithGlobalMiddlewares(middleware))

	handler := uhttp.NewHandler(uhttp.WithGet(func(r *http.Request, ret *int) interface{} {
		return r.Context().Value(ctxKey)
	}))

	expectedResponseBody := []byte(`{"manuallyAdded":"manuallyAdded"}`)

	executeHandler(handler, http.MethodGet, http.StatusOK, nil, expectedResponseBody, u, t)
}

func TestAdditionalMiddlewareHandlerSingle(t *testing.T) {
	u := uhttp.NewUHTTP()
	ctxKey := uhttp.ContextKey("manuallyAddedSingleHandler")
	middleware := testMiddleware(ctxKey, map[string]string{"manuallyAdded": "manuallyAdded"})
	handler := uhttp.NewHandler(
		uhttp.WithGet(func(r *http.Request, ret *int) interface{} {
			return r.Context().Value(ctxKey)
		}),
		uhttp.WithMiddlewares(middleware),
	)

	expectedResponseBody := []byte(`{"manuallyAdded":"manuallyAdded"}`)

	executeHandler(handler, http.MethodGet, http.StatusOK, nil, expectedResponseBody, u, t)
}

func TestAdditionalMiddlewareHandlerMultiple(t *testing.T) {
	u := uhttp.NewUHTTP()

	ctxKey := uhttp.ContextKey("manuallyAddedMultipleHandler")
	middleware := uhttp.WithContextMiddleware(ctxKey, map[string]string{"manuallyAdded": "manuallyAdded"})

	handler := uhttp.NewHandler(
		uhttp.WithGet(func(r *http.Request, ret *int) interface{} {
			return r.Context().Value(ctxKey)
		}),
		uhttp.WithMiddlewares(middleware),
	)

	expectedResponseBody := []byte(`{"manuallyAdded":"manuallyAdded"}`)

	executeHandler(handler, http.MethodGet, http.StatusOK, nil, expectedResponseBody, u, t)
}
