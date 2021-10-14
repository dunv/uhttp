package uhttp

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
)

type cache struct {
	mu      *sync.RWMutex
	pattern string
	maxAge  time.Duration
	_data   map[string]cacheEntry
}

func (c cache) Set(key string, data []byte, statusCode int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	dataCopy := make([]byte, len(data))
	copy(dataCopy, data)
	c._data[key] = cacheEntry{
		updatedOn:  time.Now(),
		data:       dataCopy,
		statusCode: statusCode,
	}
}

func (c cache) Keys() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	keys := []string{}
	for key := range c._data {
		keys = append(keys, key)
	}
	return keys
}

func (c cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c._data, key)
}

func (c cache) Get(key string) (cacheEntry, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if entry, ok := c._data[key]; ok {
		val := entry.Clone()
		return val, ok
	}
	return cacheEntry{}, false
}

func (c cache) Size() uint64 {
	c.mu.RLock()
	defer c.mu.RUnlock()

	total := uint64(0)
	for _, entry := range c._data {
		total += uint64(len(entry.data))
	}
	return total
}

func (c cache) Count() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c._data)
}

type cacheEntry struct {
	updatedOn  time.Time
	data       json.RawMessage
	statusCode int
}

func (e cacheEntry) Clone() cacheEntry {
	dataCopy := make([]byte, len(e.data))
	copy(dataCopy, e.data)
	return cacheEntry{
		updatedOn:  e.updatedOn,
		data:       dataCopy,
		statusCode: e.statusCode,
	}
}

func (c cacheEntry) String() string {
	return fmt.Sprintf("{%s - %d - '%s'}", c.updatedOn.Format(time.RFC3339Nano), c.statusCode, string(c.data))
}

func cacheHash(opts handlerOptions, r *http.Request) (string, error) {
	// Request-Body
	body := ""
	if r.Body != nil {
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return "", err
		}
		defer r.Body.Close()
		if r.Body != nil {
			r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		}
		body = string(bodyBytes)
	}

	// Request-Params
	params := r.URL.RawQuery

	// Relevant headers
	headers := []string{}
	for _, h := range opts.CacheRelevantHeaders {
		if val := r.Header.Get(h); val != "" {
			headers = append(headers, fmt.Sprintf("%s: %s", h, val))
		}
	}
	uniqueString := fmt.Sprintf("%s-%s-%s", body, params, strings.Join(headers, "-"))
	s := md5.Sum([]byte(uniqueString))
	return string(s[:]), nil
}
