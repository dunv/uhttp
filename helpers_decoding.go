package uhttp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/andybalholm/brotli"
	"github.com/klauspost/compress/flate"
	"github.com/klauspost/compress/gzip"
)

// Wraps a reader in the correct decoder based on http-headers
// inspired by https://medium.com/axiomzenteam/put-your-http-requests-on-a-diet-3e1e52333014
func DecodingReader(header http.Header, body io.ReadCloser) (io.ReadCloser, error) {
	var reader io.ReadCloser
	switch header.Get("Content-Encoding") {
	case "gzip":
		gz, err := gzip.NewReader(body)
		if err != nil {
			return nil, fmt.Errorf("could not decode gzipped response (%s)", err)
		}
		reader = gz
	case "br":
		reader = io.NopCloser(brotli.NewReader(body))
	case "deflate":
		reader = flate.NewReader(body)
	default:
		reader = body
	}

	return reader, nil
}

func decodeRequestBody(r *http.Request, model interface{}) error {
	reader, err := DecodingReader(r.Header, r.Body)
	if err != nil {
		return fmt.Errorf("err parsing request (err getting reader %s)", err)
	}

	err = json.NewDecoder(reader).Decode(model)
	if err != nil {
		return fmt.Errorf("err parsing request (err decoding %s)", err)
	}
	defer r.Body.Close()
	defer reader.Close()

	return nil
}

func decodeResponseBody(res *http.Response) ([]byte, error) {
	wrappedReader, err := DecodingReader(res.Header, res.Body)
	if err != nil {
		return nil, err
	}

	decodedResponse, err := io.ReadAll(wrappedReader)
	if err != nil {
		return nil, err
	}
	res.Body.Close()

	return decodedResponse, nil
}

func gzipEncodeRequestBody(body []byte) (io.ReadCloser, error) {
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
	return io.NopCloser(bytes.NewReader(buffer.Bytes())), nil
}
