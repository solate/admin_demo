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

type ProductInLogic struct {
	logx.Logger
	ctx            context.Context
	svcCtx         *svc.ServiceContext
	inventoryRepo  *inventoryrepo.InventoryRepo
	productRepo    *productrepo.ProductRepo
	statisticsRepo *statisticsrepo.StatisticsRepo
}

// 商品入库
func NewProductInLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProductInLogic {
	return &ProductInLogic{
		Logger:         logx.WithContext(ctx),
		ctx:            ctx,
		svcCtx:         svcCtx,
		inventoryRepo:  inventoryrepo.NewInventoryRepo(svcCtx.DB),
		productRepo:    productrepo.NewProductRepo(svcCtx.DB),
		statisticsRepo: statisticsrepo.NewStatisticsRepo(svcCtx.DB),
	}
}

func (l *ProductInLogic) ProductIn(req *types.ProductInReq) (*types.InventoryOperationResp, error) {
	// 验证商品是否存在
	productInfo, err := l.productRepo.Get(l.ctx, []predicate.Product{product.ProductID(req.ProductID)})
	if err != nil {
		l.Error("ProductIn GetProduct err:", err.Error())
		return nil, xerr.NewErrCodeMsg(xerr.ServerError, "商品不存在")
	}

	// 检查商品状态
	if productInfo.Status != 1 {
		return nil, xerr.NewErrCodeMsg(xerr.ServerError, "商品已禁用，无法入库")
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
	afterStock := beforeStock + req.Quantity

	// 创建库存记录（租户信息由 Repository 层自动处理）
	inventoryRecord := &generated.Inventory{
		InventoryID:   inventoryID,
		ProductID:     req.ProductID,
		OperationType: "in",
		Quantity:      req.Quantity,
		UnitPrice:     unitPrice,
		TotalAmount:   totalAmount,
		OperatorID:    req.OperatorID,
		Remark:        req.Remark,
		OperationTime: time.Now().UnixMilli(),
		BeforeStock:   beforeStock,
		AfterStock:    afterStock,
	}

	// 使用事务执行入库操作
	err = ent.WithTx(l.ctx, l.svcCtx.DB, func(tx *generated.Tx) error {
		// 保存库存记录
		_, err := l.inventoryRepo.CreateWithTx(l.ctx, tx, inventoryRecord)
		if err != nil {
			l.Error("ProductIn CreateInventory err:", err.Error())
			return xerr.NewErrCodeMsg(xerr.ServerError, "创建库存记录失败")
		}

		// 更新商品库存
		_, err = l.productRepo.UpdateStockWithTx(l.ctx, tx, req.ProductID, req.Quantity)
		if err != nil {
			l.Error("ProductIn UpdateStock err:", err.Error())
			return xerr.NewErrCodeMsg(xerr.ServerError, "更新商品库存失败")
		}

		// 更新库存统计
		err = l.statisticsRepo.IncrementInventoryStats(l.ctx, tx, "in", req.Quantity, totalAmount)
		if err != nil {
			l.Error("ProductIn IncrementInventoryStats err:", err.Error())
			return xerr.NewErrCodeMsg(xerr.ServerError, "更新库存统计失败")
		}

		// 更新库存数量统计（使用商品的采购价格计算库存价值，而不是入库单价）
		stockValueChange := productInfo.PurchasePrice.Mul(decimal.NewFromInt(int64(req.Quantity)))

		// 检查低库存状态变化
		lowStockChange := 0
		oldLowStock := productInfo.CurrentStock <= productInfo.MinStock
		newLowStock := afterStock <= productInfo.MinStock
		if oldLowStock != newLowStock {
			if newLowStock {
				lowStockChange = 1
			} else {
				lowStockChange = -1
			}
		}

		err = l.statisticsRepo.UpdateStockStats(l.ctx, tx, req.Quantity, stockValueChange, lowStockChange)
		if err != nil {
			l.Error("ProductIn UpdateStockStats err:", err.Error())
			return xerr.NewErrCodeMsg(xerr.ServerError, "更新库存统计失败")
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &types.InventoryOperationResp{
		InventoryID: inventoryID,
		Message:     "商品入库成功",
	}, nil
}
