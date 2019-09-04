package uhttp

import (
	"encoding/json"
	"log"
	"net/http"
	"runtime"
)

func RenderError(w http.ResponseWriter, r *http.Request, err error) {
	if err != nil {
		js, _ := json.Marshal(Error{
			Error: err.Error(),
		})

		_, fn, line, _ := runtime.Caller(1)
		w.WriteHeader(http.StatusBadRequest)
		CheckAndLogErrorSecondArg(w.Write(js))
		if customLog != nil {
			customLog.Errorf("Error in %s (%s:%d): %s", r.RequestURI, fn, line, err.Error())
		} else {
			log.Printf("Error in %s (%s:%d): %s", r.RequestURI, fn, line, err.Error())
		}
	} else {
		log.Panic("Error is nil, and trying to RenderError")
	}
}

func RenderMessageWithStatusCode(w http.ResponseWriter, r *http.Request, code int, msg string) {
	myMap := map[string]string{"msg": msg}
	js, _ := json.Marshal(myMap)
	w.WriteHeader(code)
	CheckAndLogErrorSecondArg(w.Write(js))

	_, fn, line, _ := runtime.Caller(1)
	if customLog != nil {
		customLog.Errorf("Msg in %s (%s:%d): %s", r.RequestURI, fn, line, msg)
	} else {
		log.Printf("Msg in %s (%s:%d): %s", r.RequestURI, fn, line, msg)
	}
}
