// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package role

import (
	"context"

	"admin_backend/internal/svc"
	"admin_backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SetUserRolesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 设置用户角色
func NewSetUserRolesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetUserRolesLogic {
	return &SetUserRolesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SetUserRolesLogic) SetUserRoles(req *types.SetUserRolesReq) (resp bool, err error) {
	// todo: add your logic here and delete this line

	return
}
