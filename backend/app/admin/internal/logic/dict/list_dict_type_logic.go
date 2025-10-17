package dict

import (
	"context"

	"admin_backend/app/admin/internal/repository/dictrepo"
	"admin_backend/app/admin/internal/svc"
	"admin_backend/app/admin/internal/types"
	"admin_backend/pkg/common/xerr"
	"admin_backend/pkg/ent/generated/dicttype"
	"admin_backend/pkg/ent/generated/predicate"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListDictTypeLogic struct {
	logx.Logger
	ctx      context.Context
	svcCtx   *svc.ServiceContext
	dictRepo *dictrepo.DictTypeRepo
}

// 获取字典类型列表
func NewListDictTypeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListDictTypeLogic {
	return &ListDictTypeLogic{
		Logger:   logx.WithContext(ctx),
		ctx:      ctx,
		svcCtx:   svcCtx,
		dictRepo: dictrepo.NewDictTypeRepo(svcCtx.DB),
	}
}

func (l *ListDictTypeLogic) ListDictType(req *types.DictTypeListReq) (resp *types.DictTypeListResp, err error) {
	// 1. 构建查询条件
	var predicates []predicate.DictType

	if req.Name != "" {
		predicates = append(predicates, dicttype.NameContains(req.Name))
	}
	if req.Code != "" {
		predicates = append(predicates, dicttype.CodeContains(req.Code))
	}
	if req.Status != 0 {
		predicates = append(predicates, dicttype.Status(req.Status))
	}

	// 2. 查询数据
	list, total, err := l.dictRepo.PageList(l.ctx, req.Current, req.PageSize, predicates)
	if err != nil {
		l.Error("ListDictType l.dictRepo.PageList err:", err.Error())
		return nil, xerr.NewErrCodeMsg(xerr.DbError, "查询字典类型列表失败")
	}

	// 3. 构建返回结果
	var dictTypes []*types.DictTypeInfo
	for _, item := range list {
		dictTypes = append(dictTypes, &types.DictTypeInfo{
			TypeID:      item.TypeID,
			Name:        item.Name,
			Code:        item.Code,
			Description: item.Description,
			Status:      item.Status,
			CreatedAt:   item.CreatedAt,
		})
	}

	return &types.DictTypeListResp{
		Page: &types.PageResponse{
			Total:           total,
			PageSize:        len(list),
			Current:         req.Current,
			RequestPageSize: req.PageSize,
		},
		List: dictTypes,
	}, nil
}
