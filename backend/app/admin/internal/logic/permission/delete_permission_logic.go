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

type DeletePermissionLogic struct {
	logx.Logger
	ctx            context.Context
	svcCtx         *svc.ServiceContext
	permissionRepo *permissionrepo.PermissionRepo
}

// 删除权限规则
func NewDeletePermissionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeletePermissionLogic {
	return &DeletePermissionLogic{
		Logger:         logx.WithContext(ctx),
		ctx:            ctx,
		svcCtx:         svcCtx,
		permissionRepo: permissionrepo.NewPermissionRepo(svcCtx.DB),
	}
}

func (l *DeletePermissionLogic) DeletePermission(req *types.DeletePermissionReq) (resp bool, err error) {
	// 1. 检查权限是否存在
	perm, err := l.permissionRepo.GetByPermissionID(l.ctx, req.PermissionID)
	if err != nil {
		l.Error("DeletePermission permissionRepo.GetByPermissionID err:", err.Error())
		if generated.IsNotFound(err) {
			return false, xerr.NewErrMsg("权限不存在")
		}
		return false, xerr.NewErrCodeMsg(xerr.DbError, "查询权限失败")
	}

	// 2. 软删除权限
	_, err = l.permissionRepo.Delete(l.ctx, perm.PermissionID)
	if err != nil {
		l.Error("DeletePermission permissionRepo.Delete err:", err.Error())
		return false, xerr.NewErrCodeMsg(xerr.DbError, "删除权限失败")
	}

	return true, nil
}
