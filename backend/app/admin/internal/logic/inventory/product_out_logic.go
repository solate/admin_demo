package inventory

import (
	"context"
	"time"

	"admin_backend/app/admin/internal/repository/inventoryrepo"
	"admin_backend/app/admin/internal/repository/productrepo"
	"admin_backend/app/admin/internal/repository/statisticsrepo"
	"admin_backend/app/admin/internal/svc"
	"admin_backend/app/admin/internal/types"
	"admin_backend/pkg/common/xerr"
	"admin_backend/pkg/ent"
	"admin_backend/pkg/ent/generated"
	"admin_backend/pkg/ent/generated/predicate"
	"admin_backend/pkg/ent/generated/product"
	"admin_backend/pkg/utils/idgen"

	"github.com/shopspring/decimal"
	"github.com/zeromicro/go-zero/core/logx"
)

type ProductOutLogic struct {
	logx.Logger
	ctx            context.Context
	svcCtx         *svc.ServiceContext
	inventoryRepo  *inventoryrepo.InventoryRepo
	productRepo    *productrepo.ProductRepo
	statisticsRepo *statisticsrepo.StatisticsRepo
}

// 商品出库
func NewProductOutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProductOutLogic {
	return &ProductOutLogic{
		Logger:         logx.WithContext(ctx),
		ctx:            ctx,
		svcCtx:         svcCtx,
		inventoryRepo:  inventoryrepo.NewInventoryRepo(svcCtx.DB),
		productRepo:    productrepo.NewProductRepo(svcCtx.DB),
		statisticsRepo: statisticsrepo.NewStatisticsRepo(svcCtx.DB),
	}
}

func (l *ProductOutLogic) ProductOut(req *types.ProductOutReq) (*types.InventoryOperationResp, error) {
	// 验证商品是否存在
	productInfo, err := l.productRepo.Get(l.ctx, []predicate.Product{product.ProductID(req.ProductID)})
	if err != nil {
		l.Error("ProductOut GetProduct err:", err.Error())
		return nil, xerr.NewErrCodeMsg(xerr.ServerError, "商品不存在")
	}

	// 检查商品状态
	if productInfo.Status != 1 {
		return nil, xerr.NewErrCodeMsg(xerr.ServerError, "商品已禁用，无法出库")
	}

	// 检查库存是否充足
	if productInfo.CurrentStock < req.Quantity {
		return nil, xerr.NewErrCodeMsg(xerr.ServerError, "库存不足，无法出库")
	}

	// 解析单价
	unitPrice, err := decimal.NewFromString(req.UnitPrice)
	if err != nil {
		return nil, xerr.NewErrCodeMsg(xerr.ServerError, "单价格式错误")
	}

	// 计算总金额
	totalAmount := unitPrice.Mul(decimal.NewFromInt(int64(req.Quantity)))

	// 生成库存记录ID
	inventoryID, err := idgen.GenerateUUID()
	if err != nil {
		return nil, xerr.NewErrCodeMsg(xerr.ServerError, "生成库存记录ID失败")
	}

	// 记录操作前的库存
	beforeStock := productInfo.CurrentStock
	afterStock := beforeStock - req.Quantity

	// 创建库存记录
	inventoryRecord := &generated.Inventory{
		InventoryID:   inventoryID,
		ProductID:     req.ProductID,
		OperationType: "out",
		Quantity:      req.Quantity,
		UnitPrice:     unitPrice,
		TotalAmount:   totalAmount,
		OperatorID:    req.OperatorID,
		Remark:        req.Remark,
		OperationTime: time.Now().UnixMilli(),
		BeforeStock:   beforeStock,
		AfterStock:    afterStock,
	}

	// 使用事务执行出库操作
	err = ent.WithTx(l.ctx, l.svcCtx.DB, func(tx *generated.Tx) error {
		// 保存库存记录
		_, err := l.inventoryRepo.CreateWithTx(l.ctx, tx, inventoryRecord)
		if err != nil {
			l.Error("ProductOut CreateInventory err:", err.Error())
			return xerr.NewErrCodeMsg(xerr.ServerError, "创建库存记录失败")
		}

		// 更新商品库存（出库是减少库存）
		_, err = l.productRepo.UpdateStockWithTx(l.ctx, tx, req.ProductID, -req.Quantity)
		if err != nil {
			l.Error("ProductOut UpdateStock err:", err.Error())
			return xerr.NewErrCodeMsg(xerr.ServerError, "更新商品库存失败")
		}

		// 更新库存操作统计
		err = l.statisticsRepo.IncrementInventoryStats(l.ctx, tx, "out", req.Quantity, totalAmount)
		if err != nil {
			l.Error("ProductOut IncrementInventoryStats err:", err.Error())
			return xerr.NewErrCodeMsg(xerr.ServerError, "更新库存统计失败")
		}

		// 更新库存数量统计（使用商品的采购价格计算库存价值变化，而不是出库单价）
		stockValueChange := productInfo.PurchasePrice.Mul(decimal.NewFromInt(int64(-req.Quantity)))
		lowStockChange := 0

		// 检查低库存状态变化
		oldLowStock := productInfo.CurrentStock <= productInfo.MinStock
		newLowStock := afterStock <= productInfo.MinStock
		if oldLowStock != newLowStock {
			if newLowStock {
				lowStockChange = 1
			} else {
				lowStockChange = -1
			}
		}

		err = l.statisticsRepo.UpdateStockStats(l.ctx, tx, -req.Quantity, stockValueChange, lowStockChange)
		if err != nil {
			l.Error("ProductOut UpdateStockStats err:", err.Error())
			return xerr.NewErrCodeMsg(xerr.ServerError, "更新库存统计失败")
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &types.InventoryOperationResp{
		InventoryID: inventoryID,
		Message:     "商品出库成功",
	}, nil
}
