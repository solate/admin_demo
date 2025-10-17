package tenant

import (
	"context"

	"admin_backend/app/admin/internal/repository/tenantrepo"
	"admin_backend/app/admin/internal/svc"
	"admin_backend/app/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteTenantLogic struct {
	logx.Logger
	ctx        context.Context
	svcCtx     *svc.ServiceContext
	tenantRepo *tenantrepo.TenantRepo
}

// 删除租户
func NewDeleteTenantLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteTenantLogic {
	return &DeleteTenantLogic{
		Logger:     logx.WithContext(ctx),
		ctx:        ctx,
		svcCtx:     svcCtx,
		tenantRepo: tenantrepo.NewTenantRepo(svcCtx.DB),
	}
}

func (l *DeleteTenantLogic) DeleteTenant(req *types.DeleteTenantReq) (resp bool, err error) {

	// 2. 检查租户是否存在
	tenant, err := l.tenantRepo.GetByTenantID(l.ctx, req.TenantID)
	if err != nil {
		l.Error("DeleteTenantLogic GetByTenantID err:", err.Error())
		return false, err
	}

	// 3. 软删除租户
	_, err = l.tenantRepo.DeleteByTenantID(l.ctx, tenant)
	if err != nil {
		l.Error("DeleteTenantLogic DeleteByTenantID err:", err.Error())
		return false, err
	}

	return true, nil
}
