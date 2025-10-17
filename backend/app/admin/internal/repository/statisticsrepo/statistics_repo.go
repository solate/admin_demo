package statisticsrepo

import (
	"context"

	"admin_backend/pkg/common/contextutil"
	"admin_backend/pkg/ent"
	"admin_backend/pkg/ent/generated"
	"admin_backend/pkg/ent/generated/inventory"
	"admin_backend/pkg/ent/generated/product"
	"admin_backend/pkg/ent/generated/productstatistics"

	"github.com/shopspring/decimal"
)

type StatisticsRepo struct {
	db *ent.Client
}

// NewStatisticsRepo 创建统计仓储实例
func NewStatisticsRepo(db *ent.Client) *StatisticsRepo {
	return &StatisticsRepo{db: db}
}

// GetOrCreateStatistics 获取或创建统计数据（每个租户只有一条记录）
func (r *StatisticsRepo) GetOrCreateStatistics(ctx context.Context) (*generated.ProductStatistics, error) {
	tenantCode := contextutil.GetTenantCodeFromCtx(ctx)

	// 先尝试获取现有记录
	statistics, err := r.db.ProductStatistics.Query().
		Where(productstatistics.TenantCode(tenantCode)).
		Only(ctx)

	if err == nil {
		return statistics, nil
	}

	// 如果不存在，创建新记录
	return r.CreateStatistics(ctx)
}

// CreateStatistics 创建统计数据记录
func (r *StatisticsRepo) CreateStatistics(ctx context.Context) (*generated.ProductStatistics, error) {
	tenantCode := contextutil.GetTenantCodeFromCtx(ctx)

	// 计算统计数据
	stats, err := r.CalculateProductStatistics(ctx)
	if err != nil {
		return nil, err
	}

	return r.db.ProductStatistics.Create().
		SetTenantCode(tenantCode).
		SetTotalProducts(stats["total_products"].(int)).
		SetActiveProducts(stats["active_products"].(int)).
		SetTotalStock(stats["total_stock"].(int)).
		SetTotalStockValue(stats["total_stock_value"].(decimal.Decimal)).
		SetLowStockProducts(stats["low_stock_products"].(int)).
		SetTotalInQuantity(stats["total_in_quantity"].(int)).
		SetTotalInAmount(stats["total_in_amount"].(decimal.Decimal)).
		SetTotalOutQuantity(stats["total_out_quantity"].(int)).
		SetTotalOutAmount(stats["total_out_amount"].(decimal.Decimal)).
		SetTotalSalesAmount(stats["total_sales_amount"].(decimal.Decimal)).
		SetTotalSalesQuantity(stats["total_sales_quantity"].(int)).
		Save(ctx)
}

// UpdateStatistics 更新统计数据
func (r *StatisticsRepo) UpdateStatistics(ctx context.Context) (*generated.ProductStatistics, error) {
	tenantCode := contextutil.GetTenantCodeFromCtx(ctx)

	// 计算最新统计数据
	stats, err := r.CalculateProductStatistics(ctx)
	if err != nil {
		return nil, err
	}

	// 更新记录
	_, err = r.db.ProductStatistics.Update().
		Where(productstatistics.TenantCode(tenantCode)).
		SetTotalProducts(stats["total_products"].(int)).
		SetActiveProducts(stats["active_products"].(int)).
		SetTotalStock(stats["total_stock"].(int)).
		SetTotalStockValue(stats["total_stock_value"].(decimal.Decimal)).
		SetLowStockProducts(stats["low_stock_products"].(int)).
		SetTotalInQuantity(stats["total_in_quantity"].(int)).
		SetTotalInAmount(stats["total_in_amount"].(decimal.Decimal)).
		SetTotalOutQuantity(stats["total_out_quantity"].(int)).
		SetTotalOutAmount(stats["total_out_amount"].(decimal.Decimal)).
		SetTotalSalesAmount(stats["total_sales_amount"].(decimal.Decimal)).
		SetTotalSalesQuantity(stats["total_sales_quantity"].(int)).
		Save(ctx)

	if err != nil {
		return nil, err
	}

	// 返回更新后的记录
	return r.db.ProductStatistics.Query().
		Where(productstatistics.TenantCode(tenantCode)).
		Only(ctx)
}

// IncrementProductCount 增加商品数量（在事务中调用）
func (r *StatisticsRepo) IncrementProductCount(ctx context.Context, tx *generated.Tx, isActive bool) error {
	tenantCode := contextutil.GetTenantCodeFromCtx(ctx)

	update := tx.ProductStatistics.Update().
		Where(productstatistics.TenantCode(tenantCode)).
		AddTotalProducts(1)

	if isActive {
		update = update.AddActiveProducts(1)
	}

	_, err := update.Save(ctx)
	return err
}

// DecrementProductCount 减少商品数量（在事务中调用）
func (r *StatisticsRepo) DecrementProductCount(ctx context.Context, tx *generated.Tx, isActive bool) error {
	tenantCode := contextutil.GetTenantCodeFromCtx(ctx)

	update := tx.ProductStatistics.Update().
		Where(productstatistics.TenantCode(tenantCode)).
		AddTotalProducts(-1)

	if isActive {
		update = update.AddActiveProducts(-1)
	}

	_, err := update.Save(ctx)
	return err
}

// UpdateProductStatus 更新商品状态（在事务中调用）
func (r *StatisticsRepo) UpdateProductStatus(ctx context.Context, tx *generated.Tx, isActive bool) error {
	tenantCode := contextutil.GetTenantCodeFromCtx(ctx)

	var delta int
	if isActive {
		delta = 1
	} else {
		delta = -1
	}

	_, err := tx.ProductStatistics.Update().
		Where(productstatistics.TenantCode(tenantCode)).
		AddActiveProducts(delta).
		Save(ctx)
	return err
}

// UpdateStockStats 更新库存统计（在事务中调用）
func (r *StatisticsRepo) UpdateStockStats(ctx context.Context, tx *generated.Tx, stockChange int, stockValueChange decimal.Decimal, lowStockChange int) error {
	tenantCode := contextutil.GetTenantCodeFromCtx(ctx)

	// 先获取当前统计记录
	currentStats, err := tx.ProductStatistics.Query().
		Where(productstatistics.TenantCode(tenantCode)).
		Only(ctx)
	if err != nil {
		return err
	}

	// 计算新的值
	newTotalStock := currentStats.TotalStock + stockChange
	newTotalStockValue := currentStats.TotalStockValue.Add(stockValueChange)
	newLowStockProducts := currentStats.LowStockProducts + lowStockChange

	_, err = tx.ProductStatistics.Update().
		Where(productstatistics.TenantCode(tenantCode)).
		SetTotalStock(newTotalStock).
		SetTotalStockValue(newTotalStockValue).
		SetLowStockProducts(newLowStockProducts).
		Save(ctx)
	return err
}

// IncrementInventoryStats 增加库存操作统计（在事务中调用）
func (r *StatisticsRepo) IncrementInventoryStats(ctx context.Context, tx *generated.Tx, operationType string, quantity int, amount decimal.Decimal) error {
	tenantCode := contextutil.GetTenantCodeFromCtx(ctx)

	// 先获取当前统计记录
	currentStats, err := tx.ProductStatistics.Query().
		Where(productstatistics.TenantCode(tenantCode)).
		Only(ctx)
	if err != nil {
		return err
	}

	// 计算新的值
	var newTotalInQuantity, newTotalOutQuantity, newTotalSalesQuantity int
	var newTotalInAmount, newTotalOutAmount, newTotalSalesAmount decimal.Decimal

	if operationType == "in" {
		newTotalInQuantity = currentStats.TotalInQuantity + quantity
		newTotalInAmount = currentStats.TotalInAmount.Add(amount)
		newTotalOutQuantity = currentStats.TotalOutQuantity
		newTotalOutAmount = currentStats.TotalOutAmount
		newTotalSalesQuantity = currentStats.TotalSalesQuantity
		newTotalSalesAmount = currentStats.TotalSalesAmount
	} else if operationType == "out" {
		newTotalInQuantity = currentStats.TotalInQuantity
		newTotalInAmount = currentStats.TotalInAmount
		newTotalOutQuantity = currentStats.TotalOutQuantity + quantity
		newTotalOutAmount = currentStats.TotalOutAmount.Add(amount)
		newTotalSalesQuantity = currentStats.TotalSalesQuantity + quantity
		newTotalSalesAmount = currentStats.TotalSalesAmount.Add(amount)
	}

	_, err = tx.ProductStatistics.Update().
		Where(productstatistics.TenantCode(tenantCode)).
		SetTotalInQuantity(newTotalInQuantity).
		SetTotalInAmount(newTotalInAmount).
		SetTotalOutQuantity(newTotalOutQuantity).
		SetTotalOutAmount(newTotalOutAmount).
		SetTotalSalesQuantity(newTotalSalesQuantity).
		SetTotalSalesAmount(newTotalSalesAmount).
		Save(ctx)
	return err
}

// CalculateProductStatistics 计算商品统计
func (r *StatisticsRepo) CalculateProductStatistics(ctx context.Context) (map[string]interface{}, error) {
	tenantCode := contextutil.GetTenantCodeFromCtx(ctx)

	// 1. 获取商品统计
	productStats, err := r.calculateProductStats(ctx, tenantCode)
	if err != nil {
		return nil, err
	}

	// 2. 获取库存操作统计
	inventoryStats, err := r.calculateInventoryStats(ctx, tenantCode)
	if err != nil {
		return nil, err
	}

	// 合并统计结果
	result := make(map[string]interface{})
	for k, v := range productStats {
		result[k] = v
	}
	for k, v := range inventoryStats {
		result[k] = v
	}

	return result, nil
}

// calculateProductStats 计算商品相关统计
func (r *StatisticsRepo) calculateProductStats(ctx context.Context, tenantCode string) (map[string]interface{}, error) {
	// 获取商品总数
	totalProducts, err := r.db.Product.Query().
		Where(product.DeletedAtIsNil()).
		Where(product.TenantCode(tenantCode)).
		Count(ctx)
	if err != nil {
		return nil, err
	}

	// 获取启用商品数
	activeProducts, err := r.db.Product.Query().
		Where(product.DeletedAtIsNil()).
		Where(product.TenantCode(tenantCode)).
		Where(product.Status(1)).
		Count(ctx)
	if err != nil {
		return nil, err
	}

	// 获取所有商品进行库存统计
	products, err := r.db.Product.Query().
		Where(product.DeletedAtIsNil()).
		Where(product.TenantCode(tenantCode)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	var totalStock int
	var totalStockValue decimal.Decimal
	var lowStockProducts int

	for _, product := range products {
		totalStock += product.CurrentStock
		stockValue := product.PurchasePrice.Mul(decimal.NewFromInt(int64(product.CurrentStock)))
		totalStockValue = totalStockValue.Add(stockValue)

		if product.CurrentStock <= product.MinStock {
			lowStockProducts++
		}
	}

	return map[string]interface{}{
		"total_products":     totalProducts,
		"active_products":    activeProducts,
		"total_stock":        totalStock,
		"total_stock_value":  totalStockValue,
		"low_stock_products": lowStockProducts,
	}, nil
}

// calculateInventoryStats 计算库存操作统计
func (r *StatisticsRepo) calculateInventoryStats(ctx context.Context, tenantCode string) (map[string]interface{}, error) {
	// 获取入库统计
	inRecords, err := r.db.Inventory.Query().
		Where(inventory.TenantCode(tenantCode)).
		Where(inventory.OperationType("in")).
		All(ctx)
	if err != nil {
		return nil, err
	}

	var totalInQuantity int
	var totalInAmount decimal.Decimal
	for _, record := range inRecords {
		totalInQuantity += record.Quantity
		totalInAmount = totalInAmount.Add(record.TotalAmount)
	}

	// 获取出库统计
	outRecords, err := r.db.Inventory.Query().
		Where(inventory.TenantCode(tenantCode)).
		Where(inventory.OperationType("out")).
		All(ctx)
	if err != nil {
		return nil, err
	}

	var totalOutQuantity int
	var totalOutAmount decimal.Decimal
	for _, record := range outRecords {
		totalOutQuantity += record.Quantity
		totalOutAmount = totalOutAmount.Add(record.TotalAmount)
	}

	// 销售统计（出库即销售）
	totalSalesAmount := totalOutAmount
	totalSalesQuantity := totalOutQuantity

	return map[string]interface{}{
		"total_in_quantity":    totalInQuantity,
		"total_in_amount":      totalInAmount,
		"total_out_quantity":   totalOutQuantity,
		"total_out_amount":     totalOutAmount,
		"total_sales_amount":   totalSalesAmount,
		"total_sales_quantity": totalSalesQuantity,
	}, nil
}
