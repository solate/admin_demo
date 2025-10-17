package organization

import (
	"context"

	"admin_backend/app/admin/internal/repository/organizationrepo"
	"admin_backend/app/admin/internal/svc"
	"admin_backend/app/admin/internal/types"
	"admin_backend/pkg/common/xerr"
	"admin_backend/pkg/ent/generated"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteDepartmentLogic struct {
	logx.Logger
	ctx     context.Context
	svcCtx  *svc.ServiceContext
	orgRepo *organizationrepo.DepartmentRepo
}

// 删除部门
func NewDeleteDepartmentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteDepartmentLogic {
	return &DeleteDepartmentLogic{
		Logger:  logx.WithContext(ctx),
		ctx:     ctx,
		svcCtx:  svcCtx,
		orgRepo: organizationrepo.NewDepartmentRepo(svcCtx.DB),
	}
}

func (l *DeleteDepartmentLogic) DeleteDepartment(req *types.DeleteDepartmentReq) (resp bool, err error) {
	// 1. 检查部门是否存在
	dept, err := l.orgRepo.GetByDepartmentID(l.ctx, req.DepartmentID)
	if err != nil {
		if generated.IsNotFound(err) {
			return false, xerr.NewErrCodeMsg(xerr.RecordNotFound, "部门不存在")
		}
		l.Error("GetDepartment l.orgRepo.GetByDepartmentID err: ", err.Error())
		return false, xerr.NewErrCodeMsg(xerr.DbError, "数据库查询错误")
	}

	// 2. 执行删除
	_, err = l.orgRepo.Delete(l.ctx, dept.DepartmentID)
	if err != nil {
		l.Error("DeleteDepartment Delete err:", err.Error())
		return false, xerr.NewErrCodeMsg(xerr.ServerError, "删除部门失败")
	}

	return true, nil
}
