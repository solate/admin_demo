package statistics

import (
	"context"

	"admin_backend/app/admin/internal/repository/statisticsrepo"
	"admin_backend/app/admin/internal/svc"
	"admin_backend/app/admin/internal/types"
	"admin_backend/pkg/common/xerr"

	"github.com/shopspring/decimal"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetStatisticsLogic struct {
	logx.Logger
	ctx            context.Context
	svcCtx         *svc.ServiceContext
	statisticsRepo *statisticsrepo.StatisticsRepo
}

func NewGetStatisticsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetStatisticsLogic {
	return &GetStatisticsLogic{
		Logger:         logx.WithContext(ctx),
		ctx:            ctx,
		svcCtx:         svcCtx,
		statisticsRepo: statisticsrepo.NewStatisticsRepo(svcCtx.DB),
	}
}

func (l *GetStatisticsLogic) GetStatistics(req *types.StatisticsReq) (*types.StatisticsResp, error) {
	// 使用repository计算实时统计
	stats, err := l.statisticsRepo.CalculateProductStatistics(l.ctx)
	if err != nil {
		l.Error("GetStatistics CalculateProductStatistics err:", err.Error())
		return nil, xerr.NewErrCodeMsg(xerr.ServerError, "获取统计数据失败")
	}

	return &types.StatisticsResp{
		TotalProducts:      stats["total_products"].(int),
		ActiveProducts:     stats["active_products"].(int),
		TotalStock:         stats["total_stock"].(int),
		TotalStockValue:    stats["total_stock_value"].(decimal.Decimal).String(),
		LowStockProducts:   stats["low_stock_products"].(int),
		TotalInQuantity:    stats["total_in_quantity"].(int),
		TotalInAmount:      stats["total_in_amount"].(decimal.Decimal).String(),
		TotalOutQuantity:   stats["total_out_quantity"].(int),
		TotalOutAmount:     stats["total_out_amount"].(decimal.Decimal).String(),
		TotalSalesAmount:   stats["total_sales_amount"].(decimal.Decimal).String(),
		TotalSalesQuantity: stats["total_sales_quantity"].(int),
	}, nil
}
