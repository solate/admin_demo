// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package user

import (
	"context"

	"admin_backend/internal/svc"
	"admin_backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListLoginLogLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询登录记录
func NewListLoginLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLoginLogLogic {
	return &ListLoginLogLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListLoginLogLogic) ListLoginLog(req *types.LoginLogListReq) (resp *types.LoginLogListResp, err error) {
	// todo: add your logic here and delete this line

	return
}
