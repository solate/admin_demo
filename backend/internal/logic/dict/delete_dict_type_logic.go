// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package dict

import (
	"context"

	"admin_backend/internal/svc"
	"admin_backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteDictTypeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除字典类型
func NewDeleteDictTypeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteDictTypeLogic {
	return &DeleteDictTypeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteDictTypeLogic) DeleteDictType(req *types.DeleteDictTypeReq) (resp bool, err error) {
	// todo: add your logic here and delete this line

	return
}
