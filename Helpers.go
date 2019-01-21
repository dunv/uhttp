package uhttp

import (
	"encoding/json"
	"log"
	"net/http"
	"runtime"
)

// RenderError in json
func RenderError(w http.ResponseWriter, r *http.Request, err error) {
	if err != nil {
		js, _ := json.Marshal(Error{
			Error: err.Error(),
		})

		_, fn, line, _ := runtime.Caller(1)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(js)
		log.Printf("Error in %s (%s:%d): %s", r.RequestURI, fn, line, err.Error())
	} else {
		log.Panic("Error is nil, and trying to RenderError")
	}
}

// RenderMessageWithStatusCode helper
func RenderMessageWithStatusCode(w http.ResponseWriter, r *http.Request, code int, msg string) {
	myMap := map[string]string{"msg": msg}
	js, _ := json.Marshal(myMap)
	w.WriteHeader(code)
	w.Write(js)

	_, fn, line, _ := runtime.Caller(1)
	log.Printf("Msg in %s (%s:%d): %s", r.RequestURI, fn, line, msg)
}
