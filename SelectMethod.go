package uhttp

import (
	"fmt"
	"net/http"
)

func SelectMethod(u *UHTTP, chain Middleware, handlerOpts handlerOptions) http.HandlerFunc {
	return chain(func(w http.ResponseWriter, r *http.Request) {

		// Figure out which method to invoke
		var res interface{}
		var returnCode int
		if r.Method == http.MethodGet && handlerOpts.Get != nil {
			res = handlerOpts.Get(r, &returnCode)
		} else if r.Method == http.MethodGet && handlerOpts.GetWithModel != nil {
			model := parsedModel(r)
			res = handlerOpts.GetWithModel(r, model, &returnCode)
		} else if r.Method == http.MethodPost && handlerOpts.Post != nil {
			res = handlerOpts.Post(r, &returnCode)
		} else if r.Method == http.MethodPost && handlerOpts.PostWithModel != nil {
			model := parsedModel(r)
			res = handlerOpts.PostWithModel(r, model, &returnCode)
		} else if r.Method == http.MethodDelete && handlerOpts.Delete != nil {
			res = handlerOpts.Delete(r, &returnCode)
		} else if r.Method == http.MethodDelete && handlerOpts.DeleteWithModel != nil {
			model := parsedModel(r)
			res = handlerOpts.DeleteWithModel(r, model, &returnCode)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
			u.RenderError(w, r, fmt.Errorf("method not allowed"))
			return
		}

		// Figure out, how to respond
		if res != nil {
			switch typed := res.(type) {
			case error:
				if returnCode == 0 {
					u.RenderErrorWithStatusCode(w, r, http.StatusBadRequest, typed, true)
				} else {
					u.RenderErrorWithStatusCode(w, r, returnCode, typed, true)
				}
			default:
				if returnCode == 0 {
					u.RenderWithStatusCode(w, r, http.StatusOK, typed)
				} else {
					u.RenderWithStatusCode(w, r, returnCode, typed)
				}
			}
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})
}
