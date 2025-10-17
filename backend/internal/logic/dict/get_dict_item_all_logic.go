// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package dict

import (
	"context"

	"admin_backend/internal/svc"
	"admin_backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDictItemAllLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取字典项选项列表
func NewGetDictItemAllLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDictItemAllLogic {
	return &GetDictItemAllLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDictItemAllLogic) GetDictItemAll(req *types.GetDictItemAllReq) (resp *types.DictItemAllResp, err error) {
	// todo: add your logic here and delete this line

	return
}
