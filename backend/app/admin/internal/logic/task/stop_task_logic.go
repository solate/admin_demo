package task

import (
	"context"
	"time"

	"admin_backend/app/admin/internal/repository/taskrepo"
	"admin_backend/app/admin/internal/svc"
	"admin_backend/app/admin/internal/types"
	"admin_backend/pkg/common/xerr"
	"admin_backend/pkg/ent/generated"

	"github.com/zeromicro/go-zero/core/logx"
)

type StopTaskLogic struct {
	logx.Logger
	ctx      context.Context
	svcCtx   *svc.ServiceContext
	taskRepo *taskrepo.TaskRepo
}

// NewStopTaskLogic 停止任务
func NewStopTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StopTaskLogic {
	return &StopTaskLogic{
		Logger:   logx.WithContext(ctx),
		ctx:      ctx,
		svcCtx:   svcCtx,
		taskRepo: taskrepo.NewTaskRepo(svcCtx.DB),
	}
}

func (l *StopTaskLogic) StopTask(req *types.StopTaskReq) (resp bool, err error) {
	// 1. 查询任务是否存在
	task, err := l.taskRepo.GetByTaskID(l.ctx, req.TaskID)
	if err != nil {
		if generated.IsNotFound(err) {
			return false, xerr.NewErrCodeMsg(xerr.RecordNotFound, "任务不存在")
		}
		l.Error("StopTask GetByTaskID err:", err.Error())
		return false, xerr.NewErrCodeMsg(xerr.DbError, "查询任务失败")
	}

	// 2. 检查任务状态是否可以停止
	if task.Status == "stopped" || task.Status == "completed" || task.Status == "failed" {
		return false, xerr.NewErrCodeMsg(xerr.ParamError, "任务已经结束，无法停止")
	}

	// 3. 更新任务状态
	now := time.Now()
	task.Status = "stopped"
	task.EndTime = now.UnixMilli()
	if task.StartTime != 0 {
		task.Duration = int(task.StartTime - task.EndTime)
	}
	task.Result = "任务已手动停止"

	_, err = l.taskRepo.Update(l.ctx, task)
	if err != nil {
		l.Error("StopTask Update err:", err.Error())
		return false, xerr.NewErrCodeMsg(xerr.DbError, "停止任务失败")
	}

	return true, nil
}
