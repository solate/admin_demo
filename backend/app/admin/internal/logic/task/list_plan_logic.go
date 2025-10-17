package task

import (
	"context"

	"admin_backend/app/admin/internal/repository/planrepo"
	"admin_backend/app/admin/internal/svc"
	"admin_backend/app/admin/internal/types"
	"admin_backend/pkg/common/xerr"
	"admin_backend/pkg/ent/generated/plan"
	"admin_backend/pkg/ent/generated/predicate"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListPlanLogic struct {
	logx.Logger
	ctx      context.Context
	svcCtx   *svc.ServiceContext
	planRepo *planrepo.PlanRepo
}

// 获取计划列表
func NewListPlanLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListPlanLogic {
	return &ListPlanLogic{
		Logger:   logx.WithContext(ctx),
		ctx:      ctx,
		svcCtx:   svcCtx,
		planRepo: planrepo.NewPlanRepo(svcCtx.DB),
	}
}

func (l *ListPlanLogic) ListPlan(req *types.PlanListReq) (resp *types.PlanListResp, err error) {
	// 1. 构建查询条件
	where := []predicate.Plan{}

	if req.Name != "" {
		where = append(where, plan.NameContains(req.Name))
	}
	if req.Group != "" {
		where = append(where, plan.Group(req.Group))
	}
	if req.Status != 0 {
		where = append(where, plan.Status(req.Status))
	}
	if req.PlanType != "" {
		where = append(where, plan.PlanType(req.PlanType))
	}

	// 2. 查询数据
	plans, total, err := l.planRepo.PageList(l.ctx, req.Current, req.PageSize, where)
	if err != nil {
		l.Error("ListPlan Logic PageList err:", err.Error())
		return nil, xerr.NewErrCodeMsg(xerr.DbError, "list plan page err.")
	}

	// 3. 构建返回结果
	list := make([]*types.PlanInfo, 0)
	for _, p := range plans {
		list = append(list, &types.PlanInfo{
			PlanID:        p.PlanID,
			Name:          p.Name,
			Description:   p.Description,
			Group:         p.Group,
			CronSpec:      p.CronSpec,
			Status:        p.Status,
			PlanType:      p.PlanType,
			Priority:      p.Priority,
			Timeout:       p.Timeout,
			RetryTimes:    p.RetryTimes,
			RetryInterval: p.RetryInterval,
			StartTime:     p.StartTime,
			EndTime:       p.EndTime,
			Command:       p.Command,
			Params:        p.Params,
			CreatedAt:     p.CreatedAt,
		})
	}

	return &types.PlanListResp{
		List: list,
		Page: &types.PageResponse{
			Total:           total,
			PageSize:        len(list),
			Current:         req.Current,
			RequestPageSize: req.PageSize,
		},
	}, nil
}
