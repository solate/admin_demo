package initialize

import (
	"context"
	"fmt"
	"time"

	"admin_backend/pkg/ent/generated"
	"admin_backend/pkg/utils/idgen"

	"github.com/shopspring/decimal"
)

// InitInventory 初始化库存操作记录数据
func InitInventory(ctx context.Context, tx *generated.Tx, tenantCode string, productIDs []string) error {
	// 生成库存记录ID
	ids, err := idgen.GenerateUUIDs(10)
	if err != nil {
		return err
	}

	now := time.Now().UnixMilli()
	baseTime := now - 24*60*60*1000 // 24小时前

	inventoryOps := []struct {
		id            string
		productID     string
		opType        string
		quantity      int
		unitPrice     decimal.Decimal
		totalAmount   decimal.Decimal
		operatorID    string
		remark        string
		operationTime int64
		beforeStock   int
		afterStock    int
	}{
		// 商品1的操作记录
		{ids[0], productIDs[0], "in", 10, decimal.NewFromInt(100), decimal.NewFromInt(1000), "", "初始入库", baseTime, 0, 10},
		{ids[1], productIDs[0], "in", 10, decimal.NewFromInt(100), decimal.NewFromInt(1000), "", "补货", baseTime + 2*60*60*1000, 10, 20},

		// 商品2的操作记录
		{ids[2], productIDs[1], "in", 15, decimal.NewFromInt(50), decimal.NewFromInt(750), "", "初始入库", baseTime, 0, 15},

		// 商品3的操作记录（低库存）
		{ids[3], productIDs[2], "in", 8, decimal.NewFromInt(200), decimal.NewFromInt(1600), "", "初始入库", baseTime, 0, 8},
		{ids[4], productIDs[2], "out", 5, decimal.NewFromInt(200), decimal.NewFromInt(1000), "", "销售", baseTime + 60*60*1000, 8, 3},

		// 商品4的操作记录
		{ids[5], productIDs[3], "in", 50, decimal.NewFromInt(75), decimal.NewFromInt(3750), "", "初始入库", baseTime, 0, 50},

		// 商品5的操作记录（低库存）
		{ids[6], productIDs[4], "in", 5, decimal.NewFromInt(150), decimal.NewFromInt(750), "", "初始入库", baseTime, 0, 5},
		{ids[7], productIDs[4], "out", 3, decimal.NewFromInt(150), decimal.NewFromInt(450), "", "销售", baseTime + 3*60*60*1000, 5, 2},

		// 商品6的操作记录
		{ids[8], productIDs[5], "in", 100, decimal.NewFromInt(30), decimal.NewFromInt(3000), "", "初始入库", baseTime, 0, 100},
		{ids[9], productIDs[5], "out", 0, decimal.NewFromInt(30), decimal.Zero, "", "测试", baseTime + 1*60*60*1000, 100, 100},
	}

	for _, inv := range inventoryOps {
		_, err := tx.Inventory.Create().
			SetCreatedAt(now).
			SetUpdatedAt(now).
			SetTenantCode(tenantCode).
			SetInventoryID(inv.id).
			SetProductID(inv.productID).
			SetOperationType(inv.opType).
			SetQuantity(inv.quantity).
			SetUnitPrice(inv.unitPrice).
			SetTotalAmount(inv.totalAmount).
			SetOperatorID(inv.operatorID).
			SetRemark(inv.remark).
			SetOperationTime(inv.operationTime).
			SetBeforeStock(inv.beforeStock).
			SetAfterStock(inv.afterStock).
			Save(ctx)
		if err != nil {
			return fmt.Errorf("failed creating inventory record: %v", err)
		}
		fmt.Printf("inventory: %s, product: %s, type: %s, quantity: %d\n", inv.id, inv.productID, inv.opType, inv.quantity)
	}

	return nil
}
