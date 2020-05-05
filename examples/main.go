package main

import (
	"net/http"

	"github.com/dunv/uhttp"
	"github.com/dunv/ulog"
)

func main() {

	u := uhttp.NewUHTTP()

	u.Handle("/", uhttp.WithGet(func(r *http.Request, ret *int) interface{} {
		return map[string]string{"hello": "world"}
	}))

	ulog.Fatal(u.ListenAndServe())
}
