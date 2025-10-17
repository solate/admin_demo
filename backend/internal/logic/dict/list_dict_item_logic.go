// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package dict

import (
	"context"

	"admin_backend/internal/svc"
	"admin_backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListDictItemLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取字典数据列表
func NewListDictItemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListDictItemLogic {
	return &ListDictItemLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListDictItemLogic) ListDictItem(req *types.DictItemListReq) (resp *types.DictItemListResp, err error) {
	// todo: add your logic here and delete this line

	return
}
