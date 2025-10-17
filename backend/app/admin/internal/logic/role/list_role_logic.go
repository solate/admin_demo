package role

import (
	"context"

	"admin_backend/app/admin/internal/repository/rolerepo"
	"admin_backend/app/admin/internal/svc"
	"admin_backend/app/admin/internal/types"
	"admin_backend/pkg/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListRoleLogic struct {
	logx.Logger
	ctx      context.Context
	svcCtx   *svc.ServiceContext
	roleRepo *rolerepo.RoleRepo
}

// 获取角色列表
func NewListRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListRoleLogic {
	return &ListRoleLogic{
		Logger:   logx.WithContext(ctx),
		ctx:      ctx,
		svcCtx:   svcCtx,
		roleRepo: rolerepo.NewRoleRepo(svcCtx.DB),
	}
}

func (l *ListRoleLogic) ListRole(req *types.RoleListReq) (resp *types.RoleListResp, err error) {
	// 1. 获取角色列表
	roles, total, err := l.roleRepo.PageList(l.ctx, req.Current, req.PageSize, req.Name, req.Code, req.Status)
	if err != nil {
		l.Error("ListRole List err:", err.Error())
		return nil, xerr.NewErrMsg("获取角色列表失败")
	}

	// 2. 构建响应数据
	list := make([]*types.RoleInfo, 0, len(roles))
	for _, role := range roles {
		list = append(list, &types.RoleInfo{
			RoleID:      role.RoleID,
			Name:        role.Name,
			Code:        role.Code,
			Description: role.Description,
			Status:      role.Status,
			Sort:        role.Sort,
			CreatedAt:   role.CreatedAt,
		})
	}

	// 3. 返回结果
	return &types.RoleListResp{
		Page: &types.PageResponse{
			Total:           total,
			PageSize:        len(list),
			RequestPageSize: req.PageSize,
			Current:         req.Current,
		},
		List: list,
	}, nil
}
