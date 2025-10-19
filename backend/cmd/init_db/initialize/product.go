package initialize

import (
	"context"
	"fmt"
	"time"

	"admin_backend/pkg/ent/generated"
	"admin_backend/pkg/utils/idgen"

	"github.com/shopspring/decimal"
)

// InitProduct 初始化商品数据
func InitProduct(ctx context.Context, tx *generated.Tx, tenantCode string, factoryIDs []string) ([]string, error) {
	// 生成商品ID
	ids, err := idgen.GenerateUUIDs(6)
	if err != nil {
		return nil, err
	}

	now := time.Now().UnixMilli()
	products := []struct {
		id            string
		name          string
		unit          string
		purchasePrice decimal.Decimal
		salePrice     decimal.Decimal
		currentStock  int
		minStock      int
		factoryID     string
	}{
		{ids[0], "测试商品1", "个", decimal.NewFromInt(100), decimal.NewFromInt(150), 20, 5, factoryIDs[0]},
		{ids[1], "测试商品2", "盒", decimal.NewFromInt(50), decimal.NewFromInt(80), 15, 10, factoryIDs[0]},
		{ids[2], "测试商品3", "件", decimal.NewFromInt(200), decimal.NewFromInt(300), 3, 5, factoryIDs[1]},
		{ids[3], "测试商品4", "个", decimal.NewFromInt(75), decimal.NewFromInt(120), 50, 10, factoryIDs[1]},
		{ids[4], "测试商品5", "套", decimal.NewFromInt(150), decimal.NewFromInt(250), 2, 5, factoryIDs[2]},
		{ids[5], "测试商品6", "个", decimal.NewFromInt(30), decimal.NewFromInt(50), 100, 20, factoryIDs[2]},
	}

	var productIDs []string
	for _, p := range products {
		_, err := tx.Product.Create().
			SetCreatedAt(now).
			SetUpdatedAt(now).
			SetTenantCode(tenantCode).
			SetProductID(p.id).
			SetProductName(p.name).
			SetUnit(p.unit).
			SetPurchasePrice(p.purchasePrice).
			SetSalePrice(p.salePrice).
			SetCurrentStock(p.currentStock).
			SetMinStock(p.minStock).
			SetStatus(1).
			SetFactoryID(p.factoryID).
			Save(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed creating product: %v", err)
		}
		productIDs = append(productIDs, p.id)
		fmt.Printf("product: %s, name: %s, stock: %d, minStock: %d\n", p.id, p.name, p.currentStock, p.minStock)
	}

	return productIDs, nil
}
