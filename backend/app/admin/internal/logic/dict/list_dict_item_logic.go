package dict

import (
	"context"

	"admin_backend/app/admin/internal/repository/dictrepo"
	"admin_backend/app/admin/internal/svc"
	"admin_backend/app/admin/internal/types"
	"admin_backend/pkg/common/xerr"
	"admin_backend/pkg/ent/generated/dictitem"
	"admin_backend/pkg/ent/generated/predicate"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListDictItemLogic struct {
	logx.Logger
	ctx      context.Context
	svcCtx   *svc.ServiceContext
	dictRepo *dictrepo.DictItemRepo
}

// 获取字典数据列表
func NewListDictItemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListDictItemLogic {
	return &ListDictItemLogic{
		Logger:   logx.WithContext(ctx),
		ctx:      ctx,
		svcCtx:   svcCtx,
		dictRepo: dictrepo.NewDictItemRepo(svcCtx.DB),
	}
}

func (l *ListDictItemLogic) ListDictItem(req *types.DictItemListReq) (resp *types.DictItemListResp, err error) {
	// 1. 构建查询条件
	var predicates []predicate.DictItem

	if req.Label != "" {
		predicates = append(predicates, dictitem.LabelContains(req.Label))
	}
	if req.TypeCode != "" {
		predicates = append(predicates, dictitem.TypeCode(req.TypeCode))
	}
	if req.Status != 0 {
		predicates = append(predicates, dictitem.Status(req.Status))
	}

	// 2. 查询数据
	list, total, err := l.dictRepo.PageList(l.ctx, req.Current, req.PageSize, predicates)
	if err != nil {
		l.Error("ListDictItem l.dictRepo.PageList err:", err.Error())
		return nil, xerr.NewErrCodeMsg(xerr.DbError, "查询字典数据列表失败")
	}

	// 3. 构建返回结果
	var dictItems []*types.DictItemInfo
	for _, item := range list {
		dictItems = append(dictItems, &types.DictItemInfo{
			ItemID:      item.ItemID,
			TypeCode:    item.TypeCode,
			Label:       item.Label,
			Value:       item.Value,
			Sort:        item.Sort,
			Description: item.Description,
			Status:      item.Status,
			CreatedAt:   item.CreatedAt,
		})
	}

	return &types.DictItemListResp{
		Page: &types.PageResponse{
			Total:           total,
			PageSize:        len(list),
			Current:         req.Current,
			RequestPageSize: req.PageSize,
		},
		List: dictItems,
	}, nil
}
