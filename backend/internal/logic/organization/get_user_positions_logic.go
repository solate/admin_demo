// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package organization

import (
	"context"

	"admin_backend/internal/svc"
	"admin_backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserPositionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取用户的岗位列表
func NewGetUserPositionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserPositionsLogic {
	return &GetUserPositionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserPositionsLogic) GetUserPositions(req *types.GetUserPositionsReq) (resp *types.GetUserPositionsResp, err error) {
	// todo: add your logic here and delete this line

	return
}
