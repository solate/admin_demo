package task

import (
	"context"

	"admin_backend/app/admin/internal/repository/planrepo"
	"admin_backend/app/admin/internal/svc"
	"admin_backend/app/admin/internal/types"
	"admin_backend/pkg/common/xerr"
	"admin_backend/pkg/ent/generated"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdatePlanLogic struct {
	logx.Logger
	ctx      context.Context
	svcCtx   *svc.ServiceContext
	planRepo *planrepo.PlanRepo
}

// 更新计划
func NewUpdatePlanLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatePlanLogic {
	return &UpdatePlanLogic{
		Logger:   logx.WithContext(ctx),
		ctx:      ctx,
		svcCtx:   svcCtx,
		planRepo: planrepo.NewPlanRepo(svcCtx.DB),
	}
}

func (l *UpdatePlanLogic) UpdatePlan(req *types.UpdatePlanReq) (resp bool, err error) {
	// 1. 检查计划是否存在
	plan, err := l.planRepo.GetByPlanID(l.ctx, req.PlanID)
	if err != nil {
		l.Error("UpdatePlan planRepo.GetByPlanID err:", err.Error())
		if generated.IsNotFound(err) {
			return false, xerr.NewErrMsg("计划不存在")
		}
		return false, xerr.NewErrCodeMsg(xerr.DbError, "查询计划失败")
	}

	// 2. 更新计划信息
	plan.Name = req.Name
	plan.Description = req.Description
	plan.Group = req.Group
	plan.CronSpec = req.CronSpec
	plan.Status = req.Status
	plan.Priority = req.Priority
	plan.Timeout = req.Timeout
	plan.RetryTimes = req.RetryTimes
	plan.RetryInterval = req.RetryInterval
	plan.StartTime = req.StartTime
	plan.EndTime = req.EndTime
	plan.Command = req.Command
	plan.Params = req.Params

	_, err = l.planRepo.Update(l.ctx, plan)
	if err != nil {
		l.Error("UpdatePlan Update err:", err.Error())
		return false, xerr.NewErrCodeMsg(xerr.DbError, "更新计划失败")
	}

	return true, nil
}
