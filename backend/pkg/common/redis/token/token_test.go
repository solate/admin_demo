package token

import (
	"testing"
	"time"

	"github.com/zeromicro/go-zero/core/stores/redis"
)

func TestTokenManager_ClearToken(t *testing.T) {
	// 创建Redis客户端（这里使用mock或者真实的Redis）
	// 注意：在实际测试中，你可能需要使用testcontainers或者mock Redis
	rdb := &redis.Redis{} // 这里需要根据实际情况初始化

	tm := NewTokenManager(rdb, TokenConfig{
		Prefix:     "test:token:",
		Expiration: 1 * time.Minute,
	})

	// 测试清除不存在的token
	err := tm.ClearToken("test_user")
	if err != nil {
		t.Errorf("ClearToken failed: %v", err)
	}

	// 测试清除存在的token
	// 1. 先获取一个token
	token, err := tm.AcquireToken("test_user")
	if err != nil {
		t.Errorf("AcquireToken failed: %v", err)
	}

	// 2. 验证token存在
	valid, err := tm.ValidateAndDeleteToken("test_user", token)
	if err != nil {
		t.Errorf("ValidateAndDeleteToken failed: %v", err)
	}
	if !valid {
		t.Error("Token should be valid")
	}

	// 3. 再次获取token
	token2, err := tm.AcquireToken("test_user")
	if err != nil {
		t.Errorf("AcquireToken after clear failed: %v", err)
	}

	// 4. 清除token
	err = tm.ClearToken("test_user")
	if err != nil {
		t.Errorf("ClearToken failed: %v", err)
	}

	// 5. 验证token已被清除
	valid, err = tm.ValidateAndDeleteToken("test_user", token2)
	if err != nil {
		t.Errorf("ValidateAndDeleteToken after clear failed: %v", err)
	}
	if valid {
		t.Error("Token should be invalid after clear")
	}
}
