package permission

import (
	"context"

	"admin_backend/app/admin/internal/repository/permissionrepo"
	"admin_backend/app/admin/internal/svc"
	"admin_backend/app/admin/internal/types"
	"admin_backend/pkg/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPermissionLogic struct {
	logx.Logger
	ctx            context.Context
	svcCtx         *svc.ServiceContext
	permissionRepo *permissionrepo.PermissionRepo
}

// 获取权限规则详情
func NewGetPermissionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPermissionLogic {
	return &GetPermissionLogic{
		Logger:         logx.WithContext(ctx),
		ctx:            ctx,
		svcCtx:         svcCtx,
		permissionRepo: permissionrepo.NewPermissionRepo(svcCtx.DB),
	}
}

func (l *GetPermissionLogic) GetPermission(req *types.GetPermissionReq) (resp *types.PermissionInfo, err error) {
	// 查询权限信息
	permission, err := l.permissionRepo.GetByPermissionID(l.ctx, req.PermissionID)
	if err != nil {
		l.Error("GetPermission l.permissionRepo.GetByPermissionID err: ", err.Error())
		return nil, xerr.NewErrCodeMsg(xerr.DbError, "获取权限信息失败")
	}

	resp = &types.PermissionInfo{
		PermissionID: permission.PermissionID,
		Name:         permission.Name,
		Code:         permission.Code,
		Type:         permission.Type,
		Resource:     permission.Resource,
		Action:       permission.Action,
		ParentID:     permission.ParentID,
		Description:  permission.Description,
		Status:       permission.Status,
		MenuID:       permission.MenuID,
		CreatedAt:    permission.CreatedAt,
	}

	return
}
