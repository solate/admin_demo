package dict

import (
	"context"

	"admin_backend/app/admin/internal/repository/dictrepo"
	"admin_backend/app/admin/internal/svc"
	"admin_backend/app/admin/internal/types"
	"admin_backend/pkg/common/xerr"
	"admin_backend/pkg/ent/generated"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDictItemLogic struct {
	logx.Logger
	ctx      context.Context
	svcCtx   *svc.ServiceContext
	dictRepo *dictrepo.DictItemRepo
}

// 获取字典数据详情
func NewGetDictItemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDictItemLogic {
	return &GetDictItemLogic{
		Logger:   logx.WithContext(ctx),
		ctx:      ctx,
		svcCtx:   svcCtx,
		dictRepo: dictrepo.NewDictItemRepo(svcCtx.DB),
	}
}

func (l *GetDictItemLogic) GetDictItem(req *types.GetDictItemReq) (resp *types.DictItemInfo, err error) {
	// 1. 获取字典数据项
	dictItem, err := l.dictRepo.GetByItemID(l.ctx, req.ItemID)
	if err != nil {
		if generated.IsNotFound(err) {
			return nil, xerr.NewErrCodeMsg(xerr.DbError, "字典数据项不存在")
		}
		l.Error("GetDictItem l.dictRepo.GetByItemID err: ", err.Error())
		return nil, xerr.NewErrCodeMsg(xerr.DbError, "数据库查询错误")
	}

	// 2. 返回结果
	return &types.DictItemInfo{
		ItemID:      dictItem.ItemID,
		TypeCode:    dictItem.TypeCode,
		Label:       dictItem.Label,
		Value:       dictItem.Value,
		Sort:        dictItem.Sort,
		Description: dictItem.Description,
		Status:      dictItem.Status,
		CreatedAt:   dictItem.CreatedAt,
	}, nil
}
