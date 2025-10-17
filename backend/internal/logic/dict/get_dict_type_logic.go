// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package dict

import (
	"context"

	"admin_backend/internal/svc"
	"admin_backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDictTypeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取字典类型详情
func NewGetDictTypeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDictTypeLogic {
	return &GetDictTypeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDictTypeLogic) GetDictType(req *types.GetDictTypeReq) (resp *types.DictTypeInfo, err error) {
	// todo: add your logic here and delete this line

	return
}
