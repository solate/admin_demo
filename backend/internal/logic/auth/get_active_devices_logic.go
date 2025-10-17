// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package auth

import (
	"context"

	"admin_backend/internal/svc"
	"admin_backend/internal/types"

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
	// todo: add your logic here and delete this line

	return
}
