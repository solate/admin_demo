// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package menu

import (
	"context"

	"admin_backend/internal/svc"
	"admin_backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAllMenuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取所有菜单
func NewGetAllMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAllMenuLogic {
	return &GetAllMenuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAllMenuLogic) GetAllMenu() (resp *types.MenuTreeResp, err error) {
	// todo: add your logic here and delete this line

	return
}
