// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package task

import (
	"context"

	"admin_backend/internal/svc"
	"admin_backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListTaskLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取任务列表
func NewListTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListTaskLogic {
	return &ListTaskLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListTaskLogic) ListTask(req *types.TaskListReq) (resp *types.TaskListResp, err error) {
	// todo: add your logic here and delete this line

	return
}
