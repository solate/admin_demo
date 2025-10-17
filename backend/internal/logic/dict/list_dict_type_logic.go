// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package dict

import (
	"context"

	"admin_backend/internal/svc"
	"admin_backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListDictTypeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取字典类型列表
func NewListDictTypeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListDictTypeLogic {
	return &ListDictTypeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListDictTypeLogic) ListDictType(req *types.DictTypeListReq) (resp *types.DictTypeListResp, err error) {
	// todo: add your logic here and delete this line

	return
}
