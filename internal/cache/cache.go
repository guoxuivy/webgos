package cache

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"
	"sync"
	"time"

	"webgos/common/json"

	gocache "github.com/patrickmn/go-cache"
)

const (
	DefaultExpiration time.Duration = gocache.DefaultExpiration //使用创建时的默认过期时间
	NoExpiration      time.Duration = gocache.NoExpiration      //永不过期，直到显式 Delete
)

var (
	defaultCache *Cache
	once         sync.Once
)

// 使用go-cache注意集群部署时缓存不一致问题
// 解决方法：使用分布式缓存，如Redis
type Cache struct {
	cache       *gocache.Cache
	mu          sync.RWMutex
	prefixIndex map[string]map[string]struct{}
}

// ICache 是缓存存储抽象接口。
type ICache interface {
	Get(key string) (any, bool)
	Set(key string, value any, duration time.Duration)
	Delete(key string)
	DeleteByPrefix(prefix string)
	Flush()
}

// 保留扩展其它类型缓存
func GetCache() ICache {
	once.Do(func() {
		defaultCache = &Cache{
			cache:       gocache.New(5*time.Minute, 10*time.Minute),
			prefixIndex: make(map[string]map[string]struct{}),
		}
		defaultCache.cache.OnEvicted(func(k string, v any) {
			defaultCache.removeFromIndex(k)
		})
	})
	return defaultCache
}

// 获取key的前缀，假设key格式为 "park:page@a3eracfdsfaer23423" 返回 "park:page"
func (c *Cache) keyPrefix(key string) string {
	if i := strings.IndexByte(key, '@'); i > 0 {
		return key[:i]
	}
	return ""
}
func (c *Cache) addToIndex(key string) {
	prefix := c.keyPrefix(key)
	if prefix == "" {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.prefixIndex[prefix]; !exists {
		c.prefixIndex[prefix] = make(map[string]struct{})
	}
	c.prefixIndex[prefix][key] = struct{}{}
}

func (c *Cache) removeFromIndex(key string) {
	prefix := c.keyPrefix(key)
	if prefix == "" {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if indexMap, exists := c.prefixIndex[prefix]; exists {
		delete(indexMap, key)
		if len(indexMap) == 0 {
			delete(c.prefixIndex, prefix)
		}
	}
}

// Cache 的基础方法，不包含业务类型封装。
func (c *Cache) Get(key string) (any, bool) {
	return c.cache.Get(key)
}

func (c *Cache) Set(key string, value any, duration time.Duration) {
	c.cache.Set(key, value, duration)
	c.addToIndex(key)
}

func (c *Cache) Delete(key string) {
	c.cache.Delete(key)
	c.removeFromIndex(key)
}

// 高性能的前缀删除实现
func (c *Cache) DeleteByPrefix(prefix string) {
	c.mu.RLock()
	keysMap, ok := c.prefixIndex[prefix]
	if ok {
		keys := make([]string, 0, len(keysMap))
		for k := range keysMap {
			keys = append(keys, k)
		}
		c.mu.RUnlock()

		for _, k := range keys {
			c.cache.Delete(k)
		}

		c.mu.Lock()
		delete(c.prefixIndex, prefix)
		c.mu.Unlock()
		return
	}
	c.mu.RUnlock()

	// 索引没有命中时全表扫描
	items := c.cache.Items()
	for k := range items {
		if strings.HasPrefix(k, prefix) {
			c.cache.Delete(k)
		}
	}
}

// Flush 清空所有缓存
func (c *Cache) Flush() {
	c.cache.Flush()
}

func GenerateKey(prefix string, query any) string {
	queryBytes, _ := json.Marshal(query)
	hash := md5.Sum(queryBytes)
	return fmt.Sprintf("%s@%s", prefix, hex.EncodeToString(hash[:]))
}
