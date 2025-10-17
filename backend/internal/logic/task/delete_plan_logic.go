// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package task

import (
	"context"

	"admin_backend/internal/svc"
	"admin_backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeletePlanLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除计划
func NewDeletePlanLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeletePlanLogic {
	return &DeletePlanLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeletePlanLogic) DeletePlan(req *types.DeletePlanReq) (resp bool, err error) {
	// todo: add your logic here and delete this line

	return
}
