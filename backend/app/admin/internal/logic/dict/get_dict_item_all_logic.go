package dict

import (
	"context"

	"admin_backend/app/admin/internal/repository/dictrepo"
	"admin_backend/app/admin/internal/svc"
	"admin_backend/app/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDictItemAllLogic struct {
	logx.Logger
	ctx      context.Context
	svcCtx   *svc.ServiceContext
	dictRepo *dictrepo.DictItemRepo
}

// 获取字典项选项列表
func NewGetDictItemAllLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDictItemAllLogic {
	return &GetDictItemAllLogic{
		Logger:   logx.WithContext(ctx),
		ctx:      ctx,
		svcCtx:   svcCtx,
		dictRepo: dictrepo.NewDictItemRepo(svcCtx.DB),
	}
}

func (l *GetDictItemAllLogic) GetDictItemAll(req *types.GetDictItemAllReq) (resp *types.DictItemAllResp, err error) {
	// 根据字典类型编码获取所有字典项
	list, err := l.dictRepo.GetByTypeCode(l.ctx, req.TypeCode)
	if err != nil {
		return nil, err
	}

	// 转换为前端需要的格式
	resp = &types.DictItemAllResp{
		List: make([]*types.DictItemInfo, 0, len(list)),
	}

	for _, item := range list {
		resp.List = append(resp.List, &types.DictItemInfo{
			TypeCode:    item.TypeCode,
			Label:       item.Label,
			Value:       item.Value,
			Description: item.Description,
			Sort:        item.Sort,
			Status:      item.Status,
			CreatedAt:   item.CreatedAt,
		})
	}

	return resp, nil
}
