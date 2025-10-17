// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package task

import (
	"context"

	"admin_backend/internal/svc"
	"admin_backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type StopTaskLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 停止任务
func NewStopTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StopTaskLogic {
	return &StopTaskLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *StopTaskLogic) StopTask(req *types.StopTaskReq) (resp bool, err error) {
	// todo: add your logic here and delete this line

	return
}
