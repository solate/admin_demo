package role

import (
	"context"

	"admin_backend/app/admin/internal/repository/rolerepo"
	"admin_backend/app/admin/internal/svc"
	"admin_backend/app/admin/internal/types"
	"admin_backend/pkg/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRoleLogic struct {
	logx.Logger
	ctx      context.Context
	svcCtx   *svc.ServiceContext
	roleRepo *rolerepo.RoleRepo
}

// 获取角色详情
func NewGetRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRoleLogic {
	return &GetRoleLogic{
		Logger:   logx.WithContext(ctx),
		ctx:      ctx,
		svcCtx:   svcCtx,
		roleRepo: rolerepo.NewRoleRepo(svcCtx.DB),
	}
}

func (l *GetRoleLogic) GetRole(req *types.GetRoleReq) (resp *types.RoleInfo, err error) {
	// 查询角色信息
	role, err := l.roleRepo.GetByRoleID(l.ctx, req.RoleID)
	if err != nil {
		l.Error("GetRole l.roleRepo.GetByRoleID err: ", err.Error())
		return nil, xerr.NewErrCodeMsg(xerr.DbError, "获取角色信息失败")
	}

	resp = &types.RoleInfo{
		RoleID:      role.RoleID,
		Name:        role.Name,
		Code:        role.Code,
		Description: role.Description,
		Status:      role.Status,
		Sort:        role.Sort,
		CreatedAt:   role.CreatedAt,
	}

	return
}
