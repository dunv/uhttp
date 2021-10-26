package cache

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"sync"
	"time"
	"unsafe"
)

func NewCache(
	maxAge time.Duration,
) *Cache {
	return &Cache{
		mu:     &sync.RWMutex{},
		maxAge: maxAge,
		data:   map[string]CacheEntry{},
	}
}

type Cache struct {
	mu     *sync.RWMutex
	maxAge time.Duration
	data   map[string]CacheEntry
}

func (c Cache) MaxAge() time.Duration {
	return c.maxAge
}

func (c Cache) Set(
	requestBody []byte,
	requestParams string,
	requestHeader http.Header,
	responseModel interface{},
	responseHeader http.Header,
	responseStatusCode int,
	responseBodyPlain []byte,
	responseBodyBrotli []byte,
	responseBodyGzip []byte,
	responseBodyDeflate []byte,
) {
	key := hash(requestBody, requestParams)

	c.mu.Lock()
	defer c.mu.Unlock()
	e := CacheEntry{
		updatedOn:           time.Now(),
		responseModel:       responseModel,
		responseHeader:      responseHeader,
		responseStatusCode:  responseStatusCode,
		responseBodyPlain:   responseBodyPlain,
		responseBodyBrotli:  responseBodyBrotli,
		responseBodyGzip:    responseBodyGzip,
		responseBodyDeflate: responseBodyDeflate,
	}

	c.data[key] = e
}

func (c Cache) Keys() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	keys := []string{}
	for key := range c.data {
		keys = append(keys, key)
	}
	return keys
}

func (c Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.data, key)
}

func (c Cache) Get(requestBody []byte, requestParams string) (CacheEntry, bool, string) {
	key := hash(requestBody, requestParams)

	c.mu.RLock()
	defer c.mu.RUnlock()
	if entry, ok := c.data[key]; ok {
		val := entry.Clone()
		return val, ok, key
	}
	return CacheEntry{}, false, ""
}

func (c Cache) GetByKey(key string) (CacheEntry, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if entry, ok := c.data[key]; ok {
		val := entry.Clone()
		return val, ok
	}
	return CacheEntry{}, false
}

func (c Cache) Size() uint64 {
	c.mu.RLock()
	defer c.mu.RUnlock()

	total := uint64(0)
	for _, entry := range c.data {
		total += uint64(len(entry.responseBodyPlain))
		total += uint64(len(entry.responseBodyBrotli))
		total += uint64(len(entry.responseBodyGzip))
		total += uint64(len(entry.responseBodyDeflate))
		total += uint64(unsafe.Sizeof(entry.responseModel))
	}
	return total
}

func (c Cache) Count() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.data)
}

type CacheEntry struct {
	updatedOn           time.Time
	responseModel       interface{}
	responseBodyPlain   []byte
	responseBodyGzip    []byte
	responseBodyBrotli  []byte
	responseBodyDeflate []byte
	responseHeader      http.Header
	responseStatusCode  int
}

func (e *CacheEntry) Clone() CacheEntry {
	var responseBodyPlainCopy []byte
	if e.responseBodyPlain != nil {
		responseBodyPlainCopy = make([]byte, len(e.responseBodyPlain))
		copy(responseBodyPlainCopy, e.responseBodyPlain)
	}

	var responseBodyGzipCopy []byte
	if e.responseBodyGzip != nil {
		responseBodyGzipCopy = make([]byte, len(e.responseBodyGzip))
		copy(responseBodyGzipCopy, e.responseBodyGzip)
	}

	var responseBodyBrotliCopy []byte
	if e.responseBodyBrotli != nil {
		responseBodyBrotliCopy = make([]byte, len(e.responseBodyBrotli))
		copy(responseBodyBrotliCopy, e.responseBodyBrotli)
	}

	var responseBodyDeflateCopy []byte
	if e.responseBodyDeflate != nil {
		responseBodyDeflateCopy = make([]byte, len(e.responseBodyDeflate))
		copy(responseBodyDeflateCopy, e.responseBodyDeflate)
	}

	return CacheEntry{
		updatedOn:           e.updatedOn,
		responseModel:       e.responseModel,
		responseBodyPlain:   responseBodyPlainCopy,
		responseBodyGzip:    responseBodyGzipCopy,
		responseBodyBrotli:  responseBodyBrotliCopy,
		responseBodyDeflate: responseBodyDeflateCopy,
		responseHeader:      e.responseHeader.Clone(),
		responseStatusCode:  e.responseStatusCode,
	}
}

func (e *CacheEntry) UpdatedOn() time.Time {
	return e.updatedOn
}

func (e *CacheEntry) ResponseModel() interface{} {
	return e.responseModel
}

func (e *CacheEntry) ResponseBodyPlain() []byte {
	return e.responseBodyPlain
}

func (e *CacheEntry) ResponseBodyBrotli() []byte {
	return e.responseBodyBrotli
}

func (e *CacheEntry) ResponseBodyGzip() []byte {
	return e.responseBodyGzip
}

func (e *CacheEntry) ResponseBodyDeflate() []byte {
	return e.responseBodyDeflate
}

func (e *CacheEntry) ResponseHeader() http.Header {
	return e.responseHeader
}

func (e *CacheEntry) ResponseStatusCode() int {
	return e.responseStatusCode
}

func (e *CacheEntry) String() string {
	return fmt.Sprintf("{ updated:%s statusCode:%d model:%t bodyPlain:%t bodyBr:%t bodyGzip:%t bodyDeflate:%t }",
		e.updatedOn.Format(time.RFC3339),
		e.responseStatusCode,
		e.responseModel != nil,
		e.responseBodyPlain != nil,
		e.responseBodyBrotli != nil,
		e.responseBodyGzip != nil,
		e.responseBodyDeflate != nil,
	)
}

func hash(body []byte, params string) string {
	uniqueString := fmt.Sprintf("%s-%s", body, params)
	s := md5.Sum([]byte(uniqueString))
	return string(s[:])
}
