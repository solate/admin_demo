// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package task

import (
	"context"

	"admin_backend/internal/svc"
	"admin_backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type TriggerTaskLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 手动触发任务
func NewTriggerTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TriggerTaskLogic {
	return &TriggerTaskLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TriggerTaskLogic) TriggerTask(req *types.TriggerTaskReq) (resp *types.TriggerTaskResp, err error) {
	// todo: add your logic here and delete this line

	return
}
