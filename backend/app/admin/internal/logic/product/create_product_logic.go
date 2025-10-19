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
	"admin_backend/pkg/utils/idgen"

	"github.com/shopspring/decimal"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateProductLogic struct {
	logx.Logger
	ctx            context.Context
	svcCtx         *svc.ServiceContext
	productRepo    *productrepo.ProductRepo
	statisticsRepo *statisticsrepo.StatisticsRepo
}

// 创建商品
func NewCreateProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateProductLogic {
	return &CreateProductLogic{
		Logger:         logx.WithContext(ctx),
		ctx:            ctx,
		svcCtx:         svcCtx,
		productRepo:    productrepo.NewProductRepo(svcCtx.DB),
		statisticsRepo: statisticsrepo.NewStatisticsRepo(svcCtx.DB),
	}
}

func (l *CreateProductLogic) CreateProduct(req *types.CreateProductReq) (*types.CreateProductResp, error) {
	productID, err := idgen.GenerateUUID()
	if err != nil {
		return nil, xerr.NewErrCodeMsg(xerr.ServerError, "生成商品ID失败")
	}

	if req.Status == 0 {
		req.Status = 1 // 默认启用
	}

	// 解析价格
	purchasePrice, err := decimal.NewFromString(req.PurchasePrice)
	if err != nil {
		return nil, xerr.NewErrCodeMsg(xerr.ServerError, "采购价格格式错误")
	}

	salePrice, err := decimal.NewFromString(req.SalePrice)
	if err != nil {
		return nil, xerr.NewErrCodeMsg(xerr.ServerError, "销售价格格式错误")
	}

	// 创建商品（租户信息由 Repository 层自动处理）
	newProduct := &generated.Product{
		ProductID:     productID,
		ProductName:   req.ProductName,
		Unit:          req.Unit,
		PurchasePrice: purchasePrice,
		SalePrice:     salePrice,
		CurrentStock:  req.CurrentStock,
		MinStock:      req.MinStock,
		Status:        req.Status,
		FactoryID:     req.FactoryID,
	}
	var product *generated.Product
	err = ent.WithTx(l.ctx, l.svcCtx.DB, func(tx *generated.Tx) error {
		// 创建商品
		var err error
		product, err = l.productRepo.CreateWithTx(l.ctx, tx, newProduct)
		if err != nil {
			l.Error("CreateProduct Create err:", err.Error())
			return xerr.NewErrCodeMsg(xerr.ServerError, "创建商品失败")
		}

		// 更新统计数据
		err = l.statisticsRepo.IncrementProductCount(l.ctx, tx, req.Status == 1)
		if err != nil {
			l.Error("CreateProduct IncrementProductCount err:", err.Error())
			return xerr.NewErrCodeMsg(xerr.ServerError, "更新统计数据失败")
		}

		// 更新库存统计
		stockValue := purchasePrice.Mul(decimal.NewFromInt(int64(req.CurrentStock)))
		lowStockChange := 0
		if req.CurrentStock <= req.MinStock {
			lowStockChange = 1
		}
		err = l.statisticsRepo.UpdateStockStats(l.ctx, tx, req.CurrentStock, stockValue, lowStockChange)
		if err != nil {
			l.Error("CreateProduct UpdateStockStats err:", err.Error())
			return xerr.NewErrCodeMsg(xerr.ServerError, "更新库存统计失败")
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &types.CreateProductResp{
		ProductID: product.ProductID,
	}, nil
}
