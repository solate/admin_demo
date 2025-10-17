// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package organization

import (
	"context"

	"admin_backend/internal/svc"
	"admin_backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AssignUserPositionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 分配用户岗位
func NewAssignUserPositionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AssignUserPositionLogic {
	return &AssignUserPositionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AssignUserPositionLogic) AssignUserPosition(req *types.AssignUserPositionReq) (resp bool, err error) {
	// todo: add your logic here and delete this line

	return
}
