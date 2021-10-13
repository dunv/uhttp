package uhttp

import (
	"net/http"
)

func selectMethodMiddleware(u *UHTTP, handlerOpts handlerOptions) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, retCode := executeHandlerMethod(r, u, handlerOpts)

		// Figure out, how to respond
		if res != nil {
			switch typed := res.(type) {
			case error:
				u.RenderErrorWithStatusCode(w, r, retCode, typed, u.opts.logHandlerErrors)
			default:
				u.RenderWithStatusCode(w, r, retCode, typed)
			}
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
