package product

import (
	"context"

	"admin_backend/app/admin/internal/repository/productrepo"
	"admin_backend/app/admin/internal/repository/statisticsrepo"
	"admin_backend/app/admin/internal/svc"
	"admin_backend/app/admin/internal/types"
	"admin_backend/pkg/common/xerr"
	"admin_backend/pkg/ent"
	"admin_backend/pkg/ent/generated"

	"github.com/shopspring/decimal"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateProductLogic struct {
	logx.Logger
	ctx            context.Context
	svcCtx         *svc.ServiceContext
	productRepo    *productrepo.ProductRepo
	statisticsRepo *statisticsrepo.StatisticsRepo
}

// 更新商品
func NewUpdateProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateProductLogic {
	return &UpdateProductLogic{
		Logger:         logx.WithContext(ctx),
		ctx:            ctx,
		svcCtx:         svcCtx,
		productRepo:    productrepo.NewProductRepo(svcCtx.DB),
		statisticsRepo: statisticsrepo.NewStatisticsRepo(svcCtx.DB),
	}
}

func (l *UpdateProductLogic) UpdateProduct(req *types.UpdateProductReq) (bool, error) {
	// 检查商品是否存在，获取旧数据
	oldProduct, err := l.productRepo.GetByProductID(l.ctx, req.ProductID)
	if err != nil {
		l.Error("UpdateProduct GetByProductID err:", err.Error())
		return false, xerr.NewErrCodeMsg(xerr.ServerError, "商品不存在")
	}

	// 构建更新数据（租户信息由 Repository 层自动处理）
	updateProduct := &generated.Product{
		ProductID:    req.ProductID,
		ProductName:  req.ProductName,
		Unit:         req.Unit,
		CurrentStock: req.CurrentStock,
		MinStock:     req.MinStock,
		Status:       req.Status,
		FactoryID:    req.FactoryID,
	}

	// 处理价格字段
	if req.PurchasePrice != "" {
		purchasePrice, err := decimal.NewFromString(req.PurchasePrice)
		if err != nil {
			return false, xerr.NewErrCodeMsg(xerr.ServerError, "采购价格格式错误")
		}
		updateProduct.PurchasePrice = purchasePrice
	} else {
		updateProduct.PurchasePrice = oldProduct.PurchasePrice
	}

	if req.SalePrice != "" {
		salePrice, err := decimal.NewFromString(req.SalePrice)
		if err != nil {
			return false, xerr.NewErrCodeMsg(xerr.ServerError, "销售价格格式错误")
		}
		updateProduct.SalePrice = salePrice
	} else {
		updateProduct.SalePrice = oldProduct.SalePrice
	}

	// 使用事务更新商品和统计
	err = ent.WithTx(l.ctx, l.svcCtx.DB, func(tx *generated.Tx) error {
		// 更新商品
		_, err := l.productRepo.UpdateWithTx(l.ctx, tx, updateProduct)
		if err != nil {
			l.Error("UpdateProduct Update err:", err.Error())
			return xerr.NewErrCodeMsg(xerr.ServerError, "更新商品失败")
		}

		// 检查状态是否变更
		if oldProduct.Status != req.Status {
			err = l.statisticsRepo.UpdateProductStatus(l.ctx, tx, req.Status == 1)
			if err != nil {
				l.Error("UpdateProduct UpdateProductStatus err:", err.Error())
				return xerr.NewErrCodeMsg(xerr.ServerError, "更新商品状态统计失败")
			}
		}

		// 检查库存或采购价格是否变更
		stockChange := req.CurrentStock - oldProduct.CurrentStock
		priceChanged := !updateProduct.PurchasePrice.Equal(oldProduct.PurchasePrice)

		// 如果库存数量或采购价格发生变化，需要更新库存统计
		if stockChange != 0 || priceChanged {
			var stockValueChange decimal.Decimal

			if priceChanged && stockChange != 0 {
				// 库存和价格都变化：新库存价值 - 旧库存价值
				oldStockValue := oldProduct.PurchasePrice.Mul(decimal.NewFromInt(int64(oldProduct.CurrentStock)))
				newStockValue := updateProduct.PurchasePrice.Mul(decimal.NewFromInt(int64(req.CurrentStock)))
				stockValueChange = newStockValue.Sub(oldStockValue)
			} else if priceChanged {
				// 只有价格变化：(新价格 - 旧价格) * 当前库存
				priceDiff := updateProduct.PurchasePrice.Sub(oldProduct.PurchasePrice)
				stockValueChange = priceDiff.Mul(decimal.NewFromInt(int64(oldProduct.CurrentStock)))
			} else {
				// 只有库存变化：库存变化量 * 采购价格
				stockValueChange = updateProduct.PurchasePrice.Mul(decimal.NewFromInt(int64(stockChange)))
			}

			// 检查低库存状态变化
			lowStockChange := 0
			oldLowStock := oldProduct.CurrentStock <= oldProduct.MinStock
			newLowStock := req.CurrentStock <= req.MinStock
			if oldLowStock != newLowStock {
				if newLowStock {
					lowStockChange = 1
				} else {
					lowStockChange = -1
				}
			}

			err = l.statisticsRepo.UpdateStockStats(l.ctx, tx, stockChange, stockValueChange, lowStockChange)
			if err != nil {
				l.Error("UpdateProduct UpdateStockStats err:", err.Error())
				return xerr.NewErrCodeMsg(xerr.ServerError, "更新库存统计失败")
			}
		}

		return nil
	})

	if err != nil {
		return false, err
	}

	return true, nil
}
