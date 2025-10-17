// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package tenant

import (
	"context"

	"admin_backend/internal/svc"
	"admin_backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateTenantLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建租户
func NewCreateTenantLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateTenantLogic {
	return &CreateTenantLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateTenantLogic) CreateTenant(req *types.CreateTenantReq) (resp *types.CreateTenantResp, err error) {
	// todo: add your logic here and delete this line

	return
}
