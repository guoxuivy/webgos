package cache

import (
	"reflect"
	"time"

	"webgos/common/json"
)

type pageCache struct {
	Items any `json:"items"`
	Total int `json:"total"`
}

func GetPage(key string, list any, total *int) bool {
	data, found := GetCache().Get(key)
	if !found {
		return false
	}

	var cacheData pageCache
	cacheDataBytes, err := json.Marshal(data)
	if err != nil {
		return false
	}
	if err = json.Unmarshal(cacheDataBytes, &cacheData); err != nil {
		return false
	}

	if total != nil {
		*total = cacheData.Total
	}

	itemsBytes, err := json.Marshal(cacheData.Items)
	if err != nil {
		return false
	}

	listValue := reflect.ValueOf(list)
	if listValue.Kind() != reflect.Ptr || listValue.Elem().Kind() != reflect.Slice {
		return false
	}

	if err := json.Unmarshal(itemsBytes, list); err != nil {
		return false
	}

	return true
}

func SetPage(key string, items any, total int) {
	SetPageWithExpiration(key, items, total, DefaultExpiration)
}

func SetPageWithExpiration(key string, items any, total int, expiration time.Duration) {
	GetCache().Set(key, pageCache{Items: items, Total: total}, expiration)
}
