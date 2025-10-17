// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package dict

import (
	"context"

	"admin_backend/internal/svc"
	"admin_backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateDictItemLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新字典数据
func NewUpdateDictItemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateDictItemLogic {
	return &UpdateDictItemLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateDictItemLogic) UpdateDictItem(req *types.UpdateDictItemReq) (resp bool, err error) {
	// todo: add your logic here and delete this line

	return
}
