// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package tenant

import (
	"context"

	"admin_backend/internal/svc"
	"admin_backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetTenantLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取租户详情
func NewGetTenantLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTenantLogic {
	return &GetTenantLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTenantLogic) GetTenant(req *types.GetTenantReq) (resp *types.GetTenantResp, err error) {
	// todo: add your logic here and delete this line

	return
}
