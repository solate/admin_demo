package jwt

import (
	"admin_backend/pkg/utils/idgen"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID     string `json:"user_id"`
	TenantCode string `json:"tenant_code"`
	RoleCode   string `json:"role_code"`
	Source     int    `json:"source"`             // 来源  1-电脑端pc 2-手机端h5 3-小程序miniapp
	TokenType  string `json:"token_type"`         // token类型: access 或 refresh
	TokenID    string `json:"token_id,omitempty"` // refresh token的唯一标识
	jwt.RegisteredClaims
}

type JWTConfig struct {
	AccessSecret  []byte // access token密钥
	AccessExpire  int64  // access token过期时间
	RefreshSecret []byte // refresh token密钥
	RefreshExpire int64  // refresh token过期时间
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenID      string `json:"token_id,omitempty"` // refresh token的唯一标识
}

// 生成token对 (access + refresh)
func GenerateTokenPair(userID, tenantCode, roleCode string, source int, config JWTConfig) (*TokenPair, error) {
	// 生成refresh token的唯一ID
	tokenID, err := idgen.GenerateUUID()
	if err != nil {
		return nil, fmt.Errorf("生成token ID失败: %w", err)
	}

	// 生成access token
	accessToken, err := generateTokenWithID(userID, tenantCode, roleCode, source, "access", "", config.AccessSecret, config.AccessExpire)
	if err != nil {
		return nil, fmt.Errorf("生成access token失败: %w", err)
	}

	// 生成refresh token (带tokenID)
	refreshToken, err := generateTokenWithID(userID, tenantCode, roleCode, source, "refresh", tokenID, config.RefreshSecret, config.RefreshExpire)
	if err != nil {
		return nil, fmt.Errorf("生成refresh token失败: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenID:      tokenID,
	}, nil
}

// 内部生成token方法 (支持tokenID)
func generateTokenWithID(userID, tenantCode, roleCode string, source int, tokenType, tokenID string, secret []byte, expire int64) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID:     userID,
		TenantCode: tenantCode,
		RoleCode:   roleCode,
		Source:     source,
		TokenType:  tokenType,
		TokenID:    tokenID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(expire) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// 解析access token
func ParseAccessToken(tokenString string, accessSecret []byte) (*Claims, error) {
	return parseToken(tokenString, accessSecret, "access")
}

// 解析refresh token
func ParseRefreshToken(tokenString string, refreshSecret []byte) (*Claims, error) {
	return parseToken(tokenString, refreshSecret, "refresh")
}

// RemoveBearerPrefix 去除Bearer前缀的通用函数
// 处理各种情况：Bearer token、Bearer  token、Bearer   token等
func RemoveBearerPrefix(tokenString string) string {
	tokenString = strings.TrimSpace(tokenString)
	tokenString = strings.TrimPrefix(tokenString, "Bearer")
	return strings.TrimSpace(tokenString)
}

// 内部解析token方法
func parseToken(tokenString string, secret []byte, expectedType string) (*Claims, error) {
	// 检查并去除Bearer前缀
	tokenString = RemoveBearerPrefix(tokenString)
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenExpired):
			return nil, fmt.Errorf("token已过期")
		case errors.Is(err, jwt.ErrTokenMalformed):
			return nil, fmt.Errorf("token格式错误")
		case errors.Is(err, jwt.ErrTokenSignatureInvalid):
			return nil, fmt.Errorf("token签名无效")
		default:
			return nil, fmt.Errorf("token解析失败: %v", err)
		}
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		// 验证token类型
		if expectedType != "" && claims.TokenType != expectedType {
			return nil, fmt.Errorf("token类型错误: 期望%s，实际%s", expectedType, claims.TokenType)
		}
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}

// 验证access token是否有效
func ValidateAccessToken(tokenString string, accessSecret []byte) bool {
	_, err := ParseAccessToken(tokenString, accessSecret)
	return err == nil
}

// 验证refresh token是否有效
func ValidateRefreshToken(tokenString string, refreshSecret []byte) bool {
	_, err := ParseRefreshToken(tokenString, refreshSecret)
	return err == nil
}

// RefreshTokenService 统一的 refresh token 服务接口
type RefreshTokenService interface {
	RefreshToken(ctx context.Context, refreshToken string) (*RefreshResult, error)
}

// SimpleRefreshResult 简化的刷新结果
type SimpleRefreshResult struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// UnifiedRefreshTokenHandler 统一的 refresh token 处理函数
func UnifiedRefreshTokenHandler(service RefreshTokenService, refreshToken string, ctx context.Context) (*SimpleRefreshResult, error) {
	result, err := service.RefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	return &SimpleRefreshResult{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
	}, nil
}
