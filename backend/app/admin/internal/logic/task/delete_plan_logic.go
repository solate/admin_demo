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

type DeletePlanLogic struct {
	logx.Logger
	ctx      context.Context
	svcCtx   *svc.ServiceContext
	planRepo *planrepo.PlanRepo
}

// 删除计划
func NewDeletePlanLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeletePlanLogic {
	return &DeletePlanLogic{
		Logger:   logx.WithContext(ctx),
		ctx:      ctx,
		svcCtx:   svcCtx,
		planRepo: planrepo.NewPlanRepo(svcCtx.DB),
	}
}

func (l *DeletePlanLogic) DeletePlan(req *types.DeletePlanReq) (resp bool, err error) {
	// 1. 查询计划是否存在
	plan, err := l.planRepo.GetByPlanID(l.ctx, req.PlanID)
	if err != nil {
		if generated.IsNotFound(err) {
			return false, xerr.NewErrCodeMsg(xerr.RecordNotFound, "计划不存在")
		}
		l.Error("DeletePlan GetByPlanID err:", err.Error())
		return false, xerr.NewErrCodeMsg(xerr.DbError, "查询计划失败")
	}

	// 2. 删除计划
	_, err = l.planRepo.Delete(l.ctx, plan.PlanID)
	if err != nil {
		l.Error("DeletePlan Delete err:", err.Error())
		return false, xerr.NewErrCodeMsg(xerr.DbError, "删除计划失败")
	}

	return true, nil
}
