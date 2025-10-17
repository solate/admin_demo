package menu

import (
	"context"

	"admin_backend/app/admin/internal/repository/menurepo"
	"admin_backend/app/admin/internal/svc"
	"admin_backend/app/admin/internal/types"
	"admin_backend/pkg/common/xerr"
	"admin_backend/pkg/ent/generated"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateMenuLogic struct {
	logx.Logger
	ctx      context.Context
	svcCtx   *svc.ServiceContext
	menuRepo *menurepo.MenuRepo
}

// 更新菜单
func NewUpdateMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateMenuLogic {
	return &UpdateMenuLogic{
		Logger:   logx.WithContext(ctx),
		ctx:      ctx,
		svcCtx:   svcCtx,
		menuRepo: menurepo.NewMenuRepo(svcCtx.DB),
	}
}

func (l *UpdateMenuLogic) UpdateMenu(req *types.UpdateMenuReq) (resp bool, err error) {
	// 1. 检查菜单是否存在
	menu, err := l.menuRepo.GetByMenuID(l.ctx, req.MenuID)
	if err != nil {
		l.Error("UpdateMenu menuRepo.GetByMenuID err:", err.Error())
		if generated.IsNotFound(err) {
			return false, xerr.NewErrMsg("菜单不存在")
		}
		return false, xerr.NewErrCodeMsg(xerr.DbError, "查询菜单失败")
	}

	// 2. 更新菜单信息
	if req.Name != "" {
		menu.Name = req.Name
	}
	if req.Path != "" {
		menu.Path = req.Path
	}
	if req.Component != "" {
		menu.Component = req.Component
	}
	if req.Redirect != "" {
		menu.Redirect = req.Redirect
	}
	if req.Icon != "" {
		menu.Icon = req.Icon
	}
	if req.Sort != 0 {
		menu.Sort = req.Sort
	}
	if req.Status != 0 {
		menu.Status = req.Status
	}

	if req.ParentID != "" {
		menu.ParentID = req.ParentID
	}

	_, err = l.menuRepo.Update(l.ctx, menu)
	if err != nil {
		l.Error("UpdateMenu Update err:", err.Error())
		return false, xerr.NewErrCodeMsg(xerr.DbError, "更新菜单失败")
	}

	return true, nil
}
