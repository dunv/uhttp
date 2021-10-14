package uhttp

import (
	"fmt"
	"net/http"
)

var cacheSizeHandler = func(u *UHTTP) Handler {
	return NewHandler(
		WithGet(func(r *http.Request, ret *int) interface{} {
			res := map[string]map[string]uint64{}
			totalSize := uint64(0)
			totalEntries := uint64(0)
			u.cacheLock.RLock()
			for pattern, cache := range u.cache {
				size := cache.Size()
				totalSize += size
				entries := uint64(cache.Count())
				totalEntries += entries
				res[pattern] = map[string]uint64{
					"sizeInBytes": size,
					"entries":     entries,
				}
			}
			u.cacheLock.RUnlock()
			res["total"] = map[string]uint64{
				"sizeInBytes":  totalSize,
				"totalEntries": totalEntries,
			}
			return res
		}),
	)
}

var cacheDebugHandler = func(u *UHTTP) Handler {
	return NewHandler(
		WithGet(func(r *http.Request, ret *int) interface{} {
			res := map[string]map[string]string{}
			u.cacheLock.RLock()
			for pattern, cache := range u.cache {
				res[pattern] = map[string]string{}
				keys := cache.Keys()
				for _, key := range keys {
					if data, ok := cache.Get(key); ok {
						res[pattern][fmt.Sprintf("%x", key)] = data.String()
					}
				}
			}
			u.cacheLock.RUnlock()
			return res
		}),
	)
}

var cacheClearHandler = func(u *UHTTP) Handler {
	return NewHandler(
		WithPost(func(r *http.Request, ret *int) interface{} {
			deletedEntries := 0
			u.cacheLock.RLock()
			for _, c := range u.cache {
				keys := c.Keys()
				for _, key := range keys {
					c.Delete(key)
					deletedEntries++
				}
			}
			u.cacheLock.RUnlock()
			return map[string]int{
				"deletedEntries": deletedEntries,
			}
		}),
	)
}

var specificCacheClearHandler = func(u *UHTTP, c cache) Handler {
	return NewHandler(
		WithPost(func(r *http.Request, ret *int) interface{} {
			deletedEntries := 0
			keys := c.Keys()
			for _, key := range keys {
				c.Delete(key)
				deletedEntries++
			}
			return map[string]int{
				"deletedEntries": deletedEntries,
			}
		}),
	)
}
