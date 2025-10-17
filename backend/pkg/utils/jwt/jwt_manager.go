package jwt

import (
	"context"

	"admin_backend/pkg/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

const (
	// refresh token存储key前缀
	refreshTokenPrefix = "refresh_token:%s"
	// 用户refresh token列表key前缀
	userRefreshTokensPrefix = "user_refresh_tokens:%s"
	// access token黑名单key前缀
	accessTokenBlacklistPrefix = "blacklist_access_token:%s"
	// 用户级别黑名单key前缀（用于全设备登出）
	userBlacklistPrefix = "blacklist_user:%s"
	// 默认最大设备数
	defaultMaxDevices = 5
)

// JWTManager 统一的JWT管理器，包含所有JWT相关功能
type JWTManager struct {
	config JWTConfig
	rdb    *redis.Redis
	logger logx.Logger
}

// LoginResult 登录结果
type LoginResult struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// RefreshResult 刷新结果
type RefreshResult struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// NewJWTManager 创建统一的JWT管理器
func NewJWTManager(config JWTConfig, rdb *redis.Redis) *JWTManager {
	return &JWTManager{
		config: config,
		rdb:    rdb,
		logger: logx.WithContext(context.Background()),
	}
}

// Login 用户登录，生成并存储token对
func (m *JWTManager) Login(ctx context.Context, userID, tenantCode, roleCode string, source int) (*LoginResult, error) {
	// 1. 生成token对
	tokenPair, err := m.generateTokenPair(userID, tenantCode, roleCode, source)
	if err != nil {
		m.logger.Errorf("JWTManager.Login generateTokenPair failed: userID=%s, err=%v", userID, err)
		return nil, xerr.NewErrCodeMsg(xerr.ServerError, "生成token失败")
	}

	// 2. 存储refresh token到Redis
	err = m.storeRefreshToken(ctx, userID, tokenPair.TokenID, tokenPair.RefreshToken)
	if err != nil {
		m.logger.Errorf("JWTManager.Login storeRefreshToken failed: userID=%s, tokenID=%s, err=%v", userID, tokenPair.TokenID, err)
		return nil, xerr.NewErrCodeMsg(xerr.ServerError, "存储refresh token失败")
	}

	// 3. 限制用户设备数量
	err = m.limitUserTokens(ctx, userID, defaultMaxDevices)
	if err != nil {
		// 限制设备失败不影响登录，只记录日志
		m.logger.Errorf("JWTManager.Login limitUserTokens failed: userID=%s, err=%v", userID, err)
	}

	m.logger.Infof("JWTManager.Login success: userID=%s, source=%d", userID, source)
	return &LoginResult{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	}, nil
}

// RefreshToken 刷新token
func (m *JWTManager) RefreshToken(ctx context.Context, refreshToken string) (*RefreshResult, error) {
	// 1. 验证JWT格式和签名
	claims, err := ParseRefreshToken(refreshToken, m.config.RefreshSecret)
	if err != nil {
		m.logger.Errorf("JWTManager.RefreshToken ParseRefreshToken failed: err=%v", err)
		return nil, xerr.NewErrCodeMsg(xerr.TokenInvalid, "refresh token格式无效")
	}

	// 2. 验证refresh token是否在Redis白名单中
	var oldTokenID string
	if claims.TokenID != "" {
		// 去除Bearer前缀，因为Redis中存储的是纯token
		cleanToken := RemoveBearerPrefix(refreshToken)
		valid, userID, err := m.verifyRefreshToken(ctx, claims.TokenID, cleanToken)
		if err != nil || !valid {
			m.logger.Errorf("JWTManager.RefreshToken verifyRefreshToken failed: tokenID=%s, userID=%s, valid=%v, err=%v", claims.TokenID, claims.UserID, valid, err)
			return nil, xerr.NewErrCodeMsg(xerr.TokenInvalid, "refresh token已失效")
		}

		// 确保userID一致
		if userID != claims.UserID {
			m.logger.Errorf("JWTManager.RefreshToken userID mismatch: expected=%s, actual=%s", claims.UserID, userID)
			return nil, xerr.NewErrCodeMsg(xerr.TokenInvalid, "token用户信息不匹配")
		}

		// 暂存旧token ID，在新token存储成功后再撤销
		oldTokenID = claims.TokenID
	}

	// 3. 生成新的token对
	newTokenPair, err := m.generateTokenPair(claims.UserID, claims.TenantCode, claims.RoleCode, claims.Source)
	if err != nil {
		m.logger.Errorf("JWTManager.RefreshToken generateTokenPair failed: userID=%s, err=%v", claims.UserID, err)
		return nil, xerr.NewErrCodeMsg(xerr.ServerError, "生成新token失败")
	}

	// 4. 存储新的refresh token
	err = m.storeRefreshToken(ctx, claims.UserID, newTokenPair.TokenID, newTokenPair.RefreshToken)
	if err != nil {
		m.logger.Errorf("JWTManager.RefreshToken storeRefreshToken failed: userID=%s, tokenID=%s, err=%v", claims.UserID, newTokenPair.TokenID, err)
		return nil, xerr.NewErrCodeMsg(xerr.ServerError, "存储新refresh token失败")
	}

	// 5. 新token存储成功后，撤销旧的refresh token
	if oldTokenID != "" {
		err = m.revokeRefreshToken(ctx, claims.UserID, oldTokenID)
		if err != nil {
			// 撤销旧token失败不影响刷新流程，只记录错误
			m.logger.Errorf("JWTManager.RefreshToken revokeRefreshToken failed: userID=%s, oldTokenID=%s, err=%v", claims.UserID, oldTokenID, err)
		} else {
			m.logger.Infof("JWTManager.RefreshToken revoked old token: userID=%s, oldTokenID=%s", claims.UserID, oldTokenID)
		}
	}

	m.logger.Infof("JWTManager.RefreshToken success: userID=%s, newTokenID=%s", claims.UserID, newTokenPair.TokenID)
	return &RefreshResult{
		AccessToken:  newTokenPair.AccessToken,
		RefreshToken: newTokenPair.RefreshToken,
	}, nil
}

// RevokeToken 撤销指定的refresh token
func (m *JWTManager) RevokeToken(ctx context.Context, userID, tokenID string) error {
	err := m.revokeRefreshToken(ctx, userID, tokenID)
	if err != nil {
		m.logger.Errorf("JWTManager.RevokeToken failed: userID=%s, tokenID=%s, err=%v", userID, tokenID, err)
		return xerr.NewErrCodeMsg(xerr.ServerError, "撤销token失败")
	}
	m.logger.Infof("JWTManager.RevokeToken success: userID=%s, tokenID=%s", userID, tokenID)
	return nil
}

// RevokeAllUserTokens 撤销用户的所有token（登出所有设备）
func (m *JWTManager) RevokeAllUserTokens(ctx context.Context, userID string) error {
	err := m.revokeAllUserTokens(ctx, userID)
	if err != nil {
		m.logger.Errorf("JWTManager.RevokeAllUserTokens failed: userID=%s, err=%v", userID, err)
		return xerr.NewErrCodeMsg(xerr.ServerError, "撤销所有token失败")
	}
	m.logger.Infof("JWTManager.RevokeAllUserTokens success: userID=%s", userID)
	return nil
}

// GetTokenExpireTime 获取token过期时间配置
func (m *JWTManager) GetTokenExpireTime() (accessExpire, refreshExpire int64) {
	return m.config.AccessExpire, m.config.RefreshExpire
}
