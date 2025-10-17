// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package tenant

import (
	"context"

	"admin_backend/internal/svc"
	"admin_backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListTenantLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取租户列表
func NewListTenantLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListTenantLogic {
	return &ListTenantLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListTenantLogic) ListTenant(req *types.ListTenantReq) (resp *types.ListTenantResp, err error) {
	// todo: add your logic here and delete this line

	return
}
