package planrepo

import (
	"context"
	"time"

	"admin_backend/pkg/common"
	"admin_backend/pkg/common/contextutil"
	"admin_backend/pkg/ent"
	"admin_backend/pkg/ent/generated"
	"admin_backend/pkg/ent/generated/plan"
	"admin_backend/pkg/ent/generated/predicate"
)

type PlanRepo struct {
	db *ent.Client
}

// NewPlanRepo 创建计划仓储实例
func NewPlanRepo(db *ent.Client) *PlanRepo {
	return &PlanRepo{db: db}
}

func (r *PlanRepo) Create(ctx context.Context, plan *generated.Plan) (*generated.Plan, error) {
	now := time.Now().UnixMilli()
	return r.db.Plan.Create().
		SetCreatedAt(now).
		SetUpdatedAt(now).
		SetTenantCode(plan.TenantCode).
		SetPlanID(plan.PlanID).
		SetName(plan.Name).
		SetDescription(plan.Description).
		SetGroup(plan.Group).
		SetCronSpec(plan.CronSpec).
		SetStatus(plan.Status).
		SetPlanType(plan.PlanType).
		SetPriority(plan.Priority).
		SetTimeout(plan.Timeout).
		SetRetryTimes(plan.RetryTimes).
		SetRetryInterval(plan.RetryInterval).
		SetStartTime(plan.StartTime).
		SetEndTime(plan.EndTime).
		SetCommand(plan.Command).
		SetParams(plan.Params).
		Save(ctx)
}

func (r *PlanRepo) Update(ctx context.Context, update *generated.Plan) (int, error) {
	now := time.Now().UnixMilli()
	update.UpdatedAt = now
	return r.db.Plan.Update().
		SetUpdatedAt(now).
		SetName(update.Name).
		SetDescription(update.Description).
		SetGroup(update.Group).
		SetCronSpec(update.CronSpec).
		SetStatus(update.Status).
		SetPlanType(update.PlanType).
		SetPriority(update.Priority).
		SetTimeout(update.Timeout).
		SetRetryTimes(update.RetryTimes).
		SetRetryInterval(update.RetryInterval).
		SetStartTime(update.StartTime).
		SetEndTime(update.EndTime).
		SetCommand(update.Command).
		SetParams(update.Params).
		Where(plan.PlanID(update.PlanID)).Save(ctx)
}

func (r *PlanRepo) GetByPlanID(ctx context.Context, planID string) (*generated.Plan, error) {
	return r.Get(ctx, []predicate.Plan{plan.PlanID(planID)})
}

func (r *PlanRepo) GetByName(ctx context.Context, name string) (*generated.Plan, error) {
	return r.Get(ctx, []predicate.Plan{plan.Name(name)})
}

// GetByPlanType 根据计划类型查询
func (r *PlanRepo) GetByPlanType(ctx context.Context, planType string) ([]*generated.Plan, error) {
	where := []predicate.Plan{plan.PlanType(planType)}
	where = r.defaultQuery(ctx, where)
	return r.db.Plan.Query().Where(where...).All(ctx)
}

// GetByStatus 根据状态查询计划
func (r *PlanRepo) GetByStatus(ctx context.Context, status int) ([]*generated.Plan, error) {
	where := []predicate.Plan{plan.Status(status)}
	where = r.defaultQuery(ctx, where)
	return r.db.Plan.Query().Where(where...).All(ctx)
}

// defaultQuery 默认查询条件
func (r *PlanRepo) defaultQuery(ctx context.Context, where []predicate.Plan) []predicate.Plan {
	where = append(where, plan.DeletedAtIsNil())
	where = append(where, plan.TenantCode(contextutil.GetTenantCodeFromCtx(ctx)))
	return where
}

func (r *PlanRepo) Get(ctx context.Context, where []predicate.Plan) (*generated.Plan, error) {
	where = r.defaultQuery(ctx, where)
	return r.db.Plan.Query().Where(where...).Only(ctx)
}

func (r *PlanRepo) PageList(ctx context.Context, current, limit int, where []predicate.Plan) ([]*generated.Plan, int, error) {
	where = r.defaultQuery(ctx, where)

	offset := common.Offset(current, limit)
	query := r.db.Plan.Query().Where(where...).Order(generated.Desc(plan.FieldCreatedAt))

	// 查询总数
	total, err := query.Count(ctx)
	if err != nil || total == 0 {
		return nil, 0, err
	}

	// 分页查询
	list, err := query.Offset(offset).Limit(limit).All(ctx)
	return list, total, err
}

// Delete 根据计划ID删除计划，软删除
func (r *PlanRepo) Delete(ctx context.Context, planID string) (int, error) {
	now := time.Now().UnixMilli()
	return r.db.Plan.Update().
		SetDeletedAt(now).
		Where(plan.PlanID(planID)).Save(ctx)
}

// UpdateStatus 更新计划状态
func (r *PlanRepo) UpdateStatus(ctx context.Context, planID string, status int) (int, error) {
	return r.db.Plan.Update().
		SetStatus(status).
		SetUpdatedAt(time.Now().UnixMilli()).
		Where(plan.PlanID(planID)).Save(ctx)
}
