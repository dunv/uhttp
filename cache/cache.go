package cache

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

func NewCache(maxAge time.Duration) *Cache {
	return &Cache{
		mu:     &sync.RWMutex{},
		maxAge: maxAge,
		_data:  map[string]cacheEntry{},
	}
}

type Cache struct {
	mu     *sync.RWMutex
	maxAge time.Duration
	_data  map[string]cacheEntry
}

func (c Cache) MaxAge() time.Duration {
	return c.maxAge
}

func (c Cache) Set(
	requestBody []byte,
	requestParams string,
	requestHeader http.Header,
	responseBody []byte,
	responseHeader http.Header,
	responseStatuCode int,
) {
	dataCopy := make([]byte, len(responseBody))
	copy(dataCopy, responseBody)
	key := hash(requestBody, requestParams, requestHeader)

	c.mu.Lock()
	defer c.mu.Unlock()
	c._data[key] = cacheEntry{
		updatedOn:          time.Now(),
		responseBody:       responseBody,
		responseHeader:     responseHeader,
		responseStatusCode: responseStatuCode,
	}
}

func (c Cache) Keys() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	keys := []string{}
	for key := range c._data {
		keys = append(keys, key)
	}
	return keys
}

func (c Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c._data, key)
}

func (c Cache) Get(requestBody []byte, requestParams string, requestHeader http.Header) (cacheEntry, bool, string) {
	key := hash(requestBody, requestParams, requestHeader)

	c.mu.RLock()
	defer c.mu.RUnlock()
	if entry, ok := c._data[key]; ok {
		val := entry.Clone()
		return val, ok, key
	}
	return cacheEntry{}, false, ""
}

func (c Cache) GetByKey(key string) (cacheEntry, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if entry, ok := c._data[key]; ok {
		val := entry.Clone()
		return val, ok
	}
	return cacheEntry{}, false
}

func (c Cache) Size() uint64 {
	c.mu.RLock()
	defer c.mu.RUnlock()

	total := uint64(0)
	for _, entry := range c._data {
		total += uint64(len(entry.responseBody))
	}
	return total
}

func (c Cache) Count() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c._data)
}

type cacheEntry struct {
	updatedOn          time.Time
	responseBody       []byte
	responseHeader     http.Header
	responseStatusCode int
}

func (e cacheEntry) Clone() cacheEntry {
	dataCopy := make([]byte, len(e.responseBody))
	copy(dataCopy, e.responseBody)
	return cacheEntry{
		updatedOn:          e.updatedOn,
		responseBody:       dataCopy,
		responseHeader:     e.responseHeader.Clone(),
		responseStatusCode: e.responseStatusCode,
	}
}

func (e cacheEntry) UpdatedOn() time.Time {
	return e.updatedOn
}

func (e cacheEntry) Write(w http.ResponseWriter) error {
	for k, v := range e.responseHeader {
		w.Header().Set(k, strings.Join(v, ", "))
	}
	w.WriteHeader(e.responseStatusCode)
	_, err := w.Write(e.responseBody)
	if err != nil {
		return err
	}
	return nil
}

func (c cacheEntry) String() string {
	return fmt.Sprintf("{%s - %d - '%s'}", c.updatedOn.Format(time.RFC3339Nano), c.responseStatusCode, string(c.responseBody))
}

func hash(body []byte, params string, header http.Header) string {
	headers := []string{}
	for k, h := range header {
		if k == "Accept-Encoding" {
			headers = append(headers, strings.Join(h, ", "))
		}
	}
	uniqueString := fmt.Sprintf("%s-%s-%s", body, params, strings.Join(headers, "-"))
	s := md5.Sum([]byte(uniqueString))
	return string(s[:])
}
