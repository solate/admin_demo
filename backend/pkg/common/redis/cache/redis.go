package cache

import (
	"fmt"
	"sync"

	"github.com/zeromicro/go-zero/core/stores/redis"
)

var (
	instance *redis.Redis
	once     sync.Once
)

// RedisConfig Redis配置
type RedisConfig struct {
	Host string
	Port int
	Type string
	Pass string
}

// NewRedis 获取Redis实例（单例模式）
func NewRedis(c RedisConfig) *redis.Redis {
	once.Do(func() {
		instance = redis.New(fmt.Sprintf("%s:%d", c.Host, c.Port), func(r *redis.Redis) {
			r.Type = c.Type
			r.Pass = c.Pass
		})
	})
	return instance
}
