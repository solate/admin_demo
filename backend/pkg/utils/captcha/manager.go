package captcha

import (
	"time"

	"github.com/mojocn/base64Captcha"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type CaptchaManager struct {
	store *RedisStore
}

type RedisStore struct {
	redisClient *redis.Redis
	expiration  time.Duration
	keyPrefix   string
}

// NewCaptchaManager 创建一个新的验证码管理器
func NewCaptchaManager(redisClient *redis.Redis) *CaptchaManager {
	store := &RedisStore{
		redisClient: redisClient,
		expiration:  time.Minute * 5, // 验证码5分钟过期
		keyPrefix:   "captcha:",
	}

	return &CaptchaManager{
		store: store,
	}
}

// Generate 生成图形验证码
func (m *CaptchaManager) Generate() (id, b64s, answer string, err error) {
	driver := base64Captcha.NewDriverDigit(80, 240, 4, 0.7, 80)
	c := base64Captcha.NewCaptcha(driver, m.store)
	return c.Generate()
}

// Verify 验证验证码
func (m *CaptchaManager) Verify(id, answer string) bool {
	return m.store.Verify(id, answer, true)
}

// Set 实现base64Captcha.Store接口
func (s *RedisStore) Set(id string, value string) error {
	return s.redisClient.Setex(s.keyPrefix+id, value, int(s.expiration.Seconds()))
}

// Get 实现base64Captcha.Store接口
func (s *RedisStore) Get(id string, clear bool) string {
	val, err := s.redisClient.Get(s.keyPrefix + id)
	if err != nil {
		return ""
	}
	if clear {
		_, err = s.redisClient.Del(s.keyPrefix + id)
		if err != nil {
			return ""
		}
	}
	return val
}

// Verify 实现base64Captcha.Store接口
func (s *RedisStore) Verify(id, answer string, clear bool) bool {
	val := s.Get(id, clear)
	return val == answer
}
