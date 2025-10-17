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

type GetDepartmentLogic struct {
	logx.Logger
	ctx     context.Context
	svcCtx  *svc.ServiceContext
	orgRepo *organizationrepo.DepartmentRepo
}

// 获取部门详情
func NewGetDepartmentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDepartmentLogic {
	return &GetDepartmentLogic{
		Logger:  logx.WithContext(ctx),
		ctx:     ctx,
		svcCtx:  svcCtx,
		orgRepo: organizationrepo.NewDepartmentRepo(svcCtx.DB),
	}
}

func (l *GetDepartmentLogic) GetDepartment(req *types.GetDepartmentReq) (resp *types.DepartmentInfo, err error) {
	// 1. 查询部门信息
	dept, err := l.orgRepo.GetByDepartmentID(l.ctx, req.DepartmentID)
	if err != nil {
		if generated.IsNotFound(err) {
			return nil, xerr.NewErrCodeMsg(xerr.RecordNotFound, "部门不存在")
		}
		l.Error("GetDepartment l.orgRepo.GetByDepartmentID err: ", err.Error())
		return nil, xerr.NewErrCodeMsg(xerr.DbError, "数据库查询错误")
	}

	// 2. 转换为响应结构
	resp = &types.DepartmentInfo{
		DepartmentID: dept.DepartmentID,
		Name:         dept.Name,
		ParentID:     dept.ParentID,
		Sort:         dept.Sort,
		CreatedAt:    dept.CreatedAt,
		UpdatedAt:    dept.UpdatedAt,
	}

	return resp, nil
}
