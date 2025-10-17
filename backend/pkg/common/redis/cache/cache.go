package cache

import (
	"encoding/json"
	"time"

	"github.com/zeromicro/go-zero/core/stores/redis"
)

// Cache 定义缓存接口
type Cache interface {
	// Set 设置缓存
	Set(key string, value any, expiration time.Duration) error
	// Get 获取缓存
	Get(key string, value any) error
	// Delete 删除缓存
	Delete(key string) error
	// Exists 检查key是否存在
	Exists(key string) (bool, error)
	// TTL 获取key的过期时间
	TTL(key string) (time.Duration, error)
}

// CacheManager 缓存管理器
type CacheManager struct {
	redis *redis.Redis
}

// NewCacheManager 创建缓存管理器实例
func NewCacheManager(redis *redis.Redis) *CacheManager {
	return &CacheManager{redis: redis}
}

// Set 实现Cache接口的Set方法
func (c *CacheManager) Set(key string, value any, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	if expiration > 0 {
		return c.redis.Setex(key, string(data), int(expiration.Seconds()))
	}
	return c.redis.Set(key, string(data))
}

// Get 实现Cache接口的Get方法
func (c *CacheManager) Get(key string, value any) error {
	data, err := c.redis.Get(key)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(data), value)
}

// Delete 实现Cache接口的Delete方法
func (c *CacheManager) Delete(key string) error {
	_, err := c.redis.Del(key)
	return err
}

// Exists 实现Cache接口的Exists方法
func (c *CacheManager) Exists(key string) (bool, error) {
	return c.redis.Exists(key)
}

// TTL 实现Cache接口的TTL方法
func (c *CacheManager) TTL(key string) (time.Duration, error) {
	ttl, err := c.redis.Ttl(key)
	if err != nil {
		return 0, err
	}
	return time.Duration(ttl) * time.Second, nil
}
