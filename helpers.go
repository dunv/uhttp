package uhttp

import (
	"bytes"
	"io"
	"net/http"
	"runtime"
	"strings"
)

func ExtractAndRestoreRequestBody(r *http.Request) []byte {
	if r.Body != nil {
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			return []byte{}
		}
		defer r.Body.Close()
		if r.Body != nil {
			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}
		return bodyBytes
	}
	return []byte{}
}

// figures out the first caller of the function outside of github.com/dunv/http AND net/http
// straight out of net/http/server.go
func relevantCaller() runtime.Frame {
	pc := make([]uintptr, 16)
	n := runtime.Callers(1, pc)
	frames := runtime.CallersFrames(pc[:n])
	var frame runtime.Frame
	for {
		frame, more := frames.Next()
		if !strings.HasPrefix(frame.Function, "github.com/dunv/uhttp") && !strings.HasPrefix(frame.Function, "net/http") {
			return frame
		}
		if !more {
			break
		}
	}
	return frame
}
