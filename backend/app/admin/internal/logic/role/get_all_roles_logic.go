package role

import (
	"context"

	"admin_backend/app/admin/internal/repository/rolerepo"
	"admin_backend/app/admin/internal/svc"
	"admin_backend/app/admin/internal/types"
	"admin_backend/pkg/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAllRolesLogic struct {
	logx.Logger
	ctx      context.Context
	svcCtx   *svc.ServiceContext
	roleRepo *rolerepo.RoleRepo
}

// 获取所有角色列表
func NewGetAllRolesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAllRolesLogic {
	return &GetAllRolesLogic{
		Logger:   logx.WithContext(ctx),
		ctx:      ctx,
		svcCtx:   svcCtx,
		roleRepo: rolerepo.NewRoleRepo(svcCtx.DB),
	}
}

func (l *GetAllRolesLogic) GetAllRoles() (resp *types.GetAllRolesResp, err error) {
	// 获取所有角色列表
	roleList, err := l.roleRepo.List(l.ctx, nil)
	if err != nil {
		l.Errorf("获取角色列表失败: %v", err)
		return nil, xerr.NewErrMsg("获取角色列表失败")
	}

	// 构建响应数据
	list := make([]*types.RoleInfo, 0, len(roleList))
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

	resp = &types.GetAllRolesResp{
		List: list,
	}

	return
}
