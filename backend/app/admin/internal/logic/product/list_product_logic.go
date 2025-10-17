package product

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

type ListProductLogic struct {
	logx.Logger
	ctx         context.Context
	svcCtx      *svc.ServiceContext
	productRepo *productrepo.ProductRepo
}

// 获取商品列表
func NewListProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListProductLogic {
	return &ListProductLogic{
		Logger:      logx.WithContext(ctx),
		ctx:         ctx,
		svcCtx:      svcCtx,
		productRepo: productrepo.NewProductRepo(svcCtx.DB),
	}
}

func (l *ListProductLogic) ListProduct(req *types.ProductListReq) (*types.ProductListResp, error) {
	// 构建查询条件
	var where []predicate.Product

	if req.ProductName != "" {
		where = append(where, product.ProductNameContains(req.ProductName))
	}
	if req.FactoryID != "" {
		where = append(where, product.FactoryID(req.FactoryID))
	}
	if req.Status != 0 {
		where = append(where, product.Status(req.Status))
	}

	// 分页查询
	list, total, err := l.productRepo.PageList(l.ctx, req.Current, req.PageSize, where)
	if err != nil {
		l.Error("ListProduct PageList err:", err.Error())
		return nil, xerr.NewErrCode(xerr.DbError)
	}

	// 构建响应数据
	productList := make([]*types.ProductInfo, 0)
	for _, item := range list {
		// 获取工厂名称
		factoryName := ""
		if item.FactoryID != "" {
			// 这里应该查询工厂信息，暂时使用空字符串
		}

		productList = append(productList, &types.ProductInfo{
			ProductID:     item.ProductID,
			ProductName:   item.ProductName,
			Unit:          item.Unit,
			PurchasePrice: item.PurchasePrice.String(),
			SalePrice:     item.SalePrice.String(),
			CurrentStock:  item.CurrentStock,
			MinStock:      item.MinStock,
			Status:        item.Status,
			FactoryID:     item.FactoryID,
			FactoryName:   factoryName,
			CreatedAt:     item.CreatedAt,
			UpdatedAt:     item.UpdatedAt,
		})
	}

	return &types.ProductListResp{
		List: productList,
		Page: &types.PageResponse{
			Total:           total,
			PageSize:        len(list),
			Current:         req.Current,
			RequestPageSize: req.PageSize,
		},
	}, nil
}
