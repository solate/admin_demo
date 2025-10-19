package inventory

import (
	"context"
	"time"

	"admin_backend/app/admin/internal/repository/inventoryrepo"
	"admin_backend/app/admin/internal/repository/sysuserrepo"
	"admin_backend/app/admin/internal/svc"
	"admin_backend/app/admin/internal/types"
	"admin_backend/pkg/common/xerr"
	"admin_backend/pkg/ent/generated/inventory"
	"admin_backend/pkg/ent/generated/predicate"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListInventoryLogic struct {
	logx.Logger
	ctx           context.Context
	svcCtx        *svc.ServiceContext
	inventoryRepo *inventoryrepo.InventoryRepo
	userRepo      *sysuserrepo.SysUserRepo
}

// 获取库存列表
func NewListInventoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListInventoryLogic {
	return &ListInventoryLogic{
		Logger:        logx.WithContext(ctx),
		ctx:           ctx,
		svcCtx:        svcCtx,
		inventoryRepo: inventoryrepo.NewInventoryRepo(svcCtx.DB),
		userRepo:      sysuserrepo.NewSysUserRepo(svcCtx.DB),
	}
}

func (l *ListInventoryLogic) ListInventory(req *types.InventoryListReq) (*types.InventoryListResp, error) {
	// 构建查询条件
	var where []predicate.Inventory

	if req.ProductID != "" {
		where = append(where, inventory.ProductID(req.ProductID))
	}
	if req.OperationType != "" {
		where = append(where, inventory.OperationType(req.OperationType))
	}
	if req.OperatorID != "" {
		where = append(where, inventory.OperatorID(req.OperatorID))
	}

	// 时间范围查询
	if req.StartTime != "" && req.EndTime != "" {
		startTime, err := time.Parse("2006-01-02 15:04:05", req.StartTime)
		if err == nil {
			endTime, err := time.Parse("2006-01-02 15:04:05", req.EndTime)
			if err == nil {
				where = append(where, inventory.OperationTimeGTE(startTime.UnixMilli()))
				where = append(where, inventory.OperationTimeLTE(endTime.UnixMilli()))
			}
		}
	}

	// 分页查询
	list, total, err := l.inventoryRepo.PageList(l.ctx, req.Current, req.PageSize, where)
	if err != nil {
		l.Error("ListInventory PageList err:", err.Error())
		return nil, xerr.NewErrCode(xerr.DbError)
	}

	// 构建响应数据
	inventoryList := make([]*types.InventoryInfo, 0)
	for _, item := range list {
		// 获取操作人名称
		operatorName := ""
		if item.OperatorID != "" {
			user, err := l.userRepo.GetByUserID(l.ctx, item.OperatorID)
			if err == nil {
				operatorName = user.UserName
			}
		}

		inventoryList = append(inventoryList, &types.InventoryInfo{
			InventoryID:   item.InventoryID,
			ProductID:     item.ProductID,
			ProductName:   "",
			OperationType: item.OperationType,
			Quantity:      item.Quantity,
			UnitPrice:     item.UnitPrice.String(),
			TotalAmount:   item.TotalAmount.String(),
			OperatorID:    item.OperatorID,
			OperatorName:  operatorName,
			Remark:        item.Remark,
			OperationTime: item.OperationTime,
			BeforeStock:   item.BeforeStock,
			AfterStock:    item.AfterStock,
		})
	}

	return &types.InventoryListResp{
		List: inventoryList,
		Page: &types.PageResponse{
			Total:           total,
			PageSize:        len(list),
			Current:         req.Current,
			RequestPageSize: req.PageSize,
		},
	}, nil
}
