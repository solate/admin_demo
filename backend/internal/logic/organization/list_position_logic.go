// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package organization

import (
	"context"

	"admin_backend/internal/svc"
	"admin_backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListPositionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取岗位列表
func NewListPositionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListPositionLogic {
	return &ListPositionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListPositionLogic) ListPosition(req *types.PositionListReq) (resp *types.PositionListResp, err error) {
	// todo: add your logic here and delete this line

	return
}
