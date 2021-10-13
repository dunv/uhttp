package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/dunv/uhttp"
	"github.com/dunv/ulog"
)

func main() {

	u := uhttp.NewUHTTP(
		uhttp.WithSendPanicInfoToClient(true),
		uhttp.WithExposeCacheHandlers(true),
		uhttp.WithGranularLogging(false, true, true),
	)

	u.Handle("/", uhttp.NewHandler(uhttp.WithGet(func(r *http.Request, ret *int) interface{} {
		return map[string]string{"hello": "world"}
	})))

	// force a handler-panic
	u.Handle("/forcePanic", uhttp.NewHandler(uhttp.WithGet(func(r *http.Request, ret *int) interface{} {
		var test interface{} = 5
		wrongType := test.(string)
		return wrongType
	})))

	u.Handle("/forceError", uhttp.NewHandler(uhttp.WithGet(func(r *http.Request, ret *int) interface{} {
		return errors.New("this is an error")
	})))

	u.Handle("/testCache", uhttp.NewHandler(
		uhttp.WithCache(10*time.Minute),
		uhttp.WithGet(func(r *http.Request, ret *int) interface{} {
			fmt.Println("executing GET testCacheHandler")
			return map[string]string{
				"method":    "get",
				"updatedOn": time.Now().Format(time.RFC3339Nano),
			}
		}),
	))

	// 	u.Handle("/testCacheAutomatic", uhttp.NewHandler(
	// 		uhttp.WithCache(10*time.Minute),
	// 		uhttp.WithAutomaticCacheUpdates(5*time.Second, []string{}),
	// 		uhttp.WithGet(func(r *http.Request, ret *int) interface{} {
	// 			fmt.Println("executing GET testCacheAutomaticHandler")
	// 			return map[string]string{
	// 				"method":    "get",
	// 				"updatedOn": time.Now().Format(time.RFC3339Nano),
	// 			}
	// 		}),
	// 		uhttp.WithPost(func(r *http.Request, ret *int) interface{} {
	// 			fmt.Println("executing POST testCacheAutomaticHandler")
	// 			return map[string]string{
	// 				"method":    "post",
	// 				"updatedOn": time.Now().Format(time.RFC3339Nano),
	// 			}
	// 		}),
	// 	))

	u.Handle("/test", uhttp.NewHandler(uhttp.WithGet(func(r *http.Request, ret *int) interface{} {
		writer := r.Context().Value(uhttp.CtxKeyResponseWriter).(http.ResponseWriter)
		writer.WriteHeader(http.StatusAccepted)

		ulog.LogIfErrorSecondArg(writer.Write([]byte(`{"nothing":"toSay"}` + "\n")))
		return nil
	})))

	u.Handle("/testDatei", uhttp.NewHandler(uhttp.WithGet(func(r *http.Request, ret *int) interface{} {
		writer := r.Context().Value(uhttp.CtxKeyResponseWriter).(http.ResponseWriter)
		writer.WriteHeader(http.StatusAccepted)

		f, err := os.OpenFile("./testDatei", os.O_RDONLY, 0644)
		if err != nil {
			return err
		}
		defer f.Close()

		ulog.LogIfErrorSecondArg(io.Copy(writer, f))
		return nil
	})))

	ulog.Fatal(u.ListenAndServe())
}
