package task

import (
	"context"
	"time"

	"admin_backend/app/admin/internal/repository/planrepo"
	"admin_backend/app/admin/internal/repository/taskrepo"
	"admin_backend/app/admin/internal/svc"
	"admin_backend/app/admin/internal/types"
	"admin_backend/pkg/common/contextutil"
	"admin_backend/pkg/common/xerr"
	"admin_backend/pkg/ent/generated"
	"admin_backend/pkg/utils/idgen"

	"github.com/zeromicro/go-zero/core/logx"
)

type TriggerTaskLogic struct {
	logx.Logger
	ctx      context.Context
	svcCtx   *svc.ServiceContext
	planRepo *planrepo.PlanRepo
	taskRepo *taskrepo.TaskRepo
}

// 手动触发任务
func NewTriggerTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TriggerTaskLogic {
	return &TriggerTaskLogic{
		Logger:   logx.WithContext(ctx),
		ctx:      ctx,
		svcCtx:   svcCtx,
		planRepo: planrepo.NewPlanRepo(svcCtx.DB),
		taskRepo: taskrepo.NewTaskRepo(svcCtx.DB),
	}
}

func (l *TriggerTaskLogic) TriggerTask(req *types.TriggerTaskReq) (resp *types.TriggerTaskResp, err error) {
	// 1. 查询计划是否存在
	plan, err := l.planRepo.GetByPlanID(l.ctx, req.PlanID)
	if err != nil {
		l.Error("TriggerTask GetByPlanID err: ", err.Error())
		return nil, xerr.NewErrCodeMsg(xerr.DbError, "获取计划详情失败")
	}

	// 2. 生成任务ID
	taskID, err := idgen.GenerateUUID()
	if err != nil {
		l.Error("TriggerTask GenerateUUID err:", err.Error())
		return nil, xerr.NewErrCodeMsg(xerr.ServerError, "生成任务ID失败")
	}

	// 3. 创建任务
	now := time.Now().UnixMilli()
	newTask := &generated.Task{
		TenantCode:  contextutil.GetTenantCodeFromCtx(l.ctx),
		TaskID:      taskID,
		Name:        plan.Name,
		PlanID:      plan.PlanID,
		PlanType:    plan.PlanType,
		Group:       plan.Group,
		Priority:    plan.Priority,
		Status:      "pending",
		PlannedTime: now,
	}

	task, err := l.taskRepo.Create(l.ctx, newTask)
	if err != nil {
		l.Error("TriggerTask Create err:", err.Error())
		return nil, xerr.NewErrCodeMsg(xerr.ServerError, "创建任务失败")
	}

	// 4. 返回结果
	return &types.TriggerTaskResp{
		TaskID: task.TaskID,
	}, nil
}
