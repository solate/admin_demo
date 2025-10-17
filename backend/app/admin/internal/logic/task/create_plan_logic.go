package task

import (
	"context"

	"admin_backend/app/admin/internal/repository/planrepo"
	"admin_backend/app/admin/internal/svc"
	"admin_backend/app/admin/internal/types"
	"admin_backend/pkg/common/contextutil"
	"admin_backend/pkg/common/xerr"
	"admin_backend/pkg/ent/generated"
	"admin_backend/pkg/utils/idgen"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreatePlanLogic struct {
	logx.Logger
	ctx      context.Context
	svcCtx   *svc.ServiceContext
	planRepo *planrepo.PlanRepo
}

// 创建计划
func NewCreatePlanLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreatePlanLogic {
	return &CreatePlanLogic{
		Logger:   logx.WithContext(ctx),
		ctx:      ctx,
		svcCtx:   svcCtx,
		planRepo: planrepo.NewPlanRepo(svcCtx.DB),
	}
}

func (l *CreatePlanLogic) CreatePlan(req *types.CreatePlanReq) (resp *types.CreatePlanResp, err error) {
	// 1. 检查计划名称是否已存在
	plan, err := l.planRepo.GetByName(l.ctx, req.Name)
	if err != nil && !generated.IsNotFound(err) {
		l.Error("GetPlan l.planRepo.GetByName err: ", err.Error())
		return nil, xerr.NewErrCodeMsg(xerr.DbError, "数据库查询错误")
	}
	if plan != nil {
		return nil, xerr.NewErrCodeMsg(xerr.DbRecordExist, "计划名称已存在")
	}

	// 2. 生成计划ID
	planID, err := idgen.GenerateUUID()
	if err != nil {
		l.Error("CreatePlan GenerateUUID err:", err.Error())
		return nil, xerr.NewErrCodeMsg(xerr.ServerError, "生成计划ID失败")
	}

	// 3. 设置默认值
	if req.Status == 0 {
		req.Status = 1 // 默认启用
	}

	// 4. 创建计划
	newPlan := &generated.Plan{
		TenantCode:    contextutil.GetTenantCodeFromCtx(l.ctx),
		PlanID:        planID,
		Name:          req.Name,
		Description:   req.Description,
		Group:         req.Group,
		CronSpec:      req.CronSpec,
		Status:        req.Status,
		PlanType:      req.PlanType,
		Priority:      req.Priority,
		Timeout:       req.Timeout,
		RetryTimes:    req.RetryTimes,
		RetryInterval: req.RetryInterval,
		StartTime:     req.StartTime,
		EndTime:       req.EndTime,
		Command:       req.Command,
		Params:        req.Params,
	}

	plan, err = l.planRepo.Create(l.ctx, newPlan)
	if err != nil {
		l.Error("CreatePlan Create err:", err.Error())
		return nil, xerr.NewErrCodeMsg(xerr.ServerError, "创建计划失败")
	}

	// 5. 返回结果
	return &types.CreatePlanResp{
		PlanID: newPlan.PlanID,
	}, nil
}
