package cache

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"reflect"
	"sync"
	"time"

	"webgos/common/json"

	gocache "github.com/patrickmn/go-cache"
)

const (
	DefaultExpiration time.Duration = gocache.DefaultExpiration
	NoExpiration      time.Duration = gocache.NoExpiration //永不过期 ，直到显式 Delete
)

const (
	ParkPagePrefix = "park:page"
	UserMenuPrefix = "user:menusByUserID"
)

var (
	defaultCache *Cache
	once         sync.Once
)

// 使用go-cache注意集群部署时缓存不一致问题
// 解决方法：使用分布式缓存，如Redis
type Cache struct {
	cache *gocache.Cache
}

type PageCache struct {
	Items any `json:"items"`
	Total int `json:"total"`
}

func GetCache() *Cache {
	once.Do(func() {
		defaultCache = &Cache{
			cache: gocache.New(5*time.Minute, 10*time.Minute),
		}
	})
	return defaultCache
}

func (c *Cache) GetPage(key string, list any, total *int) bool {
	data, found := c.cache.Get(key)
	if !found {
		return false
	}

	var cacheData PageCache
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

func (c *Cache) SetPage(key string, items any, total int) {
	c.cache.Set(key, PageCache{Items: items, Total: total}, 5*time.Minute)
}

func (c *Cache) SetPageWithExpiration(key string, items any, total int, expiration time.Duration) {
	c.cache.Set(key, PageCache{Items: items, Total: total}, expiration)
}

func (c *Cache) Get(key string) (any, bool) {
	return c.cache.Get(key)
}

func (c *Cache) Set(key string, value any, duration time.Duration) {
	c.cache.Set(key, value, duration)
}

func (c *Cache) Delete(key string) {
	c.cache.Delete(key)
}

func (c *Cache) DeleteByPrefix(prefix string) {
	items := c.cache.Items()
	for k := range items {
		if len(k) >= len(prefix) && k[:len(prefix)] == prefix {
			c.cache.Delete(k)
		}
	}
}

func (c *Cache) Flush() {
	c.cache.Flush()
}

func GenerateKey(prefix string, query any) string {
	queryBytes, _ := json.Marshal(query)
	hash := md5.Sum(queryBytes)
	return fmt.Sprintf("%s:%s", prefix, hex.EncodeToString(hash[:]))
}
