// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package tenant

import (
	"context"

	"admin_backend/internal/svc"
	"admin_backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteTenantLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除租户
func NewDeleteTenantLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteTenantLogic {
	return &DeleteTenantLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteTenantLogic) DeleteTenant(req *types.DeleteTenantReq) (resp bool, err error) {
	// todo: add your logic here and delete this line

	return
}
