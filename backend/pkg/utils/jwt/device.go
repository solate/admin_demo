package jwt

import (
	"context"
	"fmt"

	"admin_backend/pkg/common/xerr"
)

// revokeCurrentDeviceToken 撤销当前设备的refresh token（内部方法）
func (m *JWTManager) revokeCurrentDeviceToken(ctx context.Context, accessToken string) error {
	// 1. 解析access token获取tokenID
	claims, err := ParseAccessToken(accessToken, m.config.AccessSecret)
	if err != nil {
		m.logger.Errorf("JWTManager.revokeCurrentDeviceToken ParseAccessToken failed: err=%v", err)
		return xerr.NewErrCodeMsg(xerr.TokenInvalid, "access token无效")
	}

	// 2. 检查是否有tokenID
	if claims.TokenID == "" {
		m.logger.Errorf("JWTManager.revokeCurrentDeviceToken no tokenID in access token: userID=%s", claims.UserID)
		// 兼容旧版本token，如果没有tokenID就撤销所有token
		return m.RevokeAllUserTokens(ctx, claims.UserID)
	}

	// 3. 撤销当前设备的refresh token
	err = m.revokeRefreshToken(ctx, claims.UserID, claims.TokenID)
	if err != nil {
		m.logger.Errorf("JWTManager.revokeCurrentDeviceToken revokeRefreshToken failed: userID=%s, tokenID=%s, err=%v", claims.UserID, claims.TokenID, err)
		return xerr.NewErrCodeMsg(xerr.ServerError, "撤销当前设备token失败")
	}

	m.logger.Infof("JWTManager.revokeCurrentDeviceToken success: userID=%s, tokenID=%s", claims.UserID, claims.TokenID)
	return nil
}

// GetUserActiveTokens 获取用户当前有效的token数量
func (m *JWTManager) GetUserActiveTokens(ctx context.Context, userID string) (int, error) {
	count, err := m.getUserActiveTokens(ctx, userID)
	if err != nil {
		m.logger.Errorf("JWTManager.GetUserActiveTokens failed: userID=%s, err=%v", userID, err)
		return 0, xerr.NewErrCodeMsg(xerr.ServerError, "获取活跃token数量失败")
	}
	return count, nil
}

// limitUserTokens 限制用户的最大token数量（用于限制登录设备数）
func (m *JWTManager) limitUserTokens(ctx context.Context, userID string, maxTokens int) error {
	userTokensKey := fmt.Sprintf(userRefreshTokensPrefix, userID)

	// 获取当前token数量
	count, err := m.rdb.ScardCtx(ctx, userTokensKey)
	if err != nil {
		m.logger.Errorf("limitUserTokens Scard failed: userTokensKey=%s, err=%v", userTokensKey, err)
		return err
	}

	// 如果超过限制，删除最旧的token
	if count > int64(maxTokens) {
		// 获取所有token ID
		tokenIDs, err := m.rdb.SmembersCtx(ctx, userTokensKey)
		if err != nil {
			m.logger.Errorf("limitUserTokens Smembers failed: userTokensKey=%s, err=%v", userTokensKey, err)
			return err
		}

		// 按创建时间排序，删除最旧的
		excessCount := int(count) - maxTokens
		m.logger.Infof("limitUserTokens: userID=%s, currentCount=%d, maxTokens=%d, excessCount=%d", userID, count, maxTokens, excessCount)

		for i := 0; i < excessCount && i < len(tokenIDs); i++ {
			err = m.revokeRefreshToken(ctx, userID, tokenIDs[i])
			if err != nil {
				m.logger.Errorf("limitUserTokens revokeRefreshToken failed: userID=%s, tokenID=%s, err=%v", userID, tokenIDs[i], err)
				return err
			}
		}
	}

	return nil
}
