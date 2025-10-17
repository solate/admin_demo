package auth

import (
	"context"

	"admin_backend/app/admin/internal/svc"
	"admin_backend/app/admin/internal/types"
	"admin_backend/pkg/common/contextutil"
	"admin_backend/pkg/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetActiveDevicesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取当前用户活跃设备数量
func NewGetActiveDevicesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetActiveDevicesLogic {
	return &GetActiveDevicesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetActiveDevicesLogic) GetActiveDevices() (resp *types.GetActiveDevicesResp, err error) {
	// 1. 获取当前用户ID
	userID := contextutil.GetUserIDFromCtx(l.ctx)

	// 2. 从JWT Manager获取活跃设备数量
	count, err := l.svcCtx.JWTManager.GetUserActiveTokens(l.ctx, userID)
	if err != nil {
		l.Error("GetActiveDevices GetUserActiveTokens err:", err.Error())
		return nil, xerr.NewErrCodeMsg(xerr.ServerError, "获取活跃设备数量失败")
	}

	return &types.GetActiveDevicesResp{
		ActiveDevices: count,
	}, nil
}
