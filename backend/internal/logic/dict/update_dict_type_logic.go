// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package dict

import (
	"context"

	"admin_backend/internal/svc"
	"admin_backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateDictTypeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新字典类型
func NewUpdateDictTypeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateDictTypeLogic {
	return &UpdateDictTypeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateDictTypeLogic) UpdateDictType(req *types.UpdateDictTypeReq) (resp bool, err error) {
	// todo: add your logic here and delete this line

	return
}
