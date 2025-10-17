package menu

import (
	"context"

	"admin_backend/app/admin/internal/repository/menurepo"
	"admin_backend/app/admin/internal/svc"
	"admin_backend/app/admin/internal/types"
	"admin_backend/pkg/common/xerr"
	"admin_backend/pkg/ent/generated/menu"
	"admin_backend/pkg/ent/generated/predicate"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListMenuLogic struct {
	logx.Logger
	ctx      context.Context
	svcCtx   *svc.ServiceContext
	menuRepo *menurepo.MenuRepo
}

// 获取菜单列表
func NewListMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListMenuLogic {
	return &ListMenuLogic{
		Logger:   logx.WithContext(ctx),
		ctx:      ctx,
		svcCtx:   svcCtx,
		menuRepo: menurepo.NewMenuRepo(svcCtx.DB),
	}
}

func (l *ListMenuLogic) ListMenu(req *types.MenuListReq) (resp *types.MenuListResp, err error) {
	// 1. 构建查询条件
	where := []predicate.Menu{}

	if req.Name != "" {
		where = append(where, menu.NameContains(req.Name))
	}

	if req.Status != 0 {
		where = append(where, menu.StatusEQ(req.Status))
	}

	if req.Type != "" {
		where = append(where, menu.TypeEQ(req.Type))
	}

	// 2. 分页查询
	list, total, err := l.menuRepo.PageList(l.ctx, req.Current, req.PageSize, where)
	if err != nil {
		l.Error("ListMenu Logic PageList err:", err.Error())
		return nil, xerr.NewErrCodeMsg(xerr.DbError, "list menu page err.")
	}

	// 3. 构建返回结果
	menuList := make([]*types.MenuInfo, 0)
	for _, menu := range list {
		menuList = append(menuList, &types.MenuInfo{
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
		})
	}

	return &types.MenuListResp{
		List: menuList,
		Page: &types.PageResponse{
			Total:           total,
			PageSize:        len(list),
			Current:         req.Current,
			RequestPageSize: req.PageSize,
		},
	}, nil
}
