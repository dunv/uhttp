package main

import (
	"errors"
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
		uhttp.WithGranularLogging(true, true, true, true),
	)
	u.ExposeCacheHandlers()

	u.Handle("/hello", uhttp.NewHandler(uhttp.WithGet(func(r *http.Request, ret *int) interface{} {
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

	counter1 := 0
	u.Handle("/testCache1", uhttp.NewHandler(
		uhttp.WithCachePersistEncodings(),
		uhttp.WithCache(10*time.Minute),
		uhttp.WithGet(func(r *http.Request, ret *int) interface{} {
			counter1++
			return map[string]int{"counter1": counter1}
		}),
	))

	counter2 := 0
	u.Handle("/testCache2", uhttp.NewHandler(
		uhttp.WithCache(10*time.Minute),
		uhttp.WithGet(func(r *http.Request, ret *int) interface{} {
			counter2++
			return map[string]int{"counter2": counter2}
		}),
	))

	// 	u.Handle("/testCacheAutomatic", uhttp.NewHandler(
	// 		uhttp.WithCache(10*time.Minute),
	// 		uhttp.WithAutomaticCacheUpdates(5*time.Second),
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

		f, err := os.OpenFile("./testDatei", os.O_RDONLY, 0644)
		if err != nil {
			return err
		}
		defer f.Close()

		writer.Header().Set("Content-Type", "application/zip")
		writer.WriteHeader(http.StatusAccepted)
		ulog.LogIfErrorSecondArg(io.Copy(writer, f))

		return nil
	})))

	ulog.FatalIfError(u.RegisterStaticFilesHandler("static"))

	ulog.Fatal(u.ListenAndServe())
}
