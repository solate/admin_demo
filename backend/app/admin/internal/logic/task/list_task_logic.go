package task

import (
	"context"

	"admin_backend/app/admin/internal/repository/taskrepo"
	"admin_backend/app/admin/internal/svc"
	"admin_backend/app/admin/internal/types"
	"admin_backend/pkg/common/xerr"
	"admin_backend/pkg/ent/generated/predicate"
	"admin_backend/pkg/ent/generated/task"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListTaskLogic struct {
	logx.Logger
	ctx      context.Context
	svcCtx   *svc.ServiceContext
	taskRepo *taskrepo.TaskRepo
}

// 获取任务列表
func NewListTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListTaskLogic {
	return &ListTaskLogic{
		Logger:   logx.WithContext(ctx),
		ctx:      ctx,
		svcCtx:   svcCtx,
		taskRepo: taskrepo.NewTaskRepo(svcCtx.DB),
	}
}

func (l *ListTaskLogic) ListTask(req *types.TaskListReq) (resp *types.TaskListResp, err error) {
	// 1. 构建查询条件
	where := []predicate.Task{}

	if req.PlanID != "" {
		where = append(where, task.PlanID(req.PlanID))
	}
	if req.Name != "" {
		where = append(where, task.NameContains(req.Name))
	}
	if req.Group != "" {
		where = append(where, task.Group(req.Group))
	}
	if req.Status != "" {
		where = append(where, task.Status(req.Status))
	}
	if req.StartTime != 0 {
		where = append(where, task.StartTimeGTE(req.StartTime))
	}
	if req.EndTime != 0 {
		where = append(where, task.EndTimeLTE(req.EndTime))
	}

	// 2. 查询数据
	tasks, total, err := l.taskRepo.PageList(l.ctx, req.Current, req.PageSize, where)
	if err != nil {
		l.Error("ListTask Logic PageList err:", err.Error())
		return nil, xerr.NewErrCodeMsg(xerr.DbError, "list task page err.")
	}

	// 3. 构建返回结果
	list := make([]*types.TaskInfo, 0)
	for _, t := range tasks {
		list = append(list, &types.TaskInfo{
			TaskID:        t.TaskID,
			Name:          t.Name,
			PlanID:        t.PlanID,
			PlanType:      t.PlanType,
			Group:         t.Group,
			Priority:      t.Priority,
			Status:        t.Status,
			PlannedTime:   t.PlannedTime,
			StartTime:     t.StartTime,
			EndTime:       t.EndTime,
			Duration:      t.Duration,
			Result:        t.Result,
			Error:         t.Error,
			RetryCount:    t.RetryCount,
			NextRetryTime: t.NextRetryTime,
			CreatedAt:     t.CreatedAt,
		})
	}

	return &types.TaskListResp{
		List: list,
		Page: &types.PageResponse{
			Total:           total,
			PageSize:        len(list),
			Current:         req.Current,
			RequestPageSize: req.PageSize,
		},
	}, nil
}
