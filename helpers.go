package uhttp

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

func ExtractAndRestoreRequestBody(r *http.Request) []byte {
	if r.Body != nil {
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return []byte{}
		}
		defer r.Body.Close()
		if r.Body != nil {
			r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		}
		return bodyBytes
	}
	return []byte{}
}
