// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package task

import (
	"context"

	"admin_backend/internal/svc"
	"admin_backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdatePlanLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新计划
func NewUpdatePlanLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatePlanLogic {
	return &UpdatePlanLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdatePlanLogic) UpdatePlan(req *types.UpdatePlanReq) (resp bool, err error) {
	// todo: add your logic here and delete this line

	return
}
