package task

import (
	"context"

	"admin_backend/app/admin/internal/repository/planrepo"
	"admin_backend/app/admin/internal/svc"
	"admin_backend/app/admin/internal/types"
	"admin_backend/pkg/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPlanLogic struct {
	logx.Logger
	ctx      context.Context
	svcCtx   *svc.ServiceContext
	planRepo *planrepo.PlanRepo
}

// 获取计划详情
func NewGetPlanLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPlanLogic {
	return &GetPlanLogic{
		Logger:   logx.WithContext(ctx),
		ctx:      ctx,
		svcCtx:   svcCtx,
		planRepo: planrepo.NewPlanRepo(svcCtx.DB),
	}
}

func (l *GetPlanLogic) GetPlan(req *types.GetPlanReq) (resp *types.PlanInfo, err error) {
	// 查询计划详情
	plan, err := l.planRepo.GetByPlanID(l.ctx, req.PlanID)
	if err != nil {
		l.Error("GetPlan GetByPlanID err: ", err.Error())
		return nil, xerr.NewErrCodeMsg(xerr.DbError, "获取计划详情失败")
	}

	resp = &types.PlanInfo{
		PlanID:        plan.PlanID,
		Name:          plan.Name,
		Description:   plan.Description,
		Group:         plan.Group,
		CronSpec:      plan.CronSpec,
		Status:        plan.Status,
		PlanType:      plan.PlanType,
		Priority:      plan.Priority,
		Timeout:       plan.Timeout,
		RetryTimes:    plan.RetryTimes,
		RetryInterval: plan.RetryInterval,
		StartTime:     plan.StartTime,
		EndTime:       plan.EndTime,
		Command:       plan.Command,
		Params:        plan.Params,
		CreatedAt:     plan.CreatedAt,
	}

	return
}
