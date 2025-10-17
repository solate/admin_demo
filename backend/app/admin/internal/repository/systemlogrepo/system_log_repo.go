package systemlogrepo

import (
	"context"

	"admin_backend/app/admin/internal/svc"
	"admin_backend/pkg/ent/generated"

	"github.com/zeromicro/go-zero/core/logx"
)

type SystemLogLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSystemLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SystemLogLogic {
	return &SystemLogLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// 新增系统日志
func (l *SystemLogLogic) Add(user *generated.SysUser) {
	_, err := l.svcCtx.DB.SystemLog.Create().
		SetUserID(user.UserID).
		Save(l.ctx)
	if err != nil {
		l.Error("Add SystemLog.Create Error:", err.Error())
		return
	}
	return
}
