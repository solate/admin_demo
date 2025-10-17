package factory

import (
	"context"

	"admin_backend/app/admin/internal/repository/factoryrepo"
	"admin_backend/app/admin/internal/svc"
	"admin_backend/app/admin/internal/types"
	"admin_backend/pkg/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFactoryLogic struct {
	logx.Logger
	ctx         context.Context
	svcCtx      *svc.ServiceContext
	factoryRepo *factoryrepo.FactoryRepo
}

// 获取工厂详情
func NewGetFactoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFactoryLogic {
	return &GetFactoryLogic{
		Logger:      logx.WithContext(ctx),
		ctx:         ctx,
		svcCtx:      svcCtx,
		factoryRepo: factoryrepo.NewFactoryRepo(svcCtx.DB),
	}
}

func (l *GetFactoryLogic) GetFactory(req *types.GetFactoryReq) (*types.FactoryInfo, error) {
	// 获取工厂信息
	factory, err := l.factoryRepo.GetByFactoryID(l.ctx, req.FactoryID)
	if err != nil {
		l.Error("GetFactory GetByFactoryID err:", err.Error())
		return nil, xerr.NewErrCodeMsg(xerr.ServerError, "工厂不存在")
	}

	return &types.FactoryInfo{
		FactoryID:    factory.FactoryID,
		FactoryName:  factory.FactoryName,
		Address:      factory.Address,
		ContactPhone: factory.ContactPhone,
		Status:       factory.Status,
		CreatedAt:    factory.CreatedAt,
		UpdatedAt:    factory.UpdatedAt,
	}, nil
}
