package factoryrepo

import (
	"context"
	"time"

	"admin_backend/pkg/common"
	"admin_backend/pkg/common/contextutil"
	"admin_backend/pkg/ent"
	"admin_backend/pkg/ent/generated"
	"admin_backend/pkg/ent/generated/factory"
	"admin_backend/pkg/ent/generated/predicate"
)

type FactoryRepo struct {
	db *ent.Client
}

// NewFactoryRepo 创建工厂仓储实例
func NewFactoryRepo(db *ent.Client) *FactoryRepo {
	return &FactoryRepo{db: db}
}

func (r *FactoryRepo) Create(ctx context.Context, factory *generated.Factory) (*generated.Factory, error) {
	now := time.Now().UnixMilli()
	// 自动从上下文获取租户代码
	tenantCode := contextutil.GetTenantCodeFromCtx(ctx)

	return r.db.Factory.Create().
		SetCreatedAt(now).
		SetUpdatedAt(now).
		SetTenantCode(tenantCode). // 使用上下文中的租户代码
		SetFactoryID(factory.FactoryID).
		SetFactoryName(factory.FactoryName).
		SetAddress(factory.Address).
		SetContactPhone(factory.ContactPhone).
		SetStatus(factory.Status).
		Save(ctx)
}

func (r *FactoryRepo) Update(ctx context.Context, update *generated.Factory) (int, error) {
	now := time.Now().UnixMilli()
	update.UpdatedAt = now

	// 添加租户过滤，确保只能更新当前租户的数据
	tenantCode := contextutil.GetTenantCodeFromCtx(ctx)
	query := r.db.Factory.Update().
		SetUpdatedAt(now).
		Where(
			factory.FactoryID(update.FactoryID),
			factory.TenantCode(tenantCode), // 租户隔离，不允许跨租户更新
		)
	// 注意：不设置 TenantCode，防止租户被修改

	if update.FactoryName != "" {
		query = query.SetFactoryName(update.FactoryName)
	}
	if update.Address != "" {
		query = query.SetAddress(update.Address)
	}
	if update.ContactPhone != "" {
		query = query.SetContactPhone(update.ContactPhone)
	}
	if update.Status != 0 {
		query = query.SetStatus(update.Status)
	}

	return query.Save(ctx)
}

func (r *FactoryRepo) GetByFactoryID(ctx context.Context, factoryID string) (*generated.Factory, error) {
	return r.Get(ctx, []predicate.Factory{factory.FactoryID(factoryID)})
}

// defaultQuery 默认查询条件
func (r *FactoryRepo) defaultQuery(ctx context.Context, where []predicate.Factory) []predicate.Factory {
	where = append(where, factory.DeletedAtIsNil())
	where = append(where, factory.TenantCode(contextutil.GetTenantCodeFromCtx(ctx)))
	return where
}

func (r *FactoryRepo) Get(ctx context.Context, where []predicate.Factory) (*generated.Factory, error) {
	where = r.defaultQuery(ctx, where)
	return r.db.Factory.Query().Where(where...).Only(ctx)
}

func (r *FactoryRepo) PageList(ctx context.Context, current, limit int, where []predicate.Factory) ([]*generated.Factory, int, error) {
	where = r.defaultQuery(ctx, where)

	offset := common.Offset(current, limit)
	query := r.db.Factory.Query().Where(where...).Order(
		generated.Desc(factory.FieldCreatedAt),
		generated.Desc(factory.FieldFactoryID),
	)

	// 查询总数
	total, err := query.Count(ctx)
	if err != nil || total == 0 {
		return nil, 0, err
	}

	// 分页查询
	list, err := query.Offset(offset).Limit(limit).All(ctx)
	return list, total, err
}

// Delete 根据工厂ID删除工厂，软删除
func (r *FactoryRepo) Delete(ctx context.Context, factoryID string) (int, error) {
	now := time.Now().UnixMilli()
	return r.db.Factory.Update().
		SetDeletedAt(now).
		Where(factory.FactoryID(factoryID)).Save(ctx)
}
