package role

import (
	"context"

	"admin_backend/app/admin/internal/repository/rolerepo"
	"admin_backend/app/admin/internal/svc"
	"admin_backend/app/admin/internal/types"
	"admin_backend/pkg/common/contextutil"
	"admin_backend/pkg/ent/generated/predicate"
	"admin_backend/pkg/ent/generated/role"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserRolesLogic struct {
	logx.Logger
	ctx      context.Context
	svcCtx   *svc.ServiceContext
	roleRepo *rolerepo.RoleRepo
}

// 获取用户角色列表
func NewGetUserRolesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserRolesLogic {
	return &GetUserRolesLogic{
		Logger:   logx.WithContext(ctx),
		ctx:      ctx,
		svcCtx:   svcCtx,
		roleRepo: rolerepo.NewRoleRepo(svcCtx.DB),
	}
}

func (l *GetUserRolesLogic) GetUserRoles(req *types.GetUserRolesReq) (resp *types.GetUserRolesResp, err error) {
	tenantCode := contextutil.GetTenantCodeFromCtx(l.ctx)
	// 获取权限管理器
	pm := l.svcCtx.CasbinManager

	// 获取用户在该租户下的所有角色编码
	roleCodeList, err := pm.GetRolesForUser(req.UserID, tenantCode)
	if err != nil {
		l.Errorf("获取用户角色列表失败: %v", err)
		return nil, err
	}

	// 如果用户没有角色，则直接返回空列表
	if len(roleCodeList) == 0 {
		return &types.GetUserRolesResp{
			List: make([]*types.RoleInfo, 0),
		}, nil
	}

	// 根据角色编码获取完整的角色信息
	where := []predicate.Role{
		role.CodeIn(roleCodeList...),
	}
	roleList, err := l.roleRepo.List(l.ctx, where)
	if err != nil {
		l.Errorf("获取角色信息失败, code: %v, err: %v", roleCodeList, err)
		return nil, err
	}

	list := make([]*types.RoleInfo, 0, len(roleCodeList))
	for _, roleInfo := range roleList {
		list = append(list, &types.RoleInfo{
			RoleID:      roleInfo.RoleID,
			Name:        roleInfo.Name,
			Code:        roleInfo.Code,
			Description: roleInfo.Description,
			Status:      roleInfo.Status,
			Sort:        roleInfo.Sort,
			CreatedAt:   roleInfo.CreatedAt,
		})
	}

	resp = &types.GetUserRolesResp{
		List: list,
	}

	return
}
