// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package organization

import (
	"context"

	"admin_backend/internal/svc"
	"admin_backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListDepartmentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取部门列表
func NewListDepartmentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListDepartmentLogic {
	return &ListDepartmentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListDepartmentLogic) ListDepartment(req *types.DepartmentListReq) (resp *types.DepartmentListResp, err error) {
	// todo: add your logic here and delete this line

	return
}
