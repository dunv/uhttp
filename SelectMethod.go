package uhttp

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/dunv/ulog"
)

func SelectMethod(u *UHTTP, chain Middleware, handlerOpts handlerOptions) http.HandlerFunc {
	return chain(func(w http.ResponseWriter, r *http.Request) {

		// Figure out which method to invoke
		var returnCode int

		// We need to process the handler in a goroutine so we can recover from panics
		// this channel will be used to tell the main routine that the handler was processed
		handlerProcessed := make(chan interface{})

		if r.Method == http.MethodGet && handlerOpts.Get != nil {
			go func() {
				defer recoverFromPanic(u, handlerProcessed, r, &returnCode)
				handlerProcessed <- handlerOpts.Get(r, &returnCode)
			}()
		} else if r.Method == http.MethodGet && handlerOpts.GetWithModel != nil {
			go func() {
				defer recoverFromPanic(u, handlerProcessed, r, &returnCode)
				model := parsedModel(r)
				handlerProcessed <- handlerOpts.GetWithModel(r, model, &returnCode)
			}()
		} else if r.Method == http.MethodPost && handlerOpts.Post != nil {
			go func() {
				defer recoverFromPanic(u, handlerProcessed, r, &returnCode)
				handlerProcessed <- handlerOpts.Post(r, &returnCode)
			}()
		} else if r.Method == http.MethodPost && handlerOpts.PostWithModel != nil {
			go func() {
				defer recoverFromPanic(u, handlerProcessed, r, &returnCode)
				model := parsedModel(r)
				handlerProcessed <- handlerOpts.PostWithModel(r, model, &returnCode)
			}()
		} else if r.Method == http.MethodDelete && handlerOpts.Delete != nil {
			go func() {
				defer recoverFromPanic(u, handlerProcessed, r, &returnCode)
				handlerProcessed <- handlerOpts.Delete(r, &returnCode)
			}()
		} else if r.Method == http.MethodDelete && handlerOpts.DeleteWithModel != nil {
			go func() {
				defer recoverFromPanic(u, handlerProcessed, r, &returnCode)
				model := parsedModel(r)
				handlerProcessed <- handlerOpts.DeleteWithModel(r, model, &returnCode)
			}()
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
			u.RenderError(w, r, fmt.Errorf("method not allowed"))
			return
		}

		res := <-handlerProcessed

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

func recoverFromPanic(u *UHTTP, handlerProcessed chan interface{}, r *http.Request, returnCode *int) {
	if rec := recover(); rec != nil {
		err := fmt.Errorf("panic: handlerExecution (%s)", rec)
		u.opts.log.Errorf("panic [path: %s] %s", r.RequestURI, err)
		stack := debug.Stack()
		ulog.LogByteArrayLineByLine(stack, u.opts.log.Errorf, fmt.Sprintf("panic [path: %s] ", r.RequestURI))
		err = fmt.Errorf("%s stackTrace: %s", err, strings.ReplaceAll(string(stack), "\n", "\\n"))
		*returnCode = http.StatusInternalServerError
		if u.opts.sendPanicInfoToClient {
			handlerProcessed <- err
			return
		}
		handlerProcessed <- fmt.Errorf("internal server error")
	}
}
