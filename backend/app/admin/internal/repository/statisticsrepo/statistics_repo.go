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

	// 先检查统计记录是否存在，不存在则创建
	exists, err := tx.ProductStatistics.Query().
		Where(productstatistics.TenantCode(tenantCode)).
		Exist(ctx)
	if err != nil {
		return err
	}

	if !exists {
		// 创建初始统计记录
		_, err = tx.ProductStatistics.Create().
			SetTenantCode(tenantCode).
			SetTotalProducts(0).
			SetActiveProducts(0).
			SetTotalStock(0).
			SetTotalStockValue(decimal.Zero).
			SetLowStockProducts(0).
			SetTotalInQuantity(0).
			SetTotalInAmount(decimal.Zero).
			SetTotalOutQuantity(0).
			SetTotalOutAmount(decimal.Zero).
			SetTotalSalesAmount(decimal.Zero).
			SetTotalSalesQuantity(0).
			Save(ctx)
		if err != nil {
			return err
		}
	}

	update := tx.ProductStatistics.Update().
		Where(productstatistics.TenantCode(tenantCode)).
		AddTotalProducts(1)

	if isActive {
		update = update.AddActiveProducts(1)
	}

	_, err = update.Save(ctx)
	return err
}

// DecrementProductCount 减少商品数量（在事务中调用）
func (r *StatisticsRepo) DecrementProductCount(ctx context.Context, tx *generated.Tx, isActive bool) error {
	tenantCode := contextutil.GetTenantCodeFromCtx(ctx)

	// 先检查统计记录是否存在，不存在则创建
	exists, err := tx.ProductStatistics.Query().
		Where(productstatistics.TenantCode(tenantCode)).
		Exist(ctx)
	if err != nil {
		return err
	}

	if !exists {
		// 创建初始统计记录
		_, err = tx.ProductStatistics.Create().
			SetTenantCode(tenantCode).
			SetTotalProducts(0).
			SetActiveProducts(0).
			SetTotalStock(0).
			SetTotalStockValue(decimal.Zero).
			SetLowStockProducts(0).
			SetTotalInQuantity(0).
			SetTotalInAmount(decimal.Zero).
			SetTotalOutQuantity(0).
			SetTotalOutAmount(decimal.Zero).
			SetTotalSalesAmount(decimal.Zero).
			SetTotalSalesQuantity(0).
			Save(ctx)
		if err != nil {
			return err
		}
	}

	update := tx.ProductStatistics.Update().
		Where(productstatistics.TenantCode(tenantCode)).
		AddTotalProducts(-1)

	if isActive {
		update = update.AddActiveProducts(-1)
	}

	_, err = update.Save(ctx)
	return err
}

// UpdateProductStatus 更新商品状态（在事务中调用）
func (r *StatisticsRepo) UpdateProductStatus(ctx context.Context, tx *generated.Tx, isActive bool) error {
	tenantCode := contextutil.GetTenantCodeFromCtx(ctx)

	// 先检查统计记录是否存在，不存在则创建
	exists, err := tx.ProductStatistics.Query().
		Where(productstatistics.TenantCode(tenantCode)).
		Exist(ctx)
	if err != nil {
		return err
	}

	if !exists {
		// 创建初始统计记录
		_, err = tx.ProductStatistics.Create().
			SetTenantCode(tenantCode).
			SetTotalProducts(0).
			SetActiveProducts(0).
			SetTotalStock(0).
			SetTotalStockValue(decimal.Zero).
			SetLowStockProducts(0).
			SetTotalInQuantity(0).
			SetTotalInAmount(decimal.Zero).
			SetTotalOutQuantity(0).
			SetTotalOutAmount(decimal.Zero).
			SetTotalSalesAmount(decimal.Zero).
			SetTotalSalesQuantity(0).
			Save(ctx)
		if err != nil {
			return err
		}
	}

	var delta int
	if isActive {
		delta = 1
	} else {
		delta = -1
	}

	_, err = tx.ProductStatistics.Update().
		Where(productstatistics.TenantCode(tenantCode)).
		AddActiveProducts(delta).
		Save(ctx)
	return err
}

// UpdateStockStats 更新库存统计（在事务中调用）
func (r *StatisticsRepo) UpdateStockStats(ctx context.Context, tx *generated.Tx, stockChange int, stockValueChange decimal.Decimal, lowStockChange int) error {
	tenantCode := contextutil.GetTenantCodeFromCtx(ctx)

	// 先检查统计记录是否存在，不存在则创建
	exists, err := tx.ProductStatistics.Query().
		Where(productstatistics.TenantCode(tenantCode)).
		Exist(ctx)
	if err != nil {
		return err
	}

	if !exists {
		// 创建初始统计记录
		_, err = tx.ProductStatistics.Create().
			SetTenantCode(tenantCode).
			SetTotalProducts(0).
			SetActiveProducts(0).
			SetTotalStock(0).
			SetTotalStockValue(decimal.Zero).
			SetLowStockProducts(0).
			SetTotalInQuantity(0).
			SetTotalInAmount(decimal.Zero).
			SetTotalOutQuantity(0).
			SetTotalOutAmount(decimal.Zero).
			SetTotalSalesAmount(decimal.Zero).
			SetTotalSalesQuantity(0).
			Save(ctx)
		if err != nil {
			return err
		}
	}

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

	// 先检查统计记录是否存在，不存在则创建
	exists, err := tx.ProductStatistics.Query().
		Where(productstatistics.TenantCode(tenantCode)).
		Exist(ctx)
	if err != nil {
		return err
	}

	if !exists {
		// 创建初始统计记录
		_, err = tx.ProductStatistics.Create().
			SetTenantCode(tenantCode).
			SetTotalProducts(0).
			SetActiveProducts(0).
			SetTotalStock(0).
			SetTotalStockValue(decimal.Zero).
			SetLowStockProducts(0).
			SetTotalInQuantity(0).
			SetTotalInAmount(decimal.Zero).
			SetTotalOutQuantity(0).
			SetTotalOutAmount(decimal.Zero).
			SetTotalSalesAmount(decimal.Zero).
			SetTotalSalesQuantity(0).
			Save(ctx)
		if err != nil {
			return err
		}
	}

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

// ProductDetailStat 商品明细统计结构
type ProductDetailStat struct {
	ProductID        string
	ProductName      string
	Unit             string
	CurrentStock     int
	MinStock         int
	TotalInQuantity  int
	TotalOutQuantity int
	PurchasePrice    decimal.Decimal
	SalePrice        decimal.Decimal
	StockValue       decimal.Decimal
	Status           int
}

// CalculateProductStatistics 计算商品统计
func (r *StatisticsRepo) CalculateProductStatistics(ctx context.Context) (map[string]interface{}, error) {
	tenantCode := contextutil.GetTenantCodeFromCtx(ctx)

	// 1. 获取所有商品
	products, err := r.db.Product.Query().
		Where(product.DeletedAtIsNil()).
		Where(product.TenantCode(tenantCode)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	totalProducts := len(products)
	var totalStock int
	var totalStockValue decimal.Decimal // 按采购价计算
	var totalSalesValue decimal.Decimal // 按销售价计算
	var lowStockProducts int            // 低库存商品数
	var totalInQuantity int             // 总入库数量
	var totalInAmount decimal.Decimal   // 总入库金额
	var totalOutQuantity int            // 总出库数量
	var totalOutAmount decimal.Decimal  // 总出库金额
	productDetails := make([]*ProductDetailStat, 0, len(products))

	// 2. 遍历商品，计算统计数据
	for _, p := range products {
		// 计算该商品的入库和出库数量及金额
		inQuantity, inAmount, outQuantity, outAmount, err := r.calculateProductInventoryStats(ctx, tenantCode, p.ProductID)
		if err != nil {
			return nil, err
		}

		// 累计总库存
		totalStock += p.CurrentStock

		// 计算库存价值（采购价 * 当前库存）
		stockValue := p.PurchasePrice.Mul(decimal.NewFromInt(int64(p.CurrentStock)))
		totalStockValue = totalStockValue.Add(stockValue)

		// 累计销售价值（销售价 * 当前库存）
		salesValue := p.SalePrice.Mul(decimal.NewFromInt(int64(p.CurrentStock)))
		totalSalesValue = totalSalesValue.Add(salesValue)

		// 统计低库存商品
		if p.CurrentStock <= p.MinStock {
			lowStockProducts++
		}

		// 累计总入库和出库
		totalInQuantity += inQuantity
		totalInAmount = totalInAmount.Add(inAmount)
		totalOutQuantity += outQuantity
		totalOutAmount = totalOutAmount.Add(outAmount)

		// 添加到商品明细列表
		productDetails = append(productDetails, &ProductDetailStat{
			ProductID:        p.ProductID,
			ProductName:      p.ProductName,
			Unit:             p.Unit,
			CurrentStock:     p.CurrentStock,
			MinStock:         p.MinStock,
			TotalInQuantity:  inQuantity,
			TotalOutQuantity: outQuantity,
			PurchasePrice:    p.PurchasePrice,
			SalePrice:        p.SalePrice,
			StockValue:       stockValue,
			Status:           p.Status,
		})
	}

	return map[string]interface{}{
		"total_products":      totalProducts,
		"total_stock":         totalStock,
		"total_stock_value":   totalStockValue,
		"total_sales_value":   totalSalesValue,
		"low_stock_products":  lowStockProducts,
		"total_in_quantity":   totalInQuantity,
		"total_in_amount":     totalInAmount,
		"total_out_quantity":  totalOutQuantity,
		"total_out_amount":    totalOutAmount,
		"product_detail_list": productDetails,
	}, nil
}

// calculateProductInventoryStats 计算单个商品的入库和出库统计
func (r *StatisticsRepo) calculateProductInventoryStats(ctx context.Context, tenantCode, productID string) (
	inQuantity int, inAmount decimal.Decimal, outQuantity int, outAmount decimal.Decimal, err error) {

	// 查询该商品的所有库存操作记录
	inventories, err := r.db.Inventory.Query().
		Where(inventory.TenantCode(tenantCode)).
		Where(inventory.ProductID(productID)).
		Where(inventory.DeletedAtIsNil()).
		All(ctx)
	if err != nil {
		return 0, decimal.Zero, 0, decimal.Zero, err
	}

	// 统计入库和出库数量及金额
	for _, inv := range inventories {
		if inv.OperationType == "in" {
			inQuantity += inv.Quantity
			inAmount = inAmount.Add(inv.TotalAmount)
		} else if inv.OperationType == "out" {
			outQuantity += inv.Quantity
			outAmount = outAmount.Add(inv.TotalAmount)
		}
	}

	return inQuantity, inAmount, outQuantity, outAmount, nil
}
