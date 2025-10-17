// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package organization

import (
	"context"

	"admin_backend/internal/svc"
	"admin_backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RemoveUserPositionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 移除用户岗位
func NewRemoveUserPositionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveUserPositionLogic {
	return &RemoveUserPositionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RemoveUserPositionLogic) RemoveUserPosition(req *types.RemoveUserPositionReq) (resp bool, err error) {
	// todo: add your logic here and delete this line

	return
}
