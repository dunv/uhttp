package main

import (
	"io"
	"net/http"
	"os"

	"github.com/dunv/uhttp"
	"github.com/dunv/ulog"
)

func main() {

	u := uhttp.NewUHTTP(
		uhttp.WithSendPanicInfoToClient(true),
	)

	u.Handle("/", uhttp.NewHandler(uhttp.WithGet(func(r *http.Request, ret *int) interface{} {
		return map[string]string{"hello": "world"}
	})))

	// force a handler-panic
	u.Handle("/forcePanic", uhttp.NewHandler(uhttp.WithGet(func(r *http.Request, ret *int) interface{} {
		var test interface{}
		test = 5
		wrongType := test.(string)
		return wrongType
	})))

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
