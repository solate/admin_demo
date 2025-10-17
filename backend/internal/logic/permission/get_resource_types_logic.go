// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package permission

import (
	"context"

	"admin_backend/internal/svc"
	"admin_backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetResourceTypesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取资源类型列表
func NewGetResourceTypesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetResourceTypesLogic {
	return &GetResourceTypesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetResourceTypesLogic) GetResourceTypes() (resp *types.GetResourceTypesResp, err error) {
	// todo: add your logic here and delete this line

	return
}
