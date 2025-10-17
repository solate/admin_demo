package menu

import (
	"context"

	"admin_backend/app/admin/internal/repository/menurepo"
	"admin_backend/app/admin/internal/svc"
	"admin_backend/app/admin/internal/types"
	"admin_backend/pkg/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMenuLogic struct {
	logx.Logger
	ctx      context.Context
	svcCtx   *svc.ServiceContext
	menuRepo *menurepo.MenuRepo
}

// 获取菜单详情
func NewGetMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMenuLogic {
	return &GetMenuLogic{
		Logger:   logx.WithContext(ctx),
		ctx:      ctx,
		svcCtx:   svcCtx,
		menuRepo: menurepo.NewMenuRepo(svcCtx.DB),
	}
}

func (l *GetMenuLogic) GetMenu(req *types.GetMenuReq) (resp *types.MenuInfo, err error) {
	// 查询菜单信息
	menu, err := l.menuRepo.GetByMenuID(l.ctx, req.MenuID)
	if err != nil {
		l.Error("GetMenu l.menuRepo.GetByMenuID err: ", err.Error())
		return nil, xerr.NewErrCodeMsg(xerr.DbError, "获取菜单信息失败")
	}

	resp = &types.MenuInfo{
		MenuID:    menu.MenuID,
		Code:      menu.Code,
		ParentID:  menu.ParentID,
		Name:      menu.Name,
		Path:      menu.Path,
		Component: menu.Component,
		Redirect:  menu.Redirect,
		Icon:      menu.Icon,
		Sort:      menu.Sort,
		Type:      menu.Type,
		Status:    menu.Status,
		CreatedAt: menu.CreatedAt,
	}

	return
}
