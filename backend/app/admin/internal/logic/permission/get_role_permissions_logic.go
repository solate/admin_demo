package permission

import (
	"context"

	"admin_backend/app/admin/internal/svc"
	"admin_backend/app/admin/internal/types"
	"admin_backend/pkg/common/contextutil"

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
	tenantCode := contextutil.GetTenantCodeFromCtx(l.ctx)
	// 获取权限管理器
	pm := l.svcCtx.CasbinManager

	// 获取角色在该租户下的所有权限编码
	casbinPermissionList, err := pm.GetRolePermissions(req.RoleCode, tenantCode)
	if err != nil {
		l.Errorf("获取角色权限列表失败: %v", err)
		return nil, err
	}

	// 如果角色没有权限，则直接返回空列表
	if len(casbinPermissionList) == 0 {
		return &types.GetRolePermissionsResp{
			List: make([]*types.Permission, 0),
		}, nil
	}

	list := make([]*types.Permission, 0, len(casbinPermissionList))
	for _, permissionInfo := range casbinPermissionList {
		list = append(list, &types.Permission{
			Resource: permissionInfo[2],
			Action:   permissionInfo[3],
			Type:     permissionInfo[4],
		})
	}

	resp = &types.GetRolePermissionsResp{
		List: list,
	}

	return
}
