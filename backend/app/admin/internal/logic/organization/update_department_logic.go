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

type UpdateDepartmentLogic struct {
	logx.Logger
	ctx     context.Context
	svcCtx  *svc.ServiceContext
	orgRepo *organizationrepo.DepartmentRepo
}

// 更新部门
func NewUpdateDepartmentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateDepartmentLogic {
	return &UpdateDepartmentLogic{
		Logger:  logx.WithContext(ctx),
		ctx:     ctx,
		svcCtx:  svcCtx,
		orgRepo: organizationrepo.NewDepartmentRepo(svcCtx.DB),
	}
}

func (l *UpdateDepartmentLogic) UpdateDepartment(req *types.UpdateDepartmentReq) (resp bool, err error) {
	// 1. 检查部门是否存在
	dept, err := l.orgRepo.GetByDepartmentID(l.ctx, req.DepartmentID)
	if err != nil {
		if generated.IsNotFound(err) {
			return false, xerr.NewErrCodeMsg(xerr.RecordNotFound, "部门不存在")
		}
		l.Error("GetDepartment l.orgRepo.GetByDepartmentID err: ", err.Error())
		return false, xerr.NewErrCodeMsg(xerr.DbError, "数据库查询错误")
	}

	// 2. 如果要更新部门名称，检查新名称是否已存在
	if req.Name != "" && req.Name != dept.Name {
		existDept, err := l.orgRepo.GetByName(l.ctx, req.Name)
		if err != nil && !generated.IsNotFound(err) {
			l.Error("GetDepartment l.orgRepo.GetByName err: ", err.Error())
			return false, xerr.NewErrCodeMsg(xerr.DbError, "数据库查询错误")
		}
		if existDept != nil {
			return false, xerr.NewErrCodeMsg(xerr.DbRecordExist, "部门名称已存在")
		}
		dept.Name = req.Name
	}

	// 3. 更新其他字段
	if req.ParentID != "" {
		dept.ParentID = req.ParentID
	}
	if req.Sort != 0 {
		dept.Sort = req.Sort
	}

	// 4. 执行更新
	_, err = l.orgRepo.Update(l.ctx, dept)
	if err != nil {
		l.Error("UpdateDepartment Update err:", err.Error())
		return false, xerr.NewErrCodeMsg(xerr.ServerError, "更新部门失败")
	}

	return true, nil
}
