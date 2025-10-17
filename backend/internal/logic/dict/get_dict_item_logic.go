// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package dict

import (
	"context"

	"admin_backend/internal/svc"
	"admin_backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDictItemLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取字典数据详情
func NewGetDictItemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDictItemLogic {
	return &GetDictItemLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDictItemLogic) GetDictItem(req *types.GetDictItemReq) (resp *types.DictItemInfo, err error) {
	// todo: add your logic here and delete this line

	return
}
