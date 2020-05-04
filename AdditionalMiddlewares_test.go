package uhttp

import (
	"net/http"
	"testing"
)

func TestAdditionalMiddlewareGlobally(t *testing.T) {
	u := NewUHTTP()

	middleware := WithContextMiddleware("manuallyAddedGlobally", map[string]string{"manuallyAdded": "manuallyAdded"})
	err := AddMiddleware(middleware)
	if err != nil {
		t.Error(err)
	}

	mws := AdditionalMiddlewares()
	if len(mws) != 1 {
		t.Errorf("did not correctly keep track of middlewares")
	}

	handler := Handler{
		GetHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			Render(w, r, r.Context().Value("manuallyAddedGlobally"))
		}),
	}
	expectedResponseBody := []byte(`{"manuallyAdded":"manuallyAdded"}`)

	ExecuteHandler(handler, http.MethodGet, http.StatusOK, nil, expectedResponseBody, u, t)
}

func TestAdditionalMiddlewareHandlerSingle(t *testing.T) {
	u := NewUHTTP()

	middleware := WithContextMiddleware("manuallyAddedSingleHandler", map[string]string{"manuallyAdded": "manuallyAdded"})

	handler := Handler{
		GetHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			Render(w, r, r.Context().Value("manuallyAddedSingleHandler"))
		}),
		AddMiddleware: &middleware,
	}
	expectedResponseBody := []byte(`{"manuallyAdded":"manuallyAdded"}`)

	ExecuteHandler(handler, http.MethodGet, http.StatusOK, nil, expectedResponseBody, u, t)
}

func TestAdditionalMiddlewareHandlerMultiple(t *testing.T) {
	u := NewUHTTP()

	middleware := WithContextMiddleware("manuallyAddedMultipleHandler", map[string]string{"manuallyAdded": "manuallyAdded"})

	handler := Handler{
		GetHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			Render(w, r, r.Context().Value("manuallyAddedMultipleHandler"))
		}),
		AddMiddlewares: []Middleware{middleware},
	}
	expectedResponseBody := []byte(`{"manuallyAdded":"manuallyAdded"}`)

	ExecuteHandler(handler, http.MethodGet, http.StatusOK, nil, expectedResponseBody, u, t)
}
