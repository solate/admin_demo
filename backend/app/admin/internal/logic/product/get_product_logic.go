package product

import (
	"context"

	"admin_backend/app/admin/internal/repository/productrepo"
	"admin_backend/app/admin/internal/svc"
	"admin_backend/app/admin/internal/types"
	"admin_backend/pkg/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetProductLogic struct {
	logx.Logger
	ctx         context.Context
	svcCtx      *svc.ServiceContext
	productRepo *productrepo.ProductRepo
}

// 获取商品详情
func NewGetProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProductLogic {
	return &GetProductLogic{
		Logger:      logx.WithContext(ctx),
		ctx:         ctx,
		svcCtx:      svcCtx,
		productRepo: productrepo.NewProductRepo(svcCtx.DB),
	}
}

func (l *GetProductLogic) GetProduct(req *types.GetProductReq) (*types.ProductInfo, error) {
	// 获取商品信息
	product, err := l.productRepo.GetByProductID(l.ctx, req.ProductID)
	if err != nil {
		l.Error("GetProduct GetByProductID err:", err.Error())
		return nil, xerr.NewErrCodeMsg(xerr.ServerError, "商品不存在")
	}

	// 获取工厂名称
	factoryName := ""
	if product.FactoryID != "" {
		// 这里应该查询工厂信息，暂时使用空字符串
		// factory, err := l.factoryRepo.GetByFactoryID(l.ctx, product.FactoryID)
		// if err == nil {
		//     factoryName = factory.FactoryName
		// }
	}

	return &types.ProductInfo{
		ProductID:     product.ProductID,
		ProductName:   product.ProductName,
		Unit:          product.Unit,
		PurchasePrice: product.PurchasePrice.String(),
		SalePrice:     product.SalePrice.String(),
		CurrentStock:  product.CurrentStock,
		MinStock:      product.MinStock,
		Status:        product.Status,
		FactoryID:     product.FactoryID,
		FactoryName:   factoryName,
		CreatedAt:     product.CreatedAt,
		UpdatedAt:     product.UpdatedAt,
	}, nil
}