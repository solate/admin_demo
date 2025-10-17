package inventoryrepo

import (
	"context"
	"time"

	"admin_backend/pkg/common"
	"admin_backend/pkg/common/contextutil"
	"admin_backend/pkg/ent"
	"admin_backend/pkg/ent/generated"
	"admin_backend/pkg/ent/generated/inventory"
	"admin_backend/pkg/ent/generated/predicate"

	"github.com/shopspring/decimal"
)

type InventoryRepo struct {
	db *ent.Client
}

// NewInventoryRepo 创建库存仓储实例
func NewInventoryRepo(db *ent.Client) *InventoryRepo {
	return &InventoryRepo{db: db}
}

func (r *InventoryRepo) Create(ctx context.Context, inventory *generated.Inventory) (*generated.Inventory, error) {
	now := time.Now().UnixMilli()
	return r.db.Inventory.Create().
		SetCreatedAt(now).
		SetUpdatedAt(now).
		SetTenantCode(inventory.TenantCode).
		SetInventoryID(inventory.InventoryID).
		SetProductID(inventory.ProductID).
		SetOperationType(inventory.OperationType).
		SetQuantity(inventory.Quantity).
		SetUnitPrice(inventory.UnitPrice).
		SetTotalAmount(inventory.TotalAmount).
		SetOperatorID(inventory.OperatorID).
		SetRemark(inventory.Remark).
		SetOperationTime(inventory.OperationTime).
		SetBeforeStock(inventory.BeforeStock).
		SetAfterStock(inventory.AfterStock).
		Save(ctx)
}

// CreateWithTx 在事务中创建库存记录
func (r *InventoryRepo) CreateWithTx(ctx context.Context, tx *generated.Tx, inventory *generated.Inventory) (*generated.Inventory, error) {
	now := time.Now().UnixMilli()
	return tx.Inventory.Create().
		SetCreatedAt(now).
		SetUpdatedAt(now).
		SetTenantCode(inventory.TenantCode).
		SetInventoryID(inventory.InventoryID).
		SetProductID(inventory.ProductID).
		SetOperationType(inventory.OperationType).
		SetQuantity(inventory.Quantity).
		SetUnitPrice(inventory.UnitPrice).
		SetTotalAmount(inventory.TotalAmount).
		SetOperatorID(inventory.OperatorID).
		SetRemark(inventory.Remark).
		SetOperationTime(inventory.OperationTime).
		SetBeforeStock(inventory.BeforeStock).
		SetAfterStock(inventory.AfterStock).
		Save(ctx)
}

// defaultQuery 默认查询条件
func (r *InventoryRepo) defaultQuery(ctx context.Context, where []predicate.Inventory) []predicate.Inventory {
	where = append(where, inventory.DeletedAtIsNil())
	where = append(where, inventory.TenantCode(contextutil.GetTenantCodeFromCtx(ctx)))
	return where
}

func (r *InventoryRepo) Get(ctx context.Context, where []predicate.Inventory) (*generated.Inventory, error) {
	where = r.defaultQuery(ctx, where)
	return r.db.Inventory.Query().Where(where...).Only(ctx)
}

func (r *InventoryRepo) PageList(ctx context.Context, current, limit int, where []predicate.Inventory) ([]*generated.Inventory, int, error) {
	where = r.defaultQuery(ctx, where)

	offset := common.Offset(current, limit)
	query := r.db.Inventory.Query().Where(where...).Order(generated.Desc(inventory.FieldOperationTime))

	// 查询总数
	total, err := query.Count(ctx)
	if err != nil || total == 0 {
		return nil, 0, err
	}

	// 分页查询
	list, err := query.Offset(offset).Limit(limit).All(ctx)
	return list, total, err
}

// GetByProductID 根据商品ID获取库存记录
func (r *InventoryRepo) GetByProductID(ctx context.Context, productID string) ([]*generated.Inventory, error) {
	where := r.defaultQuery(ctx, []predicate.Inventory{inventory.ProductID(productID)})
	return r.db.Inventory.Query().Where(where...).Order(generated.Desc(inventory.FieldOperationTime)).All(ctx)
}

// GetByOperationType 根据操作类型获取库存记录
func (r *InventoryRepo) GetByOperationType(ctx context.Context, operationType string) ([]*generated.Inventory, error) {
	where := r.defaultQuery(ctx, []predicate.Inventory{inventory.OperationType(operationType)})
	return r.db.Inventory.Query().Where(where...).Order(generated.Desc(inventory.FieldOperationTime)).All(ctx)
}

// GetByTimeRange 根据时间范围获取库存记录
func (r *InventoryRepo) GetByTimeRange(ctx context.Context, startTime, endTime int64) ([]*generated.Inventory, error) {
	where := r.defaultQuery(ctx, []predicate.Inventory{
		inventory.OperationTimeGTE(startTime),
		inventory.OperationTimeLTE(endTime),
	})
	return r.db.Inventory.Query().Where(where...).Order(generated.Desc(inventory.FieldOperationTime)).All(ctx)
}

// GetStatisticsByProduct 根据商品ID获取统计信息
func (r *InventoryRepo) GetStatisticsByProduct(ctx context.Context, productID string, startTime, endTime int64) (map[string]interface{}, error) {
	where := r.defaultQuery(ctx, []predicate.Inventory{
		inventory.ProductID(productID),
		inventory.OperationTimeGTE(startTime),
		inventory.OperationTimeLTE(endTime),
	})

	query := r.db.Inventory.Query().Where(where...)

	// 入库统计
	inRecords, err := query.Clone().Where(inventory.OperationType("in")).All(ctx)
	if err != nil {
		return nil, err
	}

	// 出库统计
	outRecords, err := query.Clone().Where(inventory.OperationType("out")).All(ctx)
	if err != nil {
		return nil, err
	}

	// 计算入库总量和金额
	var inQuantity int
	var inAmount decimal.Decimal
	for _, record := range inRecords {
		inQuantity += record.Quantity
		inAmount = inAmount.Add(record.TotalAmount)
	}

	// 计算出库总量和金额
	var outQuantity int
	var outAmount decimal.Decimal
	for _, record := range outRecords {
		outQuantity += record.Quantity
		outAmount = outAmount.Add(record.TotalAmount)
	}

	return map[string]interface{}{
		"in_quantity":  inQuantity,
		"in_amount":    inAmount,
		"out_quantity": outQuantity,
		"out_amount":   outAmount,
	}, nil
}
