// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package permission

import (
	"context"

	"admin_backend/internal/svc"
	"admin_backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRolePermissionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取角色权限列表
func NewGetRolePermissionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRolePermissionsLogic {
	return &GetRolePermissionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetRolePermissionsLogic) GetRolePermissions(req *types.GetRolePermissionsReq) (resp *types.GetRolePermissionsResp, err error) {
	// todo: add your logic here and delete this line

	return
}
