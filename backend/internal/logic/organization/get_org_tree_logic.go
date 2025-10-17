// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package organization

import (
	"context"

	"admin_backend/internal/svc"
	"admin_backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOrgTreeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取组织架构树
func NewGetOrgTreeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOrgTreeLogic {
	return &GetOrgTreeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetOrgTreeLogic) GetOrgTree(req *types.GetOrgTreeReq) (resp *types.GetOrgTreeResp, err error) {
	// todo: add your logic here and delete this line

	return
}
