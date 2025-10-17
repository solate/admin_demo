package tenant

import (
	"context"

	"admin_backend/app/admin/internal/repository/tenantrepo"
	"admin_backend/app/admin/internal/svc"
	"admin_backend/app/admin/internal/types"
	"admin_backend/pkg/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateTenantLogic struct {
	logx.Logger
	ctx        context.Context
	svcCtx     *svc.ServiceContext
	tenantRepo *tenantrepo.TenantRepo
}

// 更新租户
func NewUpdateTenantLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateTenantLogic {
	return &UpdateTenantLogic{
		Logger:     logx.WithContext(ctx),
		ctx:        ctx,
		svcCtx:     svcCtx,
		tenantRepo: tenantrepo.NewTenantRepo(svcCtx.DB),
	}
}

func (l *UpdateTenantLogic) UpdateTenant(req *types.UpdateTenantReq) (resp bool, err error) {
	// 1. 参数验证

	// 2. 检查租户是否存在
	tenant, err := l.tenantRepo.GetByTenantID(l.ctx, req.TenantID)
	if err != nil {
		l.Error("UpdateTenantLogic GetByTenantID err:", err.Error())
		return false, xerr.NewErrCodeMsg(xerr.DbError, "租户不存在")
	}

	// 3. 更新租户信息
	tenant.Name = req.Name
	tenant.Description = req.Description
	tenant.Status = req.Status

	// 4. 保存更新
	_, err = l.tenantRepo.Update(l.ctx, tenant)
	if err != nil {
		l.Error("UpdateTenantLogic Update err:", err.Error())
		return false, xerr.NewErrCodeMsg(xerr.DbError, "更新租户失败")
	}

	return true, nil
}
