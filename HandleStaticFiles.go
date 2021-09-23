package uhttp

import (
	"bytes"
	"compress/gzip"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/itchio/go-brotli/enc"

	"github.com/dunv/uhelpers"
	"github.com/dunv/ulog"
)

type cachedFile struct {
	Content        []byte
	GzippedContent []byte
	BrContent      []byte
	ContentType    string
}

var filesCache = map[string]cachedFile{}

// static files handler which only works if initialized with "RegisterStaticFilesHandler"
// (only serves from initialized cache)
func staticFilesHandler(u *UHTTP) http.HandlerFunc {
	return chain(addLoggingMiddleware(u))(func(w http.ResponseWriter, r *http.Request) {
		if len(filesCache) == 0 {
			u.RenderError(w, r, errors.New("staticFilesHandler used but not initialized"))
			return
		}

		var cachedFile cachedFile
		var ok bool

		// Find file (fallback to index.html)
		if cachedFile, ok = filesCache[r.URL.Path]; !ok {
			cachedFile = filesCache["/index.html"]
		}
		w.Header().Add("Content-Type", cachedFile.ContentType)

		// If client accepts br or gzip -> return compressed
		if acceptEncoding := r.Header.Get("Accept-Encoding"); strings.Contains(acceptEncoding, "br") && u.opts.enableBrotli {
			w.Header().Add("Content-Encoding", "br")
			ulog.LogIfErrorSecondArg(w.Write(cachedFile.BrContent))
		} else if acceptEncoding := r.Header.Get("Accept-Encoding"); strings.Contains(acceptEncoding, "gzip") && u.opts.enableGzip {
			w.Header().Add("Content-Encoding", "gzip")
			ulog.LogIfErrorSecondArg(w.Write(cachedFile.GzippedContent))
		} else {
			ulog.LogIfErrorSecondArg(w.Write(cachedFile.Content))
		}
	})
}

// RegisterStaticFilesHandler which serves content from a directory and
// redirects all requests to non-existant files to index.html
// index.html must be present!
// - read all files from root directory
// - create cache for these files containing original, gzip, br
// - register handlers for main http-mux
func (u *UHTTP) RegisterStaticFilesHandler(root string) error {
	fileNames := []string{}
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if path != root && !info.IsDir() {
			fileNames = append(fileNames, path)
		}
		return nil
	})
	if err != nil {
		return err
	}

	foundMainFile := false
	for _, fileName := range fileNames {
		if strings.Contains(fileName, ".DS_Store") {
			ulog.Infof("Skipping '%s'", fileName)
			continue
		}
		fileContent, err := ioutil.ReadFile(fileName)
		if err != nil {
			return err
		}

		// Strip the root part
		var pattern string
		if strings.HasPrefix(root, "./") {
			pattern = strings.ReplaceAll(fileName, filepath.Base(root), "")
		} else {
			pattern = strings.ReplaceAll(fileName, root, "")
		}

		if strings.ReplaceAll(pattern, "/", "") == "index.html" {
			foundMainFile = true
		}

		// fmt.Println("registering", fileName, pattern, strings.ReplaceAll(pattern, "/", ""))

		// Detect content-type automatically
		contentType := http.DetectContentType(fileContent)

		// For some reason it cannot detect minified js and css files -> manual override
		if strings.HasSuffix(pattern, ".js") {
			contentType = "text/javascript; charset=utf-8"
		} else if strings.HasSuffix(pattern, ".css") {
			contentType = "text/css; charset=utf-8"
		} else if strings.HasSuffix(pattern, ".html") {
			contentType = "text/html; charset=utf-8"
		}

		cached := cachedFile{
			Content:     fileContent,
			ContentType: contentType,
		}

		if u.opts.enableGzip {
			// Compress gzip
			var gzipBuffer bytes.Buffer
			gzipWriter, err := gzip.NewWriterLevel(&gzipBuffer, gzip.BestCompression)
			if err != nil {
				return err
			}
			_, err = gzipWriter.Write(fileContent)
			if err != nil {
				return err
			}
			err = gzipWriter.Close()
			if err != nil {
				return err
			}
			cached.GzippedContent, err = ioutil.ReadAll(&gzipBuffer)
			if err != nil {
				return err
			}
		}

		if u.opts.enableBrotli {
			// Compress brotli
			var brotliBuffer bytes.Buffer
			brotliWriter := enc.NewBrotliWriter(&brotliBuffer, &enc.BrotliWriterOptions{Quality: 11})
			_, err = brotliWriter.Write(fileContent)
			if err != nil {
				return err
			}
			err = brotliWriter.Close()
			if err != nil {
				return err
			}
			cached.BrContent, err = ioutil.ReadAll(&brotliBuffer)
			if err != nil {
				return err
			}
		}

		filesCache[pattern] = cached
		if !u.opts.silentStaticFileRegistration {
			ulog.Infof("Registered http static %s (%s, gzip:%s, br:%s)",
				pattern,
				uhelpers.ByteCountIEC(int64(len(fileContent))),
				uhelpers.ByteCountIEC(int64(len(cached.GzippedContent))),
				uhelpers.ByteCountIEC(int64(len(cached.BrContent))),
			)
		}
		u.opts.serveMux.HandleFunc(pattern, staticFilesHandler(u))
	}
	ulog.Infof("Registered http static / -> /index.html")
	u.opts.serveMux.HandleFunc("/", staticFilesHandler(u))

	if !foundMainFile {
		return errors.New("could not find index.html")
	}

	return nil
}
