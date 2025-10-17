// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package dict

import (
	"context"

	"admin_backend/internal/svc"
	"admin_backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateDictTypeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建字典类型
func NewCreateDictTypeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateDictTypeLogic {
	return &CreateDictTypeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateDictTypeLogic) CreateDictType(req *types.CreateDictTypeReq) (resp *types.CreateDictTypeResp, err error) {
	// todo: add your logic here and delete this line

	return
}
