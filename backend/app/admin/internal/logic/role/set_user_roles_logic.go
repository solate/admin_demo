package role

import (
	"context"

	"admin_backend/app/admin/internal/svc"
	"admin_backend/app/admin/internal/types"
	"admin_backend/pkg/common/contextutil"

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

	tenantCode := contextutil.GetTenantCodeFromCtx(l.ctx)
	// 获取权限管理器
	pm := l.svcCtx.CasbinManager

	// 清除用户在该租户下的所有角色
	err = pm.ClearUserPermissions(req.UserID, tenantCode)
	if err != nil {
		l.Errorf("清除用户角色失败: %v", err)
		return false, err
	}

	// 为用户添加新的角色
	for _, roleCode := range req.RoleCodeList {
		err = pm.AddRoleForUser(req.UserID, roleCode, tenantCode)
		if err != nil {
			l.Errorf("添加用户角色失败: %v", err)
			return false, err
		}
	}

	return true, nil
}
