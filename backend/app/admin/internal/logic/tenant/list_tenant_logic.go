package tenant

import (
	"context"

	"admin_backend/app/admin/internal/repository/tenantrepo"
	"admin_backend/app/admin/internal/svc"
	"admin_backend/app/admin/internal/types"
	"admin_backend/pkg/common/xerr"
	"admin_backend/pkg/ent/generated/predicate"
	"admin_backend/pkg/ent/generated/tenant"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListTenantLogic struct {
	logx.Logger
	ctx        context.Context
	svcCtx     *svc.ServiceContext
	tenantRepo *tenantrepo.TenantRepo
}

// 获取租户列表
func NewListTenantLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListTenantLogic {
	return &ListTenantLogic{
		Logger:     logx.WithContext(ctx),
		ctx:        ctx,
		svcCtx:     svcCtx,
		tenantRepo: tenantrepo.NewTenantRepo(svcCtx.DB),
	}
}

func (l *ListTenantLogic) ListTenant(req *types.ListTenantReq) (*types.ListTenantResp, error) {

	// 1. 构建查询条件
	where := []predicate.Tenant{}
	// 2. 状态筛选
	if req.Status > 0 {
		where = append(where, tenant.Status(req.Status))
	}

	// 4. 分页查询
	list, total, err := l.tenantRepo.PageList(l.ctx, req.Current, req.PageSize, where)
	if err != nil {
		l.Error("ListTenant PageList err:", err.Error())
		return nil, xerr.NewErrCode(xerr.DbError)
	}

	// 5. 构建返回结果
	tenantList := make([]*types.TenantInfo, 0)
	for _, t := range list {
		tenantList = append(tenantList, &types.TenantInfo{
			TenantID:    t.TenantID,
			Code:        t.Code,
			Name:        t.Name,
			Description: t.Description,
			Status:      t.Status,
		})
	}

	return &types.ListTenantResp{
		List: tenantList,
		Page: &types.PageResponse{
			Total:           total,
			PageSize:        len(list),
			Current:         req.Current,
			RequestPageSize: req.PageSize,
		},
	}, nil
}
