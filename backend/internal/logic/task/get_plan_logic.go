// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package task

import (
	"context"

	"admin_backend/internal/svc"
	"admin_backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPlanLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取计划详情
func NewGetPlanLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPlanLogic {
	return &GetPlanLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPlanLogic) GetPlan(req *types.GetPlanReq) (resp *types.PlanInfo, err error) {
	// todo: add your logic here and delete this line

	return
}
