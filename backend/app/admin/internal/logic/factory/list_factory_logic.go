package factory

import (
	"context"

	"admin_backend/app/admin/internal/repository/factoryrepo"
	"admin_backend/app/admin/internal/svc"
	"admin_backend/app/admin/internal/types"
	"admin_backend/pkg/common/xerr"
	"admin_backend/pkg/ent/generated/factory"
	"admin_backend/pkg/ent/generated/predicate"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListFactoryLogic struct {
	logx.Logger
	ctx         context.Context
	svcCtx      *svc.ServiceContext
	factoryRepo *factoryrepo.FactoryRepo
}

// 获取工厂列表
func NewListFactoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListFactoryLogic {
	return &ListFactoryLogic{
		Logger:      logx.WithContext(ctx),
		ctx:         ctx,
		svcCtx:      svcCtx,
		factoryRepo: factoryrepo.NewFactoryRepo(svcCtx.DB),
	}
}

func (l *ListFactoryLogic) ListFactory(req *types.FactoryListReq) (*types.FactoryListResp, error) {
	// 构建查询条件
	var where []predicate.Factory

	if req.FactoryName != "" {
		where = append(where, factory.FactoryNameContains(req.FactoryName))
	}
	if req.Status != 0 {
		where = append(where, factory.Status(req.Status))
	}

	// 分页查询
	list, total, err := l.factoryRepo.PageList(l.ctx, req.Current, req.PageSize, where)
	if err != nil {
		l.Error("ListFactory PageList err:", err.Error())
		return nil, xerr.NewErrCode(xerr.DbError)
	}

	// 构建响应数据
	factoryList := make([]*types.FactoryInfo, 0)
	for _, item := range list {
		factoryList = append(factoryList, &types.FactoryInfo{
			FactoryID:    item.FactoryID,
			FactoryName:  item.FactoryName,
			Address:      item.Address,
			ContactPhone: item.ContactPhone,
			Status:       item.Status,
			CreatedAt:    item.CreatedAt,
			UpdatedAt:    item.UpdatedAt,
		})
	}

	return &types.FactoryListResp{
		List: factoryList,
		Page: &types.PageResponse{
			Total:           total,
			PageSize:        len(list),
			Current:         req.Current,
			RequestPageSize: req.PageSize,
		},
	}, nil
}
