package factory

import (
	"context"

	"admin_backend/app/admin/internal/repository/factoryrepo"
	"admin_backend/app/admin/internal/svc"
	"admin_backend/app/admin/internal/types"
	"admin_backend/pkg/common/xerr"
	"admin_backend/pkg/ent/generated"
	"admin_backend/pkg/utils/idgen"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateFactoryLogic struct {
	logx.Logger
	ctx         context.Context
	svcCtx      *svc.ServiceContext
	factoryRepo *factoryrepo.FactoryRepo
}

// 创建工厂
func NewCreateFactoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateFactoryLogic {
	return &CreateFactoryLogic{
		Logger:      logx.WithContext(ctx),
		ctx:         ctx,
		svcCtx:      svcCtx,
		factoryRepo: factoryrepo.NewFactoryRepo(svcCtx.DB),
	}
}

func (l *CreateFactoryLogic) CreateFactory(req *types.CreateFactoryReq) (*types.CreateFactoryResp, error) {
	factoryID, err := idgen.GenerateUUID()
	if err != nil {
		return nil, xerr.NewErrCodeMsg(xerr.ServerError, "生成工厂ID失败")
	}

	if req.Status == 0 {
		req.Status = 1 // 默认启用
	}

	// 创建工厂（租户信息由 Repository 层自动处理）
	newFactory := &generated.Factory{
		FactoryID:    factoryID,
		FactoryName:  req.FactoryName,
		Address:      req.Address,
		ContactPhone: req.ContactPhone,
		Status:       req.Status,
	}
	factory, err := l.factoryRepo.Create(l.ctx, newFactory)
	if err != nil {
		l.Error("CreateFactory Create err:", err.Error())
		return nil, xerr.NewErrCodeMsg(xerr.ServerError, "创建工厂失败")
	}

	return &types.CreateFactoryResp{
		FactoryID: factory.FactoryID,
	}, nil
}
