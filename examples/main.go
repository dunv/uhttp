package main

import (
	"net/http"

	"github.com/dunv/uhttp"
	"github.com/dunv/ulog"
)

func main() {

	u := uhttp.NewUHTTP()

	u.Handle("/", uhttp.NewHandler(uhttp.WithGet(func(r *http.Request, ret *int) interface{} {
		return map[string]string{"hello": "world"}
	})))

	u.Handle("/test", uhttp.NewHandler(uhttp.WithGet(func(r *http.Request, ret *int) interface{} {
		writer := r.Context().Value(uhttp.CtxKeyResponseWriter).(http.ResponseWriter)
		writer.WriteHeader(http.StatusAccepted)

		ulog.LogIfErrorSecondArg(writer.Write([]byte(`{"nothing":"toSay"}` + "\n")))
		return nil
	})))

	ulog.Fatal(u.ListenAndServe())
}
