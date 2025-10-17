// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package organization

import (
	"context"

	"admin_backend/internal/svc"
	"admin_backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPositionUsersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取岗位下的用户列表
func NewGetPositionUsersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPositionUsersLogic {
	return &GetPositionUsersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPositionUsersLogic) GetPositionUsers(req *types.GetPositionUsersReq) (resp *types.GetPositionUsersResp, err error) {
	// todo: add your logic here and delete this line

	return
}
