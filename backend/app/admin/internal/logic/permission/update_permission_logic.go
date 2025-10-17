package permission

import (
	"context"

	"admin_backend/app/admin/internal/repository/permissionrepo"
	"admin_backend/app/admin/internal/svc"
	"admin_backend/app/admin/internal/types"
	"admin_backend/pkg/common/xerr"
	"admin_backend/pkg/ent/generated"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdatePermissionLogic struct {
	logx.Logger
	ctx            context.Context
	svcCtx         *svc.ServiceContext
	permissionRepo *permissionrepo.PermissionRepo
}

// 更新权限规则
func NewUpdatePermissionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatePermissionLogic {
	return &UpdatePermissionLogic{
		Logger:         logx.WithContext(ctx),
		ctx:            ctx,
		svcCtx:         svcCtx,
		permissionRepo: permissionrepo.NewPermissionRepo(svcCtx.DB),
	}
}

func (l *UpdatePermissionLogic) UpdatePermission(req *types.UpdatePermissionReq) (resp bool, err error) {
	// 1. 检查权限是否存在
	permission, err := l.permissionRepo.GetByPermissionID(l.ctx, req.PermissionID)
	if err != nil {
		l.Error("UpdatePermission permissionRepo.GetByPermissionID err:", err.Error())
		if generated.IsNotFound(err) {
			return false, xerr.NewErrMsg("权限不存在")
		}
		return false, xerr.NewErrCodeMsg(xerr.DbError, "查询权限失败")
	}

	// 2. 更新权限信息
	if req.Name != "" {
		permission.Name = req.Name
	}
	if req.Type != "" {
		permission.Type = req.Type
	}
	if req.Resource != "" {
		permission.Resource = req.Resource
	}
	if req.Action != "" {
		permission.Action = req.Action
	}
	if req.ParentID != "" {
		permission.ParentID = req.ParentID
	}
	if req.Description != "" {
		permission.Description = req.Description
	}
	if req.Status != 0 {
		permission.Status = req.Status
	}
	if req.MenuID != "" {
		permission.MenuID = req.MenuID
	}

	_, err = l.permissionRepo.Update(l.ctx, permission)
	if err != nil {
		l.Error("UpdatePermission Update err:", err.Error())
		return false, xerr.NewErrCodeMsg(xerr.DbError, "更新权限失败")
	}

	return true, nil
}
