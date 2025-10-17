// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package task

import (
	"context"

	"admin_backend/internal/svc"
	"admin_backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreatePlanLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建计划
func NewCreatePlanLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreatePlanLogic {
	return &CreatePlanLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreatePlanLogic) CreatePlan(req *types.CreatePlanReq) (resp *types.CreatePlanResp, err error) {
	// todo: add your logic here and delete this line

	return
}
