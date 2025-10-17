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

type DeleteProductLogic struct {
	logx.Logger
	ctx            context.Context
	svcCtx         *svc.ServiceContext
	productRepo    *productrepo.ProductRepo
	statisticsRepo *statisticsrepo.StatisticsRepo
}

// 删除商品
func NewDeleteProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteProductLogic {
	return &DeleteProductLogic{
		Logger:         logx.WithContext(ctx),
		ctx:            ctx,
		svcCtx:         svcCtx,
		productRepo:    productrepo.NewProductRepo(svcCtx.DB),
		statisticsRepo: statisticsrepo.NewStatisticsRepo(svcCtx.DB),
	}
}

func (l *DeleteProductLogic) DeleteProduct(req *types.DeleteProductReq) (bool, error) {
	// 检查商品是否存在，获取商品信息
	product, err := l.productRepo.GetByProductID(l.ctx, req.ProductID)
	if err != nil {
		l.Error("DeleteProduct GetByProductID err:", err.Error())
		return false, xerr.NewErrCodeMsg(xerr.ServerError, "商品不存在")
	}

	// 使用事务删除商品和更新统计
	err = ent.WithTx(l.ctx, l.svcCtx.DB, func(tx *generated.Tx) error {
		// 软删除商品
		_, err := l.productRepo.DeleteWithTx(l.ctx, tx, req.ProductID)
		if err != nil {
			l.Error("DeleteProduct Delete err:", err.Error())
			return xerr.NewErrCodeMsg(xerr.ServerError, "删除商品失败")
		}

		// 更新统计数据
		err = l.statisticsRepo.DecrementProductCount(l.ctx, tx, product.Status == 1)
		if err != nil {
			l.Error("DeleteProduct DecrementProductCount err:", err.Error())
			return xerr.NewErrCodeMsg(xerr.ServerError, "更新商品统计失败")
		}

		// 更新库存统计（减少库存）
		stockValueChange := product.PurchasePrice.Mul(decimal.NewFromInt(int64(-product.CurrentStock)))
		lowStockChange := 0
		if product.CurrentStock <= product.MinStock {
			lowStockChange = -1
		}

		err = l.statisticsRepo.UpdateStockStats(l.ctx, tx, -product.CurrentStock, stockValueChange, lowStockChange)
		if err != nil {
			l.Error("DeleteProduct UpdateStockStats err:", err.Error())
			return xerr.NewErrCodeMsg(xerr.ServerError, "更新库存统计失败")
		}

		return nil
	})

	if err != nil {
		return false, err
	}

	return true, nil
}
