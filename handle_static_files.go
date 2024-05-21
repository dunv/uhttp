package uhttp

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/andybalholm/brotli"
	"github.com/dunv/uhelpers"
	"github.com/klauspost/compress/flate"
	"github.com/klauspost/compress/gzip"
)

type cachedFile struct {
	Content        []byte
	GzippedContent []byte
	BrContent      []byte
	DeflateContent []byte
	ContentType    string
}

var filesCache = map[string]cachedFile{}

// static files handler which only works if initialized with "RegisterStaticFilesHandler"
// (only serves from initialized cache)
func StaticFilesHandler(u *UHTTP) http.HandlerFunc {
	return chain(addLoggingMiddleware(u, nil, true))(func(w http.ResponseWriter, r *http.Request) {
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
			if _, err := w.Write(cachedFile.BrContent); err != nil {
				u.Log().Errorf("%s", err)
			}
		} else if acceptEncoding := r.Header.Get("Accept-Encoding"); strings.Contains(acceptEncoding, "gzip") && u.opts.enableGzip {
			w.Header().Add("Content-Encoding", "gzip")
			if _, err := w.Write(cachedFile.GzippedContent); err != nil {
				u.Log().Errorf("%s", err)
			}
		} else if acceptEncoding := r.Header.Get("Accept-Encoding"); strings.Contains(acceptEncoding, "deflate") && u.opts.enableDeflate {
			w.Header().Add("Content-Encoding", "deflate")
			if _, err := w.Write(cachedFile.DeflateContent); err != nil {
				u.Log().Errorf("%s", err)
			}
		} else {
			if _, err := w.Write(cachedFile.Content); err != nil {
				u.Log().Errorf("%s", err)
			}
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
		if path != root && !info.IsDir() && !strings.HasSuffix(path, ".brotli") && !strings.HasSuffix(path, ".gz") && !strings.HasSuffix(path, ".deflate") {
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
			u.Log().Infof("Skipping '%s'", fileName)
			continue
		}
		fileContent, err := os.ReadFile(fileName)
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

		// Cannot identify the following from the file itself
		// Overrides (taken from)
		// https://wiki.selfhtml.org/wiki/MIME-Type/%C3%9Cbersicht
		if strings.HasSuffix(pattern, ".js") {
			contentType = "text/javascript; charset=utf-8"
		} else if strings.HasSuffix(pattern, ".css") {
			contentType = "text/css; charset=utf-8"
		} else if strings.HasSuffix(pattern, ".html") {
			contentType = "text/html; charset=utf-8"
		} else if strings.HasSuffix(pattern, ".svg") {
			contentType = "image/svg+xml; charset=utf-8"
		}

		cached := cachedFile{
			Content:     fileContent,
			ContentType: contentType,
		}

		if u.opts.enableGzip {
			if _, err := os.Stat(fmt.Sprintf("%s.gz", fileName)); err == nil {
				u.Log().Infof("http static: gzip-compressed file %s exists, not compressing again", pattern)
				cached.GzippedContent, err = os.ReadFile(fmt.Sprintf("%s.gz", fileName))
				if err != nil {
					return err
				}
			} else {
				var buffer bytes.Buffer
				u.Log().Infof("http static: gzip-compressed file %s does not exist. Compressing", pattern)
				gzipWriter, err := gzip.NewWriterLevel(&buffer, u.opts.gzipCompressionLevel)
				if err != nil {
					return err
				}
				if _, err = gzipWriter.Write(fileContent); err != nil {
					return err
				}
				if err := gzipWriter.Flush(); err != nil {
					return err
				}
				if err := gzipWriter.Close(); err != nil {
					return err
				}
				cached.GzippedContent = buffer.Bytes()
			}
		}

		if u.opts.enableBrotli {
			if _, err := os.Stat(fmt.Sprintf("%s.brotli", fileName)); err == nil {
				u.Log().Infof("http static: brotli-compressed file %s exists, not compressing again", pattern)
				cached.BrContent, err = os.ReadFile(fmt.Sprintf("%s.brotli", fileName))
				if err != nil {
					return err
				}
			} else {
				var buffer bytes.Buffer
				u.Log().Infof("http static: brotli-compressed file %s does not exist. Compressing", pattern)
				brotliWriter := brotli.NewWriterLevel(&buffer, u.opts.brotliCompressionLevel)
				if _, err = brotliWriter.Write(fileContent); err != nil {
					return err
				}
				if err := brotliWriter.Flush(); err != nil {
					return err
				}
				if err := brotliWriter.Close(); err != nil {
					return err
				}
				cached.BrContent = buffer.Bytes()
			}
		}

		if u.opts.enableDeflate {
			if _, err := os.Stat(fmt.Sprintf("%s.deflate", fileName)); err == nil {
				u.Log().Infof("http static: deflate-compressed file %s exists, not compressing again", pattern)
				cached.DeflateContent, err = os.ReadFile(fmt.Sprintf("%s.deflate", fileName))
				if err != nil {
					return err
				}
			} else {
				var buffer bytes.Buffer
				u.Log().Infof("http static: deflate-compressed file %s does not exist. Compressing", pattern)
				deflateWriter, err := flate.NewWriter(&buffer, u.opts.deflateCompressionLevel)
				if err != nil {
					return err
				}
				if _, err = deflateWriter.Write(fileContent); err != nil {
					return err
				}
				if err := deflateWriter.Flush(); err != nil {
					return err
				}
				if err := deflateWriter.Close(); err != nil {
					return err
				}
				cached.DeflateContent = buffer.Bytes()
			}
		}

		filesCache[pattern] = cached
		if !u.opts.silentStaticFileRegistration {
			u.Log().Infof("Registered http static %s (%s, gzip:%s, br:%s, deflate:%s)",
				pattern,
				uhelpers.FormatByteCountIEC(int64(len(fileContent))),
				uhelpers.FormatByteCountIEC(int64(len(cached.GzippedContent))),
				uhelpers.FormatByteCountIEC(int64(len(cached.BrContent))),
				uhelpers.FormatByteCountIEC(int64(len(cached.DeflateContent))),
			)
		}
		u.opts.serveMux.HandleFunc(pattern, StaticFilesHandler(u))
	}
	u.Log().Infof("Registered http static / -> /index.html")
	u.opts.serveMux.HandleFunc("/", StaticFilesHandler(u))

	if !foundMainFile {
		return errors.New("could not find index.html")
	}

	return nil
}
