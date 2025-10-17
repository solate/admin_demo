package tenant

import (
	"context"

	"admin_backend/app/admin/internal/repository/tenantrepo"
	"admin_backend/app/admin/internal/svc"
	"admin_backend/app/admin/internal/types"
	"admin_backend/pkg/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetTenantLogic struct {
	logx.Logger
	ctx        context.Context
	svcCtx     *svc.ServiceContext
	tenantRepo *tenantrepo.TenantRepo
}

// 获取租户详情
func NewGetTenantLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTenantLogic {
	return &GetTenantLogic{
		Logger:     logx.WithContext(ctx),
		ctx:        ctx,
		svcCtx:     svcCtx,
		tenantRepo: tenantrepo.NewTenantRepo(svcCtx.DB),
	}
}

func (l *GetTenantLogic) GetTenant(req *types.GetTenantReq) (*types.GetTenantResp, error) {

	// 获取租户信息
	tenant, err := l.tenantRepo.GetByTenantID(l.ctx, req.TenantID)
	if err != nil {
		l.Error("GetTenant l.tenantRepo.GetByTenantID err: ", err)
		return nil, xerr.NewErrCodeMsg(xerr.DbError, "获取租户信息失败")
	}

	return &types.GetTenantResp{
		TenantInfo: types.TenantInfo{
			TenantID:    tenant.TenantID,
			Code:        tenant.Code,
			Name:        tenant.Name,
			Description: tenant.Description,
			Status:      tenant.Status,
		},
	}, nil
}
