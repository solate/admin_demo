// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package menu

import (
	"context"

	"admin_backend/internal/svc"
	"admin_backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListMenuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取菜单列表
func NewListMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListMenuLogic {
	return &ListMenuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListMenuLogic) ListMenu(req *types.MenuListReq) (resp *types.MenuListResp, err error) {
	// todo: add your logic here and delete this line

	return
}
