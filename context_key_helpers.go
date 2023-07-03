package uhttp

import "reflect"

func ContextKeysFromMap(m interface{}) []ContextKey {
	s := reflect.ValueOf(m)
	if s.Kind() != reflect.Map {
		panic("StringKeysFromMap() given a non-map type")
	}

	keys := make([]ContextKey, len(s.MapKeys()))
	i := 0
	for _, mapKey := range s.MapKeys() {
		key, ok := mapKey.Interface().(ContextKey)
		if !ok {
			panic("StringKeysFromMap() given a non-string key")
		}
		keys[i] = key
		i++
	}
	return keys
}
