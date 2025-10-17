package factory

import (
	"context"

	"admin_backend/app/admin/internal/repository/factoryrepo"
	"admin_backend/app/admin/internal/svc"
	"admin_backend/app/admin/internal/types"
	"admin_backend/pkg/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteFactoryLogic struct {
	logx.Logger
	ctx         context.Context
	svcCtx      *svc.ServiceContext
	factoryRepo *factoryrepo.FactoryRepo
}

// 删除工厂
func NewDeleteFactoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteFactoryLogic {
	return &DeleteFactoryLogic{
		Logger:      logx.WithContext(ctx),
		ctx:         ctx,
		svcCtx:      svcCtx,
		factoryRepo: factoryrepo.NewFactoryRepo(svcCtx.DB),
	}
}

func (l *DeleteFactoryLogic) DeleteFactory(req *types.DeleteFactoryReq) (bool, error) {
	// 检查工厂是否存在
	_, err := l.factoryRepo.GetByFactoryID(l.ctx, req.FactoryID)
	if err != nil {
		l.Error("DeleteFactory GetByFactoryID err:", err.Error())
		return false, xerr.NewErrCodeMsg(xerr.ServerError, "工厂不存在")
	}

	// 软删除工厂
	_, err = l.factoryRepo.Delete(l.ctx, req.FactoryID)
	if err != nil {
		l.Error("DeleteFactory Delete err:", err.Error())
		return false, xerr.NewErrCodeMsg(xerr.ServerError, "删除工厂失败")
	}

	return true, nil
}
