package permission

import (
	"context"

	"admin_backend/app/admin/internal/repository/permissionrepo"
	"admin_backend/app/admin/internal/svc"
	"admin_backend/app/admin/internal/types"
	"admin_backend/pkg/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListPermissionLogic struct {
	logx.Logger
	ctx            context.Context
	svcCtx         *svc.ServiceContext
	permissionRepo *permissionrepo.PermissionRepo
}

// 获取权限规则列表
func NewListPermissionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListPermissionLogic {
	return &ListPermissionLogic{
		Logger:         logx.WithContext(ctx),
		ctx:            ctx,
		svcCtx:         svcCtx,
		permissionRepo: permissionrepo.NewPermissionRepo(svcCtx.DB),
	}
}

func (l *ListPermissionLogic) ListPermission(req *types.ListPermissionReq) (resp *types.ListPermissionResp, err error) {
	// 调用仓储层获取权限列表
	list, total, err := l.permissionRepo.PageList(l.ctx, req.Current, req.PageSize, req.Name, req.Code, req.Status)
	if err != nil {
		l.Error("ListPermission Logic PageList err:", err.Error())
		return nil, xerr.NewErrCodeMsg(xerr.DbError, "list permission page err.")
	}

	// 构建返回结果
	permissionList := make([]*types.PermissionInfo, 0)
	for _, p := range list {
		permissionList = append(permissionList, &types.PermissionInfo{
			PermissionID: p.PermissionID,
			Name:         p.Name,
			Code:         p.Code,
			Type:         p.Type,
			Resource:     p.Resource,
			Action:       p.Action,
			ParentID:     p.ParentID,
			Description:  p.Description,
			Status:       p.Status,
			MenuID:       p.MenuID,
		})
	}

	return &types.ListPermissionResp{
		List: permissionList,
		Page: &types.PageResponse{
			Total:           total,
			PageSize:        len(list),
			Current:         req.Current,
			RequestPageSize: req.PageSize,
		},
	}, nil
}
