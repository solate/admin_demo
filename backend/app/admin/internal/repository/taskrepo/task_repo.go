package taskrepo

import (
	"context"
	"time"

	"admin_backend/pkg/common"
	"admin_backend/pkg/common/contextutil"
	"admin_backend/pkg/ent"
	"admin_backend/pkg/ent/generated"
	"admin_backend/pkg/ent/generated/predicate"
	"admin_backend/pkg/ent/generated/task"
)

type TaskRepo struct {
	db *ent.Client
}

// NewTaskRepo 创建任务仓储实例
func NewTaskRepo(db *ent.Client) *TaskRepo {
	return &TaskRepo{db: db}
}

func (r *TaskRepo) Create(ctx context.Context, task *generated.Task) (*generated.Task, error) {
	now := time.Now().UnixMilli()
	return r.db.Task.Create().
		SetCreatedAt(now).
		SetUpdatedAt(now).
		SetTenantCode(task.TenantCode).
		SetTaskID(task.TaskID).
		SetName(task.Name).
		SetPlanID(task.PlanID).
		SetPlanType(task.PlanType).
		SetGroup(task.Group).
		SetPriority(task.Priority).
		SetStatus(task.Status).
		SetPlannedTime(task.PlannedTime).
		SetStartTime(task.StartTime).
		SetEndTime(task.EndTime).
		SetDuration(task.Duration).
		SetResult(task.Result).
		SetError(task.Error).
		SetRetryCount(task.RetryCount).
		SetNextRetryTime(task.NextRetryTime).
		Save(ctx)
}

func (r *TaskRepo) Update(ctx context.Context, update *generated.Task) (int, error) {
	now := time.Now().UnixMilli()
	update.UpdatedAt = now
	return r.db.Task.Update().
		SetUpdatedAt(now).
		SetName(update.Name).
		SetPlanType(update.PlanType).
		SetGroup(update.Group).
		SetPriority(update.Priority).
		SetStatus(update.Status).
		SetPlannedTime(update.PlannedTime).
		SetStartTime(update.StartTime).
		SetEndTime(update.EndTime).
		SetDuration(update.Duration).
		SetResult(update.Result).
		SetError(update.Error).
		SetRetryCount(update.RetryCount).
		SetNextRetryTime(update.NextRetryTime).
		Where(task.TaskID(update.TaskID)).Save(ctx)
}

func (r *TaskRepo) GetByTaskID(ctx context.Context, taskID string) (*generated.Task, error) {
	return r.Get(ctx, []predicate.Task{task.TaskID(taskID)})
}

func (r *TaskRepo) GetByPlanID(ctx context.Context, planID string) ([]*generated.Task, error) {
	where := []predicate.Task{task.PlanID(planID)}
	where = r.defaultQuery(ctx, where)
	return r.db.Task.Query().Where(where...).All(ctx)
}

// GetByStatus 根据状态查询任务
func (r *TaskRepo) GetByStatus(ctx context.Context, status string) ([]*generated.Task, error) {
	where := []predicate.Task{task.Status(status)}
	where = r.defaultQuery(ctx, where)
	return r.db.Task.Query().Where(where...).All(ctx)
}

// defaultQuery 默认查询条件
func (r *TaskRepo) defaultQuery(ctx context.Context, where []predicate.Task) []predicate.Task {
	where = append(where, task.DeletedAtIsNil())
	where = append(where, task.TenantCode(contextutil.GetTenantCodeFromCtx(ctx)))
	return where
}

func (r *TaskRepo) Get(ctx context.Context, where []predicate.Task) (*generated.Task, error) {
	where = r.defaultQuery(ctx, where)
	return r.db.Task.Query().Where(where...).Only(ctx)
}

func (r *TaskRepo) PageList(ctx context.Context, current, limit int, where []predicate.Task) ([]*generated.Task, int, error) {
	where = r.defaultQuery(ctx, where)

	offset := common.Offset(current, limit)
	query := r.db.Task.Query().Where(where...).Order(generated.Desc(task.FieldCreatedAt))

	// 查询总数
	total, err := query.Count(ctx)
	if err != nil || total == 0 {
		return nil, 0, err
	}

	// 分页查询
	list, err := query.Offset(offset).Limit(limit).All(ctx)
	return list, total, err
}

// Delete 根据任务ID删除任务，软删除
func (r *TaskRepo) Delete(ctx context.Context, taskID string) (int, error) {
	now := time.Now().UnixMilli()
	return r.db.Task.Update().
		SetDeletedAt(now).
		Where(task.TaskID(taskID)).Save(ctx)
}

// UpdateStatus 更新任务状态
func (r *TaskRepo) UpdateStatus(ctx context.Context, taskID string, status string) (int, error) {
	return r.db.Task.Update().
		SetStatus(status).
		SetUpdatedAt(time.Now().UnixMilli()).
		Where(task.TaskID(taskID)).Save(ctx)
}

// UpdateTaskResult 更新任务执行结果
func (r *TaskRepo) UpdateTaskResult(ctx context.Context, taskID string, result string, err string) (int, error) {
	return r.db.Task.Update().
		SetResult(result).
		SetError(err).
		SetUpdatedAt(time.Now().UnixMilli()).
		Where(task.TaskID(taskID)).Save(ctx)
}
