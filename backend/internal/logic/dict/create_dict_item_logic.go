// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package dict

import (
	"context"

	"admin_backend/internal/svc"
	"admin_backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateDictItemLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建字典数据
func NewCreateDictItemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateDictItemLogic {
	return &CreateDictItemLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateDictItemLogic) CreateDictItem(req *types.CreateDictItemReq) (resp *types.CreateDictItemResp, err error) {
	// todo: add your logic here and delete this line

	return
}
