package inventory

import (
	"context"

	"admin_backend/app/admin/internal/repository/productrepo"
	"admin_backend/app/admin/internal/svc"
	"admin_backend/app/admin/internal/types"
	"admin_backend/pkg/common/xerr"
	"admin_backend/pkg/ent/generated/predicate"
	"admin_backend/pkg/ent/generated/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetProductStockLogic struct {
	logx.Logger
	ctx         context.Context
	svcCtx      *svc.ServiceContext
	productRepo *productrepo.ProductRepo
}

// 获取商品库存信息
func NewGetProductStockLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProductStockLogic {
	return &GetProductStockLogic{
		Logger:      logx.WithContext(ctx),
		ctx:         ctx,
		svcCtx:      svcCtx,
		productRepo: productrepo.NewProductRepo(svcCtx.DB),
	}
}

func (l *GetProductStockLogic) GetProductStock(req *types.ProductStockReq) (*types.ProductStockResp, error) {
	// 构建查询条件
	var where []predicate.Product

	if req.ProductID != "" {
		where = append(where, product.ProductID(req.ProductID))
	}
	if req.FactoryID != "" {
		where = append(where, product.FactoryID(req.FactoryID))
	}

	// 查询商品列表
	products, err := l.productRepo.GetAll(l.ctx, where)
	if err != nil {
		l.Error("GetProductStock GetAll err:", err.Error())
		return nil, xerr.NewErrCode(xerr.DbError)
	}

	// 构建响应数据
	var stockList []*types.StockInfo
	for _, item := range products {
		// 检查是否低库存
		isLowStock := item.CurrentStock <= item.MinStock

		// 获取工厂名称
		factoryName := "" // 这里应该查询工厂信息

		stockList = append(stockList, &types.StockInfo{
			ProductID:     item.ProductID,
			ProductName:   item.ProductName,
			FactoryID:     item.FactoryID,
			FactoryName:   factoryName,
			CurrentStock:  item.CurrentStock,
			MinStock:      item.MinStock,
			Unit:          item.Unit,
			PurchasePrice: item.PurchasePrice.String(),
			SalePrice:     item.SalePrice.String(),
			Status:        item.Status,
			IsLowStock:    isLowStock,
		})
	}

	return &types.ProductStockResp{
		List: stockList,
	}, nil
}
