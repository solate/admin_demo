package task

import (
	"context"

	"admin_backend/app/admin/internal/repository/taskrepo"
	"admin_backend/app/admin/internal/svc"
	"admin_backend/app/admin/internal/types"
	"admin_backend/pkg/common/xerr"
	"admin_backend/pkg/ent/generated"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteTaskLogic struct {
	logx.Logger
	ctx      context.Context
	svcCtx   *svc.ServiceContext
	taskRepo *taskrepo.TaskRepo
}

// 删除任务
func NewDeleteTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteTaskLogic {
	return &DeleteTaskLogic{
		Logger:   logx.WithContext(ctx),
		ctx:      ctx,
		svcCtx:   svcCtx,
		taskRepo: taskrepo.NewTaskRepo(svcCtx.DB),
	}
}

func (l *DeleteTaskLogic) DeleteTask(req *types.DeleteTaskReq) (resp bool, err error) {
	// 1. 检查任务是否存在
	task, err := l.taskRepo.GetByTaskID(l.ctx, req.TaskID)
	if err != nil {
		l.Error("DeleteTask taskRepo.GetByTaskID err:", err.Error())
		if generated.IsNotFound(err) {
			return false, xerr.NewErrMsg("任务不存在")
		}
		return false, xerr.NewErrCodeMsg(xerr.DbError, "查询任务失败")
	}

	// 2. 软删除任务
	_, err = l.taskRepo.Delete(l.ctx, task.TaskID)
	if err != nil {
		l.Error("DeleteTask taskRepo.Delete err:", err.Error())
		return false, xerr.NewErrCodeMsg(xerr.DbError, "删除任务失败")
	}

	return true, nil
}
