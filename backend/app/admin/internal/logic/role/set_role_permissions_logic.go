package role

import (
	"context"

	"admin_backend/app/admin/internal/svc"
	"admin_backend/app/admin/internal/types"
	"admin_backend/pkg/common/contextutil"

	"github.com/zeromicro/go-zero/core/logx"
)

type SetRolePermissionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 设置角色权限
func NewSetRolePermissionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetRolePermissionsLogic {
	return &SetRolePermissionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SetRolePermissionsLogic) SetRolePermissions(req *types.SetRolePermissionsReq) (resp bool, err error) {
	tenantCode := contextutil.GetTenantCodeFromCtx(l.ctx)
	// 获取权限管理器
	pm := l.svcCtx.CasbinManager

	// 清除角色在该租户下的所有权限
	err = pm.ClearRolePermissions(req.RoleCode, tenantCode)
	if err != nil {
		l.Errorf("清除角色权限失败: %v", err)
		return false, err
	}

	// 为角色添加新的权限
	var policies [][]string
	for _, item := range req.PermissionList {
		policies = append(policies, []string{req.RoleCode, tenantCode, item.Resource, item.Action, item.Type})
	}
	err = pm.BatchAddPermissions(policies)
	if err != nil {
		l.Errorf("添加角色权限失败: %v", err)
		return false, err
	}

	return true, nil
}
