package tenantrepo

import (
	"context"
	"time"

	"admin_backend/pkg/common"
	"admin_backend/pkg/ent"
	"admin_backend/pkg/ent/generated"
	"admin_backend/pkg/ent/generated/predicate"
	"admin_backend/pkg/ent/generated/tenant"
)

type TenantRepo struct {
	db *ent.Client
}

// NewTenantRepo 创建租户仓储实例
func NewTenantRepo(db *ent.Client) *TenantRepo {
	return &TenantRepo{db: db}
}

func (r *TenantRepo) Create(ctx context.Context, tenant *generated.Tenant) (*generated.Tenant, error) {
	now := time.Now().UnixMilli()
	return r.db.Tenant.Create().
		SetCreatedAt(now).
		SetUpdatedAt(now).
		SetTenantID(tenant.TenantID).
		SetCode(tenant.Code).
		SetName(tenant.Name).
		SetDescription(tenant.Description).
		SetStatus(tenant.Status).
		Save(ctx)
}

func (r *TenantRepo) Update(ctx context.Context, update *generated.Tenant) (int, error) {
	now := time.Now().UnixMilli()
	update.UpdatedAt = now
	return r.db.Tenant.Update().
		SetStatus(update.Status).
		SetDescription(update.Description).
		SetName(update.Name).
		Where(tenant.TenantID(update.TenantID)).Save(ctx)
}

func (r *TenantRepo) GetByTenantID(ctx context.Context, tenantID string) (*generated.Tenant, error) {
	return r.Get(ctx, []predicate.Tenant{tenant.TenantID(tenantID)})
}

// defaultQuery 默认查询条件
func (r *TenantRepo) defaultQuery(where []predicate.Tenant) []predicate.Tenant {
	where = append(where, tenant.DeletedAtIsNil())
	return where
}
func (r *TenantRepo) Get(ctx context.Context, where []predicate.Tenant) (*generated.Tenant, error) {
	where = r.defaultQuery(where)
	return r.db.Tenant.Query().Where(where...).Only(ctx)
}

func (r *TenantRepo) PageList(ctx context.Context, current, limit int, where []predicate.Tenant) ([]*generated.Tenant, int, error) {

	where = r.defaultQuery(where)

	offset := common.Offset(current, limit)
	query := r.db.Tenant.Query().Where(where...).Order(generated.Desc(tenant.FieldCreatedAt))

	// 查询总数
	total, err := query.Count(ctx)
	if err != nil || total == 0 {
		return nil, 0, err
	}

	// 分页查询
	list, err := query.Offset(offset).Limit(limit).All(ctx)
	return list, total, err
}

func (r *TenantRepo) DeleteByTenantID(ctx context.Context, delete *generated.Tenant) (int, error) {
	now := time.Now().UnixMilli()
	delete.DeletedAt = &now
	return r.db.Tenant.Update().
		SetDeletedAt(now).
		Where(tenant.TenantID(delete.TenantID)).Save(ctx)
}
