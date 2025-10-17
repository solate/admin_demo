package token

import (
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/redis"
)

// TokenConfig 令牌配置
type TokenConfig struct {
	Prefix     string
	Expiration time.Duration
}

// DefaultConfig 默认配置
var DefaultConfig = TokenConfig{
	Prefix:     "token:",
	Expiration: 5 * time.Minute,
}

// TokenManager 防重令牌管理器
type TokenManager struct {
	rdb    *redis.Redis
	config TokenConfig
}

// luaScripts 统一管理Lua脚本
var luaScripts = struct {
	// 获取令牌并设置过期时间
	acquire string
	// 验证并删除令牌
	validateAndDelete string
}{
	acquire: `
		local key = KEYS[1]
		local value = ARGV[1]
		local expiration = tonumber(ARGV[2])
		if redis.call('setnx', key, value) == 1 then
			redis.call('pexpire', key, expiration)
			return 1
		end
		return 0
	`,
	validateAndDelete: `
		local key = KEYS[1]
		local value = ARGV[1]
		local current = redis.call('get', key)
		if current == value then
			redis.call('del', key)
			return 1
		end
		return 0
	`,
}

// NewTokenManager 创建防重令牌管理器
func NewTokenManager(rdb *redis.Redis, config ...TokenConfig) *TokenManager {
	cfg := DefaultConfig
	if len(config) > 0 {
		cfg = config[0]
	}
	return &TokenManager{
		rdb:    rdb,
		config: cfg,
	}
}

// AcquireToken 获取防重令牌
// key 为业务唯一标识，如：用户ID+操作类型
// expiration 为可选的令牌过期时间，如果为0则使用默认过期时间
func (tm *TokenManager) AcquireToken(key string, expiration ...time.Duration) (string, error) {
	exp := tm.config.Expiration
	if len(expiration) > 0 && expiration[0] > 0 {
		exp = expiration[0]
	}

	// 生成唯一的token值
	tokenValue := fmt.Sprintf("%d", time.Now().UnixNano())
	redisKey := tm.getRedisKey(key)

	// 使用Lua脚本实现原子操作
	result, err := tm.rdb.Eval(
		luaScripts.acquire,
		[]string{redisKey},
		[]string{tokenValue, fmt.Sprintf("%d", exp.Milliseconds())},
	)
	if err != nil {
		return "", fmt.Errorf("acquire token failed: %v", err)
	}

	if result.(int64) != 1 {
		return "", fmt.Errorf("token already exists")
	}

	return tokenValue, nil
}

// ValidateAndDeleteToken 验证并删除令牌（原子操作）
func (tm *TokenManager) ValidateAndDeleteToken(key, value string) (bool, error) {
	redisKey := tm.getRedisKey(key)
	result, err := tm.rdb.Eval(
		luaScripts.validateAndDelete,
		[]string{redisKey},
		[]string{value},
	)
	if err != nil {
		return false, fmt.Errorf("validate and delete token failed: %v", err)
	}

	return result.(int64) == 1, nil
}

// ClearToken 清除指定key的令牌
func (tm *TokenManager) ClearToken(key string) error {
	redisKey := tm.getRedisKey(key)
	_, err := tm.rdb.Del(redisKey)
	if err != nil {
		return fmt.Errorf("clear token failed: %v", err)
	}
	return nil
}

// getRedisKey 生成Redis key
func (tm *TokenManager) getRedisKey(key string) string {
	return tm.config.Prefix + key
}
