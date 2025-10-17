package jwt

import (
	"context"
	"fmt"
	"time"

	"admin_backend/pkg/common/xerr"
)

// AddToBlacklist 将access token添加到黑名单
func (m *JWTManager) AddToBlacklist(ctx context.Context, accessToken string) error {
	// 1. 解析access token获取信息
	claims, err := ParseAccessToken(accessToken, m.config.AccessSecret)
	if err != nil {
		m.logger.Errorf("JWTManager.AddToBlacklist ParseAccessToken failed: err=%v", err)
		return xerr.NewErrCodeMsg(xerr.TokenInvalid, "access token无效")
	}

	// 2. 计算token剩余有效时间
	now := time.Now()
	if claims.ExpiresAt.Before(now) {
		// token已过期，无需加入黑名单
		m.logger.Infof("JWTManager.AddToBlacklist token already expired: userID=%s", claims.UserID)
		return nil
	}

	// 3. 计算剩余有效时间作为Redis过期时间
	remainingTime := claims.ExpiresAt.Time.Sub(now)
	if remainingTime <= 0 {
		return nil
	}

	// 4. 将token加入黑名单
	tokenKey := fmt.Sprintf(accessTokenBlacklistPrefix, accessToken)
	tokenData := map[string]string{
		"user_id":        claims.UserID,
		"blacklisted_at": fmt.Sprintf("%d", now.Unix()),
		"expires_at":     fmt.Sprintf("%d", claims.ExpiresAt.Unix()),
	}

	err = m.rdb.HmsetCtx(ctx, tokenKey, tokenData)
	if err != nil {
		m.logger.Errorf("JWTManager.AddToBlacklist Hmset failed: tokenKey=%s, err=%v", tokenKey, err)
		return xerr.NewErrCodeMsg(xerr.ServerError, "添加token到黑名单失败")
	}

	// 5. 设置过期时间（token的剩余有效时间）
	err = m.rdb.ExpireCtx(ctx, tokenKey, int(remainingTime.Seconds()))
	if err != nil {
		m.logger.Errorf("JWTManager.AddToBlacklist Expire failed: tokenKey=%s, err=%v", tokenKey, err)
		return xerr.NewErrCodeMsg(xerr.ServerError, "设置黑名单过期时间失败")
	}

	m.logger.Infof("JWTManager.AddToBlacklist success: userID=%s, remainingTime=%v", claims.UserID, remainingTime)
	return nil
}

// IsInBlacklist 检查access token是否在黑名单中
func (m *JWTManager) IsInBlacklist(ctx context.Context, accessToken string) (bool, error) {
	tokenKey := fmt.Sprintf(accessTokenBlacklistPrefix, accessToken)

	exists, err := m.rdb.ExistsCtx(ctx, tokenKey)
	if err != nil {
		m.logger.Errorf("JWTManager.IsInBlacklist Exists failed: tokenKey=%s, err=%v", tokenKey, err)
		return false, xerr.NewErrCodeMsg(xerr.ServerError, "检查黑名单失败")
	}

	if exists {
		m.logger.Infof("JWTManager.IsInBlacklist token is blacklisted: %s", tokenKey)
		return true, nil
	}

	return false, nil
}

// ValidateAccessTokenWithBlacklist 验证access token并检查黑名单
func (m *JWTManager) ValidateAccessTokenWithBlacklist(ctx context.Context, accessToken string) (*Claims, error) {
	// 1. 基本的JWT验证
	claims, err := ParseAccessToken(accessToken, m.config.AccessSecret)
	if err != nil {
		m.logger.Errorf("JWTManager.ValidateAccessTokenWithBlacklist ParseAccessToken failed: err=%v", err)
		return nil, xerr.NewErrCodeMsg(xerr.TokenInvalid, "access token无效")
	}

	// 2. 检查用户级别黑名单（优先级更高）
	isUserBlacklisted, err := m.IsUserInBlacklist(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}

	if isUserBlacklisted {
		m.logger.Errorf("JWTManager.ValidateAccessTokenWithBlacklist user is blacklisted: userID=%s", claims.UserID)
		return nil, xerr.NewErrCodeMsg(xerr.TokenInvalid, "账户已登出，请重新登录")
	}

	// 3. 检查token级别黑名单
	isBlacklisted, err := m.IsInBlacklist(ctx, accessToken)
	if err != nil {
		return nil, err
	}

	if isBlacklisted {
		m.logger.Errorf("JWTManager.ValidateAccessTokenWithBlacklist token is blacklisted: userID=%s", claims.UserID)
		return nil, xerr.NewErrCodeMsg(xerr.TokenInvalid, "token已失效，请重新登录")
	}

	return claims, nil
}

// AddUserToBlacklist 将用户添加到黑名单（用于全设备登出）
func (m *JWTManager) AddUserToBlacklist(ctx context.Context, userID string, blacklistDuration int64) error {
	userKey := fmt.Sprintf(userBlacklistPrefix, userID)

	userData := map[string]string{
		"user_id":        userID,
		"blacklisted_at": fmt.Sprintf("%d", time.Now().Unix()),
	}

	err := m.rdb.HmsetCtx(ctx, userKey, userData)
	if err != nil {
		m.logger.Errorf("JWTManager.AddUserToBlacklist Hmset failed: userKey=%s, err=%v", userKey, err)
		return xerr.NewErrCodeMsg(xerr.ServerError, "添加用户到黑名单失败")
	}

	// 设置过期时间
	err = m.rdb.ExpireCtx(ctx, userKey, int(blacklistDuration))
	if err != nil {
		m.logger.Errorf("JWTManager.AddUserToBlacklist Expire failed: userKey=%s, err=%v", userKey, err)
		return xerr.NewErrCodeMsg(xerr.ServerError, "设置用户黑名单过期时间失败")
	}

	m.logger.Infof("JWTManager.AddUserToBlacklist success: userID=%s, duration=%ds", userID, blacklistDuration)
	return nil
}

// IsUserInBlacklist 检查用户是否在黑名单中
func (m *JWTManager) IsUserInBlacklist(ctx context.Context, userID string) (bool, error) {
	userKey := fmt.Sprintf(userBlacklistPrefix, userID)

	exists, err := m.rdb.ExistsCtx(ctx, userKey)
	if err != nil {
		m.logger.Errorf("JWTManager.IsUserInBlacklist Exists failed: userKey=%s, err=%v", userKey, err)
		return false, xerr.NewErrCodeMsg(xerr.ServerError, "检查用户黑名单失败")
	}

	if exists {
		m.logger.Infof("JWTManager.IsUserInBlacklist user is blacklisted: userID=%s", userID)
		return true, nil
	}

	return false, nil
}

// LogoutSingleDevice 单设备登出（封装完整流程）
func (m *JWTManager) LogoutSingleDevice(ctx context.Context, accessToken string) error {
	// 1. 将access token加入黑名单（立即失效）
	err := m.AddToBlacklist(ctx, accessToken)
	if err != nil {
		m.logger.Errorf("JWTManager.LogoutSingleDevice AddToBlacklist failed: err=%v", err)
		return err
	}

	// 2. 撤销当前设备的refresh token
	err = m.revokeCurrentDeviceToken(ctx, accessToken)
	if err != nil {
		m.logger.Errorf("JWTManager.LogoutSingleDevice revokeCurrentDeviceToken failed: err=%v", err)
		return err
	}

	m.logger.Infof("JWTManager.LogoutSingleDevice success")
	return nil
}

// LogoutAllDevices 全设备登出（封装完整流程）
func (m *JWTManager) LogoutAllDevices(ctx context.Context, userID string) error {
	// 1. 将用户添加到黑名单（立即失效所有设备的access token）
	err := m.AddUserToBlacklist(ctx, userID, m.config.AccessExpire)
	if err != nil {
		m.logger.Errorf("JWTManager.LogoutAllDevices AddUserToBlacklist failed: userID=%s, err=%v", userID, err)
		return err
	}

	// 2. 撤销用户的所有refresh token
	err = m.RevokeAllUserTokens(ctx, userID)
	if err != nil {
		m.logger.Errorf("JWTManager.LogoutAllDevices RevokeAllUserTokens failed: userID=%s, err=%v", userID, err)
		return err
	}

	m.logger.Infof("JWTManager.LogoutAllDevices success: userID=%s", userID)
	return nil
}
