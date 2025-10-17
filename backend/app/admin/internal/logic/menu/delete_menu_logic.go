package menu

import (
	"context"

	"admin_backend/app/admin/internal/repository/menurepo"
	"admin_backend/app/admin/internal/repository/permissionrepo"
	"admin_backend/app/admin/internal/svc"
	"admin_backend/app/admin/internal/types"
	"admin_backend/pkg/common/xerr"
	"admin_backend/pkg/ent/generated"
	"admin_backend/pkg/ent/generated/permission"
	"admin_backend/pkg/ent/generated/predicate"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteMenuLogic struct {
	logx.Logger
	ctx            context.Context
	svcCtx         *svc.ServiceContext
	menuRepo       *menurepo.MenuRepo
	permissionRepo *permissionrepo.PermissionRepo
}

// 删除菜单
func NewDeleteMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteMenuLogic {
	return &DeleteMenuLogic{
		Logger:         logx.WithContext(ctx),
		ctx:            ctx,
		svcCtx:         svcCtx,
		menuRepo:       menurepo.NewMenuRepo(svcCtx.DB),
		permissionRepo: permissionrepo.NewPermissionRepo(svcCtx.DB),
	}
}

func (l *DeleteMenuLogic) DeleteMenu(req *types.DeleteMenuReq) (resp bool, err error) {

	// 1. 检查菜单是否存在
	menu, err := l.menuRepo.GetByMenuID(l.ctx, req.MenuID)
	if err != nil {
		l.Error("DeleteMenu menuRepo.GetByMenuID err:", err.Error())
		if generated.IsNotFound(err) {
			return false, xerr.NewErrMsg("菜单不存在")
		}
		return false, xerr.NewErrCodeMsg(xerr.DbError, "查询菜单失败")
	}

	where := []predicate.Permission{
		permission.MenuIDEQ(menu.MenuID),
	}
	permission, err := l.permissionRepo.Get(l.ctx, where)
	if err != nil {
		l.Error("DeleteMenu permissionRepo.Get err:", err.Error())
		if generated.IsNotFound(err) {
			return false, xerr.NewErrMsg("菜单不存在")
		}
		return false, xerr.NewErrCodeMsg(xerr.DbError, "查询菜单失败")
	}

	// 2. 软删除权限
	_, err = l.permissionRepo.Delete(l.ctx, permission.PermissionID)
	if err != nil {
		l.Error("DeleteMenu permissionRepo.Delete err:", err.Error())
		return false, xerr.NewErrCodeMsg(xerr.DbError, "删除菜单失败")
	}

	// 3. 软删除菜单
	_, err = l.menuRepo.Delete(l.ctx, menu.MenuID)
	if err != nil {
		l.Error("DeleteMenu menuRepo.Delete err:", err.Error())
		return false, xerr.NewErrCodeMsg(xerr.DbError, "删除菜单失败")
	}

	// 4. casbin 删除菜单权限
	// p     | admin | default | sys_user   | view   | menu   |
	// 2 sys_user
	pm := l.svcCtx.CasbinManager
	err = pm.RemoveFilteredPolicy(2, menu.Code)
	if err != nil {
		l.Errorf("删除菜单权限失败: %v", err)
	}

	return true, nil
}
