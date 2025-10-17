package jwt

import (
	"context"
	"fmt"
	"time"

	"admin_backend/pkg/utils/idgen"
)

// storeRefreshToken 存储refresh token到Redis
func (m *JWTManager) storeRefreshToken(ctx context.Context, userID, tokenID, refreshToken string) error {
	m.logger.Infof("storeRefreshToken: userID=%s, tokenID=%s, expire=%ds", userID, tokenID, m.config.RefreshExpire)

	// 1. 存储refresh token本身
	tokenKey := fmt.Sprintf(refreshTokenPrefix, tokenID)
	tokenData := map[string]string{
		"user_id":       userID,
		"refresh_token": refreshToken,
		"created_at":    fmt.Sprintf("%d", time.Now().Unix()),
	}

	err := m.rdb.HmsetCtx(ctx, tokenKey, tokenData)
	if err != nil {
		m.logger.Errorf("storeRefreshToken Hmset failed: tokenKey=%s, err=%v", tokenKey, err)
		return fmt.Errorf("存储token到Redis失败: %w", err)
	}

	// 设置过期时间
	err = m.rdb.ExpireCtx(ctx, tokenKey, int(m.config.RefreshExpire))
	if err != nil {
		m.logger.Errorf("storeRefreshToken Expire failed: tokenKey=%s, err=%v", tokenKey, err)
		// 如果设置过期失败，清理已创建的数据
		m.rdb.DelCtx(ctx, tokenKey)
		return fmt.Errorf("设置token过期时间失败: %w", err)
	}

	// 2. 将token ID加入用户的token列表（支持多设备登录）
	userTokensKey := fmt.Sprintf(userRefreshTokensPrefix, userID)
	_, err = m.rdb.SaddCtx(ctx, userTokensKey, tokenID)
	if err != nil {
		m.logger.Errorf("storeRefreshToken Sadd failed: userTokensKey=%s, tokenID=%s, err=%v", userTokensKey, tokenID, err)
		// 如果添加到用户列表失败，清理已创建的token数据
		m.rdb.DelCtx(ctx, tokenKey)
		return fmt.Errorf("添加token到用户列表失败: %w", err)
	}

	// 设置用户token列表的过期时间
	err = m.rdb.ExpireCtx(ctx, userTokensKey, int(m.config.RefreshExpire))
	if err != nil {
		m.logger.Errorf("storeRefreshToken user tokens Expire failed: userTokensKey=%s, err=%v", userTokensKey, err)
		// 用户token列表过期时间设置失败不是致命错误，只记录日志
		m.logger.Errorf("设置用户token列表过期时间失败，但token存储成功: userID=%s", userID)
	}

	// 验证存储是否成功
	exists, err := m.rdb.ExistsCtx(ctx, tokenKey)
	if err != nil || !exists {
		m.logger.Errorf("storeRefreshToken verification failed: tokenKey=%s, exists=%v, err=%v", tokenKey, exists, err)
		return fmt.Errorf("token存储验证失败")
	}

	m.logger.Infof("storeRefreshToken success: userID=%s, tokenID=%s", userID, tokenID)
	return nil
}

// verifyRefreshToken 验证refresh token是否有效
func (m *JWTManager) verifyRefreshToken(ctx context.Context, tokenID, refreshToken string) (bool, string, error) {
	tokenKey := fmt.Sprintf(refreshTokenPrefix, tokenID)

	// 检查token是否存在
	exists, err := m.rdb.ExistsCtx(ctx, tokenKey)
	if err != nil {
		m.logger.Errorf("verifyRefreshToken Exists failed: tokenKey=%s, err=%v", tokenKey, err)
		return false, "", err
	}

	if !exists {
		m.logger.Errorf("verifyRefreshToken token not exists: tokenID=%s", tokenID)
		return false, "", fmt.Errorf("refresh token不存在或已过期")
	}

	// 获取存储的refresh token和用户ID
	data, err := m.rdb.HmgetCtx(ctx, tokenKey, "user_id", "refresh_token")
	if err != nil {
		m.logger.Errorf("verifyRefreshToken Hmget failed: tokenKey=%s, err=%v", tokenKey, err)
		return false, "", err
	}

	if len(data) != 2 || data[1] != refreshToken {
		m.logger.Errorf("verifyRefreshToken token mismatch: tokenID=%s", tokenID)
		return false, "", fmt.Errorf("refresh token无效")
	}

	m.logger.Infof("verifyRefreshToken success: tokenID=%s, userID=%s", tokenID, data[0])
	return true, data[0], nil
}

// revokeRefreshToken 撤销指定的refresh token
func (m *JWTManager) revokeRefreshToken(ctx context.Context, userID, tokenID string) error {
	m.logger.Infof("revokeRefreshToken: userID=%s, tokenID=%s", userID, tokenID)

	// 1. 删除token
	tokenKey := fmt.Sprintf(refreshTokenPrefix, tokenID)
	_, err := m.rdb.DelCtx(ctx, tokenKey)
	if err != nil {
		m.logger.Errorf("revokeRefreshToken Del token failed: tokenKey=%s, err=%v", tokenKey, err)
		return err
	}

	// 2. 从用户token列表中移除
	userTokensKey := fmt.Sprintf(userRefreshTokensPrefix, userID)
	_, err = m.rdb.SremCtx(ctx, userTokensKey, tokenID)
	if err != nil {
		m.logger.Errorf("revokeRefreshToken Srem failed: userTokensKey=%s, tokenID=%s, err=%v", userTokensKey, tokenID, err)
		return err
	}

	m.logger.Infof("revokeRefreshToken success: userID=%s, tokenID=%s", userID, tokenID)
	return nil
}

// revokeAllUserTokens 撤销用户的所有refresh token（用于登出所有设备）
func (m *JWTManager) revokeAllUserTokens(ctx context.Context, userID string) error {
	m.logger.Infof("revokeAllUserTokens: userID=%s", userID)

	userTokensKey := fmt.Sprintf(userRefreshTokensPrefix, userID)

	// 获取用户的所有token ID
	tokenIDs, err := m.rdb.SmembersCtx(ctx, userTokensKey)
	if err != nil {
		m.logger.Errorf("revokeAllUserTokens Smembers failed: userTokensKey=%s, err=%v", userTokensKey, err)
		return err
	}

	// 删除所有token
	for _, tokenID := range tokenIDs {
		tokenKey := fmt.Sprintf(refreshTokenPrefix, tokenID)
		_, err := m.rdb.DelCtx(ctx, tokenKey)
		if err != nil {
			m.logger.Errorf("revokeAllUserTokens Del token failed: tokenKey=%s, err=%v", tokenKey, err)
		}
	}

	// 清空用户token列表
	_, err = m.rdb.DelCtx(ctx, userTokensKey)
	if err != nil {
		m.logger.Errorf("revokeAllUserTokens Del user tokens failed: userTokensKey=%s, err=%v", userTokensKey, err)
		return err
	}

	m.logger.Infof("revokeAllUserTokens success: userID=%s, revokedCount=%d", userID, len(tokenIDs))
	return nil
}

// getUserActiveTokens 获取用户当前有效的token数量
func (m *JWTManager) getUserActiveTokens(ctx context.Context, userID string) (int, error) {
	userTokensKey := fmt.Sprintf(userRefreshTokensPrefix, userID)
	count, err := m.rdb.ScardCtx(ctx, userTokensKey)
	if err != nil {
		m.logger.Errorf("getUserActiveTokens Scard failed: userTokensKey=%s, err=%v", userTokensKey, err)
		return 0, err
	}
	return int(count), nil
}

// generateTokenPair 生成token对
func (m *JWTManager) generateTokenPair(userID, tenantCode, roleCode string, source int) (*TokenPair, error) {
	// 生成refresh token的唯一ID
	tokenID, err := idgen.GenerateUUID()
	if err != nil {
		return nil, fmt.Errorf("生成token ID失败: %w", err)
	}

	// 生成access token (也包含tokenID，用于识别设备)
	accessToken, err := generateTokenWithID(userID, tenantCode, roleCode, source, "access", tokenID, m.config.AccessSecret, m.config.AccessExpire)
	if err != nil {
		return nil, fmt.Errorf("生成access token失败: %w", err)
	}

	// 生成refresh token (带tokenID)
	refreshToken, err := generateTokenWithID(userID, tenantCode, roleCode, source, "refresh", tokenID, m.config.RefreshSecret, m.config.RefreshExpire)
	if err != nil {
		return nil, fmt.Errorf("生成refresh token失败: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenID:      tokenID,
	}, nil
}
