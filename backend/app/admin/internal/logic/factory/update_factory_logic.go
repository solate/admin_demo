package factory

import (
	"context"

	"admin_backend/app/admin/internal/repository/factoryrepo"
	"admin_backend/app/admin/internal/svc"
	"admin_backend/app/admin/internal/types"
	"admin_backend/pkg/common/xerr"
	"admin_backend/pkg/ent/generated"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateFactoryLogic struct {
	logx.Logger
	ctx         context.Context
	svcCtx      *svc.ServiceContext
	factoryRepo *factoryrepo.FactoryRepo
}

// 更新工厂
func NewUpdateFactoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateFactoryLogic {
	return &UpdateFactoryLogic{
		Logger:      logx.WithContext(ctx),
		ctx:         ctx,
		svcCtx:      svcCtx,
		factoryRepo: factoryrepo.NewFactoryRepo(svcCtx.DB),
	}
}

func (l *UpdateFactoryLogic) UpdateFactory(req *types.UpdateFactoryReq) (bool, error) {
	// 检查工厂是否存在
	_, err := l.factoryRepo.GetByFactoryID(l.ctx, req.FactoryID)
	if err != nil {
		l.Error("UpdateFactory GetByFactoryID err:", err.Error())
		return false, xerr.NewErrCodeMsg(xerr.ServerError, "工厂不存在")
	}

	// 构建更新数据（租户信息由 Repository 层自动处理）
	updateFactory := &generated.Factory{
		FactoryID:    req.FactoryID,
		FactoryName:  req.FactoryName,
		Address:      req.Address,
		ContactPhone: req.ContactPhone,
		Status:       req.Status,
	}

	// 更新工厂
	_, err = l.factoryRepo.Update(l.ctx, updateFactory)
	if err != nil {
		l.Error("UpdateFactory Update err:", err.Error())
		return false, xerr.NewErrCodeMsg(xerr.ServerError, "更新工厂失败")
	}

	return true, nil
}
