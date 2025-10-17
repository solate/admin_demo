package svc

import (
	"context"

	"admin_backend/app/admin/internal/config"
	"admin_backend/app/admin/internal/middleware"
	"admin_backend/pkg/common/casbin"

	"admin_backend/pkg/ent"
	"admin_backend/pkg/ent/generated"

	"admin_backend/pkg/utils/jwt"

	"admin_backend/pkg/common/redis/cache"
	"admin_backend/pkg/utils/captcha"
	"admin_backend/pkg/utils/cron"

	_ "github.com/lib/pq"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
)

type ServiceContext struct {
	Config           config.Config
	DB               *ent.Client
	Redis            *redis.Redis
	AuthMiddleware   rest.Middleware
	ClientMiddleware rest.Middleware         // 客户端认证中间件
	CaptchaManager   *captcha.CaptchaManager // 验证码管理
	CasbinManager    *casbin.CasbinManager
	JWTManager       *jwt.JWTManager // JWT统一管理器（合并了原来的TokenManager和RefreshTokenManager）
	CronManager      *cron.CronManager
}

func NewServiceContext(c config.Config) *ServiceContext {

	client := initOrm(c)
	rdb := initRedis(c)

	// 先初始化JWT管理器
	jwtManager := initJWTManager(c, rdb)

	return &ServiceContext{
		Config:         c,
		DB:             client,
		Redis:          initRedis(c),
		AuthMiddleware: middleware.NewAuthMiddleware(c, jwtManager).Handle,
		CaptchaManager: captcha.NewCaptchaManager(rdb),
		CasbinManager:  casbin.NewCasbinManager(client),
		JWTManager:     jwtManager, // 使用已初始化的JWT管理器
		CronManager:    cron.NewCronManager(),
	}
}

// initOrm
func initOrm(c config.Config) *ent.Client {
	ops := make([]generated.Option, 0)
	if c.ShowSQL {
		ops = append(ops, generated.Debug())
	}
	client, err := ent.NewClient(context.Background(), c.DataSource, ops...)
	if err != nil {
		logx.Errorf("ent.Open error: %v", err)
		panic(err)
	}
	return client
}

func initRedis(c config.Config) *redis.Redis {
	return cache.NewRedis(cache.RedisConfig{
		Host: c.Redis.Host,
		Port: c.Redis.Port,
		Type: c.Redis.Type,
		Pass: c.Redis.Pass,
	})
}

// 初始化JWT管理器（统一管理器）
func initJWTManager(c config.Config, rdb *redis.Redis) *jwt.JWTManager {
	jwtConfig := jwt.JWTConfig{
		AccessSecret:  []byte(c.JwtAuth.AccessSecret),
		AccessExpire:  c.JwtAuth.AccessExpire,
		RefreshSecret: []byte(c.JwtAuth.RefreshSecret),
		RefreshExpire: c.JwtAuth.RefreshExpire,
	}
	return jwt.NewJWTManager(jwtConfig, rdb)
}
