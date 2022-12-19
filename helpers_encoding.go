package uhttp

import (
	"bytes"
	"io"
	"net/http"
	"strings"

	"github.com/andybalholm/brotli"
	"github.com/klauspost/compress/flate"
	"github.com/klauspost/compress/gzip"
)

const (
	HEADER_ACCEPT_ENCODING  = "Accept-Encoding"
	HEADER_CONTENT_ENCODING = "Content-Encoding"
	ENCODING_PLAIN          = ""
	ENCODING_BROTLI         = "br"
	ENCODING_GZIP           = "gzip"
	ENCODING_DEFLATE        = "deflate"
)

type nopCloser struct{ internal io.Writer }

func (n nopCloser) Write(p []byte) (int, error) {
	return n.internal.Write(p)
}
func (n nopCloser) Close() error {
	return nil
}

func (u *UHTTP) encodingWriter(w io.Writer, encoding string) io.WriteCloser {
	switch encoding {
	case ENCODING_BROTLI:
		return brotli.NewWriterLevel(w, u.opts.brotliCompressionLevel)
	case ENCODING_GZIP:
		// we check that we are using a supported level when assigning the option
		ww, _ := gzip.NewWriterLevel(w, u.opts.gzipCompressionLevel)
		return ww
	case ENCODING_DEFLATE:
		// we check that we are using a supported level when assigning the option
		ww, _ := flate.NewWriter(w, u.opts.deflateCompressionLevel)
		return ww
	default:
		return nopCloser{w}
	}
}

func (u *UHTTP) determineEncoding(r *http.Request, statusCode int) string {
	// The go-http-client implementation decodes gzip out-of-the-box, but only if it gets 200 OK
	// For now: use the same behavior here
	acceptEncoding := r.Header.Get(HEADER_ACCEPT_ENCODING)
	if statusCode == http.StatusOK {
		if u.opts.enableBrotli && strings.Contains(acceptEncoding, ENCODING_BROTLI) {
			return ENCODING_BROTLI
		} else if u.opts.enableGzip && strings.Contains(acceptEncoding, ENCODING_GZIP) {
			return ENCODING_GZIP
		} else if u.opts.enableDeflate && strings.Contains(acceptEncoding, ENCODING_DEFLATE) {
			return ENCODING_DEFLATE
		}
	}
	return ENCODING_PLAIN
}

func (u *UHTTP) compressJSON(encoding string, data []byte) ([]byte, error) {
	var b bytes.Buffer
	ew := u.encodingWriter(&b, encoding)
	_, err := ew.Write(data)
	if err != nil {
		return nil, err
	}
	err = ew.Close()
	if err != nil {
		return nil, err
	}
	encoded, err := io.ReadAll(&b)
	if err != nil {
		return nil, err
	}
	return encoded, nil
}
