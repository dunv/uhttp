package uhttp

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
)

// inspired by https://medium.com/axiomzenteam/put-your-http-requests-on-a-diet-3e1e52333014
func ReaderHelper(header http.Header, body io.ReadCloser) (io.Reader, error) {
	var reader io.Reader
	switch header.Get("Content-Encoding") {
	case "gzip":
		gz, err := gzip.NewReader(body)
		if err != nil {
			return nil, fmt.Errorf("could not decode gzipped response (%s)", err)
		}
		defer gz.Close()
		reader = gz
	default:
		reader = body
	}

	return reader, nil
}
