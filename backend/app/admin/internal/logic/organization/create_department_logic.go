package organization

import (
	"context"

	"admin_backend/app/admin/internal/repository/organizationrepo"
	"admin_backend/app/admin/internal/svc"
	"admin_backend/app/admin/internal/types"
	"admin_backend/pkg/common/contextutil"
	"admin_backend/pkg/common/xerr"
	"admin_backend/pkg/ent/generated"
	"admin_backend/pkg/utils/idgen"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateDepartmentLogic struct {
	logx.Logger
	ctx     context.Context
	svcCtx  *svc.ServiceContext
	orgRepo *organizationrepo.DepartmentRepo
}

// 创建部门
func NewCreateDepartmentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateDepartmentLogic {
	return &CreateDepartmentLogic{
		Logger:  logx.WithContext(ctx),
		ctx:     ctx,
		svcCtx:  svcCtx,
		orgRepo: organizationrepo.NewDepartmentRepo(svcCtx.DB),
	}
}

func (l *CreateDepartmentLogic) CreateDepartment(req *types.CreateDepartmentReq) (resp *types.CreateDepartmentResp, err error) {
	// 1. 参数校验
	if req.Name == "" {
		return nil, xerr.NewErrCodeMsg(xerr.ParamError, "部门名称不能为空")
	}

	// 2. 检查部门名称是否已存在
	dept, err := l.orgRepo.GetByName(l.ctx, req.Name)
	if err != nil && !generated.IsNotFound(err) {
		l.Error("GetDepartment l.orgRepo.GetByName err: ", err.Error())
		return nil, xerr.NewErrCodeMsg(xerr.DbError, "数据库查询错误")
	}
	if dept != nil {
		l.Errorf("部门名称已存在: %s", req.Name)
		return nil, xerr.NewErrCodeMsg(xerr.DbRecordExist, "部门名称已存在")
	}

	// 2. 生成部门ID
	deptID, err := idgen.GenerateUUID()
	if err != nil {
		l.Error("CreateDepartment GenerateUUID err:", err.Error())
		return nil, xerr.NewErrCodeMsg(xerr.ServerError, "生成部门ID失败")
	}

	// 3. 设置默认值
	if req.Sort == 0 {
		req.Sort = 100 // 默认排序值
	}

	// 4. 创建部门
	newDept := &generated.Department{
		TenantCode:   contextutil.GetTenantCodeFromCtx(l.ctx),
		DepartmentID: deptID,
		Name:         req.Name,
		ParentID:     req.ParentID,
		Sort:         req.Sort,
	}

	_, err = l.orgRepo.Create(l.ctx, newDept)
	if err != nil {
		l.Error("CreateDepartment Create err:", err.Error())
		return nil, xerr.NewErrCodeMsg(xerr.ServerError, "创建部门失败")
	}

	// 5. 返回结果
	return &types.CreateDepartmentResp{
		DepartmentID: newDept.DepartmentID,
	}, nil
}
