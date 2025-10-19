package productrepo

import (
	"context"
	"time"

	"admin_backend/pkg/common"
	"admin_backend/pkg/common/contextutil"
	"admin_backend/pkg/ent"
	"admin_backend/pkg/ent/generated"
	"admin_backend/pkg/ent/generated/predicate"
	"admin_backend/pkg/ent/generated/product"
)

type ProductRepo struct {
	db *ent.Client
}

// NewProductRepo 创建商品仓储实例
func NewProductRepo(db *ent.Client) *ProductRepo {
	return &ProductRepo{db: db}
}

func (r *ProductRepo) Create(ctx context.Context, product *generated.Product) (*generated.Product, error) {
	now := time.Now().UnixMilli()
	// 自动从上下文获取租户代码
	tenantCode := contextutil.GetTenantCodeFromCtx(ctx)

	return r.db.Product.Create().
		SetCreatedAt(now).
		SetUpdatedAt(now).
		SetTenantCode(tenantCode). // 使用上下文中的租户代码
		SetProductID(product.ProductID).
		SetProductName(product.ProductName).
		SetUnit(product.Unit).
		SetPurchasePrice(product.PurchasePrice).
		SetSalePrice(product.SalePrice).
		SetCurrentStock(product.CurrentStock).
		SetMinStock(product.MinStock).
		SetStatus(product.Status).
		SetFactoryID(product.FactoryID).
		Save(ctx)
}

// CreateWithTx 在事务中创建商品
func (r *ProductRepo) CreateWithTx(ctx context.Context, tx *generated.Tx, product *generated.Product) (*generated.Product, error) {
	now := time.Now().UnixMilli()
	// 自动从上下文获取租户代码
	tenantCode := contextutil.GetTenantCodeFromCtx(ctx)

	return tx.Product.Create().
		SetCreatedAt(now).
		SetUpdatedAt(now).
		SetTenantCode(tenantCode). // 使用上下文中的租户代码
		SetProductID(product.ProductID).
		SetProductName(product.ProductName).
		SetUnit(product.Unit).
		SetPurchasePrice(product.PurchasePrice).
		SetSalePrice(product.SalePrice).
		SetCurrentStock(product.CurrentStock).
		SetMinStock(product.MinStock).
		SetStatus(product.Status).
		SetFactoryID(product.FactoryID).
		Save(ctx)
}

func (r *ProductRepo) Update(ctx context.Context, update *generated.Product) (int, error) {
	now := time.Now().UnixMilli()
	update.UpdatedAt = now

	// 添加租户过滤，确保只能更新当前租户的数据
	tenantCode := contextutil.GetTenantCodeFromCtx(ctx)
	query := r.db.Product.Update().
		SetUpdatedAt(now).
		Where(
			product.ProductID(update.ProductID),
			product.TenantCode(tenantCode), // 租户隔离，不允许跨租户更新
		)
	// 注意：不设置 TenantCode，防止租户被修改

	if update.ProductName != "" {
		query = query.SetProductName(update.ProductName)
	}
	if update.Unit != "" {
		query = query.SetUnit(update.Unit)
	}
	if !update.PurchasePrice.IsZero() {
		query = query.SetPurchasePrice(update.PurchasePrice)
	}
	if !update.SalePrice.IsZero() {
		query = query.SetSalePrice(update.SalePrice)
	}
	if update.CurrentStock != 0 {
		query = query.SetCurrentStock(update.CurrentStock)
	}
	if update.MinStock != 0 {
		query = query.SetMinStock(update.MinStock)
	}
	if update.Status != 0 {
		query = query.SetStatus(update.Status)
	}
	if update.FactoryID != "" {
		query = query.SetFactoryID(update.FactoryID)
	}

	return query.Save(ctx)
}

// UpdateWithTx 在事务中更新商品
func (r *ProductRepo) UpdateWithTx(ctx context.Context, tx *generated.Tx, update *generated.Product) (int, error) {
	now := time.Now().UnixMilli()
	update.UpdatedAt = now

	// 添加租户过滤，确保只能更新当前租户的数据
	tenantCode := contextutil.GetTenantCodeFromCtx(ctx)
	query := tx.Product.Update().
		SetUpdatedAt(now).
		Where(
			product.ProductID(update.ProductID),
			product.TenantCode(tenantCode), // 租户隔离，不允许跨租户更新
		)
	// 注意：不设置 TenantCode，防止租户被修改

	if update.ProductName != "" {
		query = query.SetProductName(update.ProductName)
	}
	if update.Unit != "" {
		query = query.SetUnit(update.Unit)
	}
	if !update.PurchasePrice.IsZero() {
		query = query.SetPurchasePrice(update.PurchasePrice)
	}
	if !update.SalePrice.IsZero() {
		query = query.SetSalePrice(update.SalePrice)
	}
	if update.CurrentStock != 0 {
		query = query.SetCurrentStock(update.CurrentStock)
	}
	if update.MinStock != 0 {
		query = query.SetMinStock(update.MinStock)
	}
	if update.Status != 0 {
		query = query.SetStatus(update.Status)
	}
	if update.FactoryID != "" {
		query = query.SetFactoryID(update.FactoryID)
	}

	return query.Save(ctx)
}

func (r *ProductRepo) GetByProductID(ctx context.Context, productID string) (*generated.Product, error) {
	return r.Get(ctx, []predicate.Product{product.ProductID(productID)})
}

// defaultQuery 默认查询条件
func (r *ProductRepo) defaultQuery(ctx context.Context, where []predicate.Product) []predicate.Product {
	where = append(where, product.DeletedAtIsNil())
	where = append(where, product.TenantCode(contextutil.GetTenantCodeFromCtx(ctx)))
	return where
}

func (r *ProductRepo) Get(ctx context.Context, where []predicate.Product) (*generated.Product, error) {
	where = r.defaultQuery(ctx, where)
	return r.db.Product.Query().Where(where...).Only(ctx)
}

func (r *ProductRepo) PageList(ctx context.Context, current, limit int, where []predicate.Product) ([]*generated.Product, int, error) {
	where = r.defaultQuery(ctx, where)

	offset := common.Offset(current, limit)
	query := r.db.Product.Query().Where(where...).Order(generated.Desc(product.FieldCreatedAt))

	// 查询总数
	total, err := query.Count(ctx)
	if err != nil || total == 0 {
		return nil, 0, err
	}

	// 分页查询
	list, err := query.Offset(offset).Limit(limit).All(ctx)
	return list, total, err
}

// Delete 根据商品ID删除商品，软删除
func (r *ProductRepo) Delete(ctx context.Context, productID string) (int, error) {
	now := time.Now().UnixMilli()
	return r.db.Product.Update().
		SetDeletedAt(now).
		Where(product.ProductID(productID)).Save(ctx)
}

// DeleteWithTx 在事务中删除商品，软删除
func (r *ProductRepo) DeleteWithTx(ctx context.Context, tx *generated.Tx, productID string) (int, error) {
	now := time.Now().UnixMilli()
	return tx.Product.Update().
		SetDeletedAt(now).
		Where(product.ProductID(productID)).Save(ctx)
}

// UpdateStock 更新商品库存
func (r *ProductRepo) UpdateStock(ctx context.Context, productID string, quantity int) (int, error) {
	now := time.Now().UnixMilli()
	return r.db.Product.Update().
		SetUpdatedAt(now).
		AddCurrentStock(quantity).
		Where(product.ProductID(productID)).Save(ctx)
}

// UpdateStockWithTx 在事务中更新商品库存
func (r *ProductRepo) UpdateStockWithTx(ctx context.Context, tx *generated.Tx, productID string, quantity int) (int, error) {
	now := time.Now().UnixMilli()
	return tx.Product.Update().
		SetUpdatedAt(now).
		AddCurrentStock(quantity).
		Where(product.ProductID(productID)).Save(ctx)
}

// GetLowStockProducts 获取低库存商品
func (r *ProductRepo) GetLowStockProducts(ctx context.Context) ([]*generated.Product, error) {
	where := r.defaultQuery(ctx, []predicate.Product{})
	// 这里需要手动实现低库存查询，因为需要比较两个字段
	products, err := r.db.Product.Query().Where(where...).All(ctx)
	if err != nil {
		return nil, err
	}

	var lowStockProducts []*generated.Product
	for _, p := range products {
		if p.CurrentStock <= p.MinStock {
			lowStockProducts = append(lowStockProducts, p)
		}
	}

	return lowStockProducts, nil
}

// GetProductsByFactory 根据工厂ID获取商品
func (r *ProductRepo) GetProductsByFactory(ctx context.Context, factoryID string) ([]*generated.Product, error) {
	where := r.defaultQuery(ctx, []predicate.Product{product.FactoryID(factoryID)})
	return r.db.Product.Query().Where(where...).All(ctx)
}

// GetAll 获取所有符合条件的商品
func (r *ProductRepo) GetAll(ctx context.Context, where []predicate.Product) ([]*generated.Product, error) {
	where = r.defaultQuery(ctx, where)
	return r.db.Product.Query().Where(where...).All(ctx)
}
