// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package auth

import (
	"context"

	"admin_backend/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type LogoutAllLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 用户登出（所有设备）
func NewLogoutAllLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogoutAllLogic {
	return &LogoutAllLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LogoutAllLogic) LogoutAll() (resp bool, err error) {
	// todo: add your logic here and delete this line

	return
}
