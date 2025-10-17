package role

import (
	"context"

	"admin_backend/app/admin/internal/repository/rolerepo"
	"admin_backend/app/admin/internal/svc"
	"admin_backend/app/admin/internal/types"
	"admin_backend/pkg/common/contextutil"
	"admin_backend/pkg/common/xerr"
	"admin_backend/pkg/ent/generated"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteRoleLogic struct {
	logx.Logger
	ctx      context.Context
	svcCtx   *svc.ServiceContext
	roleRepo *rolerepo.RoleRepo
}

// 删除角色
func NewDeleteRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteRoleLogic {
	return &DeleteRoleLogic{
		Logger:   logx.WithContext(ctx),
		ctx:      ctx,
		svcCtx:   svcCtx,
		roleRepo: rolerepo.NewRoleRepo(svcCtx.DB),
	}
}

func (l *DeleteRoleLogic) DeleteRole(req *types.DeleteRoleReq) (resp bool, err error) {
	// 1. 检查角色是否存在
	role, err := l.roleRepo.GetByRoleID(l.ctx, req.RoleID)
	if err != nil {
		l.Error("DeleteRole roleRepo.GetByRoleID err:", err.Error())
		if generated.IsNotFound(err) {
			return false, xerr.NewErrMsg("角色不存在")
		}
		return false, xerr.NewErrCodeMsg(xerr.DbError, "查询角色失败")
	}

	// 2. 判断是否具有权限
	pm := l.svcCtx.CasbinManager
	tenantCode := contextutil.GetTenantCodeFromCtx(l.ctx)
	casbinPermissionList, err := pm.GetRolePermissions(role.Code, tenantCode)
	if err != nil {
		l.Errorf("获取角色权限列表失败: %v", err)
		return false, xerr.NewErrCodeMsg(xerr.DbError, "查询角色失败")
	}
	if len(casbinPermissionList) != 0 {
		return false, xerr.NewErrCodeMsg(xerr.ForbiddenError, "角色存在权限，不能删除")
	}

	// 2. 软删除角色
	_, err = l.roleRepo.Delete(l.ctx, role.RoleID)
	if err != nil {
		l.Error("DeleteRole roleRepo.Update err:", err.Error())
		return false, xerr.NewErrCodeMsg(xerr.DbError, "删除角色失败")
	}

	return true, nil
}
