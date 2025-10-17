package tenant

import (
	"context"

	"admin_backend/app/admin/internal/repository/tenantrepo"
	"admin_backend/app/admin/internal/svc"
	"admin_backend/app/admin/internal/types"
	"admin_backend/pkg/common/xerr"
	"admin_backend/pkg/ent/generated"
	"admin_backend/pkg/utils/idgen"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateTenantLogic struct {
	logx.Logger
	ctx        context.Context
	svcCtx     *svc.ServiceContext
	tenantRepo *tenantrepo.TenantRepo
}

// 创建租户
func NewCreateTenantLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateTenantLogic {
	return &CreateTenantLogic{
		Logger:     logx.WithContext(ctx),
		ctx:        ctx,
		svcCtx:     svcCtx,
		tenantRepo: tenantrepo.NewTenantRepo(svcCtx.DB),
	}
}

func (l *CreateTenantLogic) CreateTenant(req *types.CreateTenantReq) (*types.CreateTenantResp, error) {

	tenantID, err := idgen.GenerateUUID()
	if err != nil {
		return nil, xerr.NewErrCodeMsg(xerr.ServerError, "生成租户ID失败")
	}

	if req.Status == 0 {
		req.Status = 1 // 默认启用
	}

	// 创建租户
	newTenant := &generated.Tenant{
		TenantID:    tenantID,
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
		Status:      req.Status,
	}
	tenant, err := l.tenantRepo.Create(l.ctx, newTenant)
	if err != nil {
		l.Error("CreateTenant Create err:", err.Error())
		return nil, xerr.NewErrCodeMsg(xerr.ServerError, "创建租户失败")
	}

	return &types.CreateTenantResp{
		TenantID: tenant.TenantID,
	}, nil
}
