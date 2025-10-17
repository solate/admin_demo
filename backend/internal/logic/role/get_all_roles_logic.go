// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package role

import (
	"context"

	"admin_backend/internal/svc"
	"admin_backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAllRolesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取所有角色列表
func NewGetAllRolesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAllRolesLogic {
	return &GetAllRolesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAllRolesLogic) GetAllRoles() (resp *types.GetAllRolesResp, err error) {
	// todo: add your logic here and delete this line

	return
}
