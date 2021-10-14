package cache

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"sync"
	"time"
	"unsafe"
)

func NewCache(maxAge time.Duration, persistEncodings bool) *Cache {
	return &Cache{
		mu:               &sync.RWMutex{},
		maxAge:           maxAge,
		persistEncodings: persistEncodings,
		data:             map[string]CacheEntry{},
	}
}

type Cache struct {
	mu               *sync.RWMutex
	maxAge           time.Duration
	persistEncodings bool
	data             map[string]CacheEntry
}

func (c Cache) MaxAge() time.Duration {
	return c.maxAge
}

func (c Cache) Set(
	requestBody []byte,
	requestParams string,
	requestHeader http.Header,
	responseModel interface{},
	responseBody []byte,
	responseHeader http.Header,
	responseStatusCode int,
) {
	dataCopy := make([]byte, len(responseBody))
	copy(dataCopy, responseBody)
	key := hash(requestBody, requestParams)

	c.mu.Lock()
	defer c.mu.Unlock()
	e := CacheEntry{
		updatedOn:          time.Now(),
		responseModel:      responseModel,
		responseHeader:     responseHeader,
		responseStatusCode: responseStatusCode,
	}
	if c.persistEncodings {
		encoding := responseHeader.Get("Content-Encoding")
		switch encoding {
		case "br":
			e.responseBodyBrotli = responseBody
		case "gzip":
			e.responseBodyGzip = responseBody
		case "deflate":
			e.responseBodyDeflate = responseBody
		case "":
			e.responseBody = responseBody
		}
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
		total += uint64(len(entry.responseBody))
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
	responseBody        []byte
	responseBodyGzip    []byte
	responseBodyBrotli  []byte
	responseBodyDeflate []byte
	responseHeader      http.Header
	responseStatusCode  int
}

func (e *CacheEntry) Clone() CacheEntry {
	var responseBodyCopy []byte
	if e.responseBody != nil {
		responseBodyCopy = make([]byte, len(e.responseBody))
		copy(responseBodyCopy, e.responseBody)
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
		responseBody:        responseBodyCopy,
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

func (e *CacheEntry) ResponseBody() []byte {
	return e.responseBody
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
	return fmt.Sprintf("{ updated:%s statusCode:%d model:%t body:%t bodyBr:%t bodyGzip:%t bodyDeflate:%t }",
		e.updatedOn.Format(time.RFC3339),
		e.responseStatusCode,
		e.responseModel != nil,
		e.responseBody != nil,
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
