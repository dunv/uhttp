package uhttp

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
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

func GzipDecodeRequestBody(r *http.Request, model interface{}) error {
	reader, err := ReaderHelper(r.Header, r.Body)
	if err != nil {
		return fmt.Errorf("err parsing request (err getting reader %s)", err)
	}

	err = json.NewDecoder(reader).Decode(model)
	if err != nil {
		return fmt.Errorf("err parsing request (err decoding %s)", err)
	}
	defer r.Body.Close()

	return nil
}

func GzipDecodeResponseBody(res *http.Response) ([]byte, error) {
	wrappedReader, err := ReaderHelper(res.Header, res.Body)
	if err != nil {
		return nil, err
	}

	decodedResponse, err := ioutil.ReadAll(wrappedReader)
	if err != nil {
		return nil, err
	}
	res.Body.Close()

	return decodedResponse, nil
}

func GzipEncodeRequestBody(body []byte) (io.ReadCloser, error) {
	buffer := bytes.NewBuffer([]byte{})
	writer, err := gzip.NewWriterLevel(buffer, 5)
	if err != nil {
		return nil, fmt.Errorf("could not initialize gzip writer (%s)", err)
	}
	_, err = writer.Write(body)
	if err != nil {
		return nil, fmt.Errorf("could not write to gzip writer (%s)", err)
	}
	err = writer.Flush()
	if err != nil {
		return nil, fmt.Errorf("could not flush gzip writer (%s)", err)
	}
	return ioutil.NopCloser(bytes.NewReader(buffer.Bytes())), nil
}
