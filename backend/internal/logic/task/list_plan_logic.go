// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package task

import (
	"context"

	"admin_backend/internal/svc"
	"admin_backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListPlanLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取计划列表
func NewListPlanLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListPlanLogic {
	return &ListPlanLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListPlanLogic) ListPlan(req *types.PlanListReq) (resp *types.PlanListResp, err error) {
	// todo: add your logic here and delete this line

	return
}
