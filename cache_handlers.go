package uhttp

import (
	"fmt"
	"net/http"
)

var cacheDetailsHandler = func(u *UHTTP, middlewares ...Middleware) Handler {
	return NewHandler(
		WithMiddlewares(middlewares...),
		WithOptionalGet(R{
			"parsable": BOOL,
		}),
		WithGet(func(r *http.Request, ret *int) interface{} {
			parsable := GetAsBool("parsable", r)
			res := make(map[string]map[string]interface{})
			u.cacheLock.RLock()
			defer u.cacheLock.RUnlock()
			for pattern, c := range u.cache {
				res[pattern] = make(map[string]interface{})
				keys := c.Keys()
				for _, key := range keys {
					if data, ok := c.GetByKey(key); ok {
						if parsable != nil && *parsable {
							stats, err := data.Stats(c)
							if err != nil {
								return err
							}
							res[pattern][fmt.Sprintf("%x", key)] = stats

						} else {
							res[pattern][fmt.Sprintf("%x", key)] = data.String()
						}
					}
				}
			}
			return res
		}),
	)
}

var cacheClearHandler = func(u *UHTTP, middlewares ...Middleware) Handler {
	return NewHandler(
		WithMiddlewares(middlewares...),
		WithOptionalGet(R{
			"path": STRING,
			"hash": STRING,
		}),
		WithPost(func(r *http.Request, ret *int) interface{} {
			deletedEntries := 0
			path := GetAsString("path", r)
			hash := GetAsString("hash", r)

			u.cacheLock.RLock()
			defer u.cacheLock.RUnlock()

			for _, c := range u.cache {
				if path != nil && *path != c.HandlerPattern() {
					continue
				}

				keys := c.Keys()
				for _, key := range keys {
					if hash != nil && *hash != key {
						continue
					}

					c.Delete(key)
					deletedEntries++
				}
			}
			return map[string]int{
				"deletedEntries": deletedEntries,
			}
		}),
	)
}

var cacheSizeHandler = func(u *UHTTP, middlewares ...Middleware) Handler {
	return NewHandler(
		WithMiddlewares(middlewares...),
		WithGet(func(r *http.Request, ret *int) interface{} {
			res := make(map[string]map[string]interface{})
			totalSize := uint64(0)
			totalEntries := uint64(0)
			u.cacheLock.RLock()
			for pattern, cache := range u.cache {
				size := cache.Size()
				totalSize += size
				entries := uint64(cache.Count())
				totalEntries += entries
				res[pattern] = map[string]interface{}{
					"sizeInBytes": size,
					"entries":     entries,
				}
			}
			u.cacheLock.RUnlock()
			res["total"] = map[string]interface{}{
				"sizeInBytes": totalSize,
				"entries":     totalEntries,
			}
			return res
		}),
	)
}
