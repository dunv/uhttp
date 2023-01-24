package uhttp

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/dunv/uhelpers"
)

func selectMethodMiddleware(u *UHTTP, handlerOpts handlerOptions) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, retCode := executeHandlerMethod(r, u, handlerOpts)

		// Figure out how to respond
		// if res == nil -> do not do ANYTHING with response (no header, no body)
		if res != nil {
			switch typed := res.(type) {
			case error:
				u.RenderErrorWithStatusCode(w, r, retCode, typed, u.opts.logHandlerErrors)
			default:
				u.RenderWithStatusCode(w, r, retCode, typed)
			}
			return
		}
	}
}

func executeHandlerMethod(r *http.Request, u *UHTTP, handlerOpts handlerOptions) (interface{}, int) {
	// Figure out which method to invoke
	var returnCode int

	// We need to process the handler in a goroutine so we can recover from panics
	// this channel will be used to tell the main routine that the handler was processed
	handlerProcessed := make(chan interface{})

	if r.Method == http.MethodGet && handlerOpts.get != nil {
		go func() {
			defer recoverFromPanic(u, handlerProcessed, r, &returnCode)
			handlerProcessed <- handlerOpts.get(r, &returnCode)
		}()
	} else if r.Method == http.MethodGet && handlerOpts.getWithModel != nil {
		go func() {
			defer recoverFromPanic(u, handlerProcessed, r, &returnCode)
			model := parsedModel(r)
			handlerProcessed <- handlerOpts.getWithModel(r, model, &returnCode)
		}()
	} else if r.Method == http.MethodPost && handlerOpts.post != nil {
		go func() {
			defer recoverFromPanic(u, handlerProcessed, r, &returnCode)
			handlerProcessed <- handlerOpts.post(r, &returnCode)
		}()
	} else if r.Method == http.MethodPost && handlerOpts.postWithModel != nil {
		go func() {
			defer recoverFromPanic(u, handlerProcessed, r, &returnCode)
			model := parsedModel(r)
			handlerProcessed <- handlerOpts.postWithModel(r, model, &returnCode)
		}()
	} else if r.Method == http.MethodDelete && handlerOpts.delete != nil {
		go func() {
			defer recoverFromPanic(u, handlerProcessed, r, &returnCode)
			handlerProcessed <- handlerOpts.delete(r, &returnCode)
		}()
	} else if r.Method == http.MethodDelete && handlerOpts.deleteWithModel != nil {
		go func() {
			defer recoverFromPanic(u, handlerProcessed, r, &returnCode)
			model := parsedModel(r)
			handlerProcessed <- handlerOpts.deleteWithModel(r, model, &returnCode)
		}()
	} else {
		return fmt.Errorf("method not allowed"), http.StatusMethodNotAllowed
	}

	res := <-handlerProcessed
	if res != nil {
		switch res.(type) {
		case error:
			if returnCode == 0 {
				returnCode = http.StatusBadRequest
			}
		default:
			if returnCode == 0 {
				returnCode = http.StatusOK
			}
		}
	}

	return res, returnCode
}

func recoverFromPanic(u *UHTTP, handlerProcessed chan interface{}, r *http.Request, returnCode *int) {
	if rec := recover(); rec != nil {
		err := fmt.Errorf("panic: handlerExecution (%s)", rec)
		u.opts.log.Sugar().Errorf("panic [path: %s] %s", r.RequestURI, err)
		stack := debug.Stack()
		uhelpers.CallForByteArrayLineByLine(stack, u.opts.log.Sugar().Errorf, fmt.Sprintf("panic [path: %s] ", r.RequestURI))
		err = fmt.Errorf("%s stackTrace: %s", err, strings.ReplaceAll(string(stack), "\n", "\\n"))
		*returnCode = http.StatusInternalServerError

		// let caller know if a panic happened
		// do this asynchronously
		if u.opts.handleHandlerPanics != nil {
			for _, fn := range u.opts.handleHandlerPanics {
				go func(fn func(r *http.Request, err error)) {
					fn(r, err)
				}(fn)
			}
		}

		if u.opts.sendPanicInfoToClient {
			handlerProcessed <- err
			return
		}
		handlerProcessed <- fmt.Errorf("internal server error")
	}
}
