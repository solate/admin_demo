package jwt

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

func TestParseToken(t *testing.T) {
	config := JWTConfig{
		AccessExpire:  3600,
		AccessSecret:  []byte("test"),
		RefreshExpire: 86400,
		RefreshSecret: []byte("refresh_test"),
	}

	// 生成token对
	tokenPair, err := GenerateTokenPair("1", "test", "", 1, config)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("access_token:", tokenPair.AccessToken)
	fmt.Println("refresh_token:", tokenPair.RefreshToken)

	// 测试解析带Bearer前缀的access token
	claims, err := ParseAccessToken("Bearer "+tokenPair.AccessToken, []byte("test"))
	if err != nil {
		t.Fatal(err)
	}
	t.Log("解析带Bearer前缀的access token成功:", claims.UserID, claims.TenantCode, claims.TokenType)

	// 测试解析不带Bearer前缀的access token
	claims, err = ParseAccessToken(tokenPair.AccessToken, []byte("test"))
	if err != nil {
		t.Fatal(err)
	}
	t.Log("解析不带Bearer前缀的access token成功:", claims.UserID, claims.TenantCode, claims.TokenType)

	// 测试解析refresh token
	refreshClaims, err := ParseRefreshToken(tokenPair.RefreshToken, []byte("refresh_test"))
	if err != nil {
		t.Fatal(err)
	}
	t.Log("解析refresh token成功:", refreshClaims.UserID, refreshClaims.TenantCode, refreshClaims.TokenType, refreshClaims.TokenID)

	// 测试无效的token
	_, err = ParseAccessToken("invalid_token", []byte("test"))
	if err == nil {
		t.Fatal("期望解析无效token时返回错误，但没有")
	}
	t.Log("无效access token测试通过")

	// 测试token类型验证
	_, err = ParseRefreshToken(tokenPair.AccessToken, []byte("refresh_test"))
	if err == nil {
		t.Fatal("期望使用access token解析refresh token时返回错误，但没有")
	}
	t.Log("token类型验证测试通过")
}

func TestRemoveBearerPrefix(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "标准Bearer前缀",
			input:    "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
			expected: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
		},
		{
			name:     "多个空格Bearer前缀",
			input:    "Bearer  eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
			expected: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
		},
		{
			name:     "多个空格Bearer前缀2",
			input:    "Bearer   eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
			expected: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
		},
		{
			name:     "首尾空格",
			input:    "  Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...  ",
			expected: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
		},
		{
			name:     "无Bearer前缀",
			input:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
			expected: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
		},
		{
			name:     "空字符串",
			input:    "",
			expected: "",
		},
		{
			name:     "只有空格",
			input:    "   ",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RemoveBearerPrefix(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestJWTManager_Integration(t *testing.T) {
	// 初始化测试环境
	rdb := redis.New("localhost:6379")
	config := JWTConfig{
		AccessSecret:  []byte("test_access_secret"),
		AccessExpire:  3600, // 1小时
		RefreshSecret: []byte("test_refresh_secret"),
		RefreshExpire: 604800, // 7天
	}

	jwtManager := NewJWTManager(config, rdb)
	ctx := context.Background()

	// 测试数据
	userID := "test_user_integration"
	tenantCode := "test_tenant"
	roleCode := "admin"
	source := 1

	t.Run("TestCompleteFlow", func(t *testing.T) {
		// 1. 测试登录
		loginResult, err := jwtManager.Login(ctx, userID, tenantCode, roleCode, source)
		assert.NoError(t, err)
		assert.NotEmpty(t, loginResult.AccessToken)
		assert.NotEmpty(t, loginResult.RefreshToken)

		// 2. 测试验证access token
		claims, err := ParseAccessToken(loginResult.AccessToken, []byte(config.AccessSecret))
		assert.NoError(t, err)
		assert.Equal(t, userID, claims.UserID)

		// 3. 测试刷新token
		time.Sleep(1 * time.Second) // 确保时间戳不同
		refreshResult, err := jwtManager.RefreshToken(ctx, loginResult.RefreshToken)
		assert.NoError(t, err)
		assert.NotEmpty(t, refreshResult.AccessToken)
		assert.NotEmpty(t, refreshResult.RefreshToken)
		assert.NotEqual(t, loginResult.AccessToken, refreshResult.AccessToken)

		// 4. 测试获取活跃token数量
		count, err := jwtManager.GetUserActiveTokens(ctx, userID)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, count, 1)

		// 5. 测试撤销token
		// 解析 refresh token 获取 tokenID
		refreshClaims, err := ParseRefreshToken(refreshResult.RefreshToken, config.RefreshSecret)
		assert.NoError(t, err)

		err = jwtManager.RevokeToken(ctx, userID, refreshClaims.TokenID)
		assert.NoError(t, err)

		// 6. 清理
		err = jwtManager.RevokeAllUserTokens(ctx, userID)
		assert.NoError(t, err)
	})

	t.Run("TestBlacklistFlow", func(t *testing.T) {
		testUserID := userID + "_blacklist"

		// 生成token
		loginResult, err := jwtManager.Login(ctx, testUserID, tenantCode, roleCode, source)
		assert.NoError(t, err)

		// 验证token有效
		_, err = jwtManager.ValidateAccessTokenWithBlacklist(ctx, loginResult.AccessToken)
		assert.NoError(t, err)

		// 执行单设备登出
		err = jwtManager.LogoutSingleDevice(ctx, loginResult.AccessToken)
		assert.NoError(t, err)

		// 验证token已失效
		_, err = jwtManager.ValidateAccessTokenWithBlacklist(ctx, loginResult.AccessToken)
		assert.Error(t, err)

		// 清理
		err = jwtManager.RevokeAllUserTokens(ctx, testUserID)
		assert.NoError(t, err)
	})

	t.Run("TestMultiDeviceFlow", func(t *testing.T) {
		testUserID := userID + "_multi_device"

		// 生成多个token（模拟多设备）
		device1, err := jwtManager.Login(ctx, testUserID, tenantCode, roleCode, 1)
		assert.NoError(t, err)

		time.Sleep(1 * time.Second)
		device2, err := jwtManager.Login(ctx, testUserID, tenantCode, roleCode, 2)
		assert.NoError(t, err)

		// 验证两个设备的token都有效
		_, err = jwtManager.ValidateAccessTokenWithBlacklist(ctx, device1.AccessToken)
		assert.NoError(t, err)
		_, err = jwtManager.ValidateAccessTokenWithBlacklist(ctx, device2.AccessToken)
		assert.NoError(t, err)

		// 执行全设备登出
		err = jwtManager.LogoutAllDevices(ctx, testUserID)
		assert.NoError(t, err)

		// 验证所有设备的token都失效
		_, err = jwtManager.ValidateAccessTokenWithBlacklist(ctx, device1.AccessToken)
		assert.Error(t, err)
		_, err = jwtManager.ValidateAccessTokenWithBlacklist(ctx, device2.AccessToken)
		assert.Error(t, err)

		// 清理
		err = jwtManager.RevokeAllUserTokens(ctx, testUserID)
		assert.NoError(t, err)
	})

	// 清理测试数据
	t.Cleanup(func() {
		// 清理所有测试相关的Redis数据
		patterns := []string{
			"refresh_token:*",
			"user_refresh_tokens:*",
			"blacklist_access_token:*",
			"blacklist_user:*",
		}

		for _, pattern := range patterns {
			if keys, err := rdb.Keys(pattern); err == nil {
				for _, key := range keys {
					rdb.Del(key)
				}
			}
		}
	})
}

// 解析token
func TestParseToken2(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMTM4MTU2OTc2NjI1NTQxOTc0IiwidGVuYW50X2NvZGUiOiJkZWZhdWx0Iiwicm9sZV9jb2RlIjoiIiwic291cmNlIjoxLCJ0b2tlbl90eXBlIjoicmVmcmVzaCIsInRva2VuX2lkIjoiMTM4NzM0NjAxNTkyOTI4OTc5IiwiZXhwIjoxNzU3ODE1NDY2LCJuYmYiOjE3NTUyMjM0NjYsImlhdCI6MTc1NTIyMzQ2Nn0.9G_xwoy0-9pWFhARxi_dKwywwbydNuciSHyujwfFQ4E"
	claims, err := ParseRefreshToken(token, []byte("B0OXteNV0+ZPynWKvOR75jSSoYnd/p9LvZVWqGmeh9Q="))
	fmt.Println(claims, err)
	fmt.Println("UserID:", claims.UserID)
	fmt.Println("TenantCode:", claims.TenantCode)
	fmt.Println("RoleCode:", claims.RoleCode)
	fmt.Println("Source:", claims.Source)
	fmt.Println("TokenType:", claims.TokenType)
	fmt.Println("TokenID:", claims.TokenID) // 这是Redis中的key
	fmt.Println("ExpiresAt:", claims.RegisteredClaims.ExpiresAt)
	fmt.Println("IssuedAt:", claims.RegisteredClaims.IssuedAt)
	fmt.Println("NotBefore:", claims.RegisteredClaims.NotBefore)

	// 模拟检查Redis中是否存在这个token
	fmt.Printf("\n检查Redis Key应该是: refresh_token:%s\n", claims.TokenID)
	fmt.Printf("用户token列表Key应该是: user_refresh_tokens:%s\n", claims.UserID)
}
