package organization

import (
	"context"

	"admin_backend/app/admin/internal/repository/organizationrepo"
	"admin_backend/app/admin/internal/svc"
	"admin_backend/app/admin/internal/types"
	"admin_backend/pkg/common/xerr"
	"admin_backend/pkg/ent/generated/department"
	"admin_backend/pkg/ent/generated/predicate"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListDepartmentLogic struct {
	logx.Logger
	ctx     context.Context
	svcCtx  *svc.ServiceContext
	orgRepo *organizationrepo.DepartmentRepo
}

// 获取部门列表
func NewListDepartmentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListDepartmentLogic {
	return &ListDepartmentLogic{
		Logger:  logx.WithContext(ctx),
		ctx:     ctx,
		svcCtx:  svcCtx,
		orgRepo: organizationrepo.NewDepartmentRepo(svcCtx.DB),
	}
}

func (l *ListDepartmentLogic) ListDepartment(req *types.DepartmentListReq) (resp *types.DepartmentListResp, err error) {
	// 1. 构建查询条件
	var predicates []predicate.Department
	if req.Name != "" {
		predicates = append(predicates, department.NameContains(req.Name))
	}
	if req.ParentID != "" {
		predicates = append(predicates, department.ParentID(req.ParentID))
	}

	// 2. 查询部门列表
	list, total, err := l.orgRepo.PageList(l.ctx, req.Current, req.PageSize, predicates)
	if err != nil {
		l.Error("ListDepartment PageList err:", err.Error())
		return nil, xerr.NewErrCodeMsg(xerr.DbError, "查询部门列表失败")
	}

	// 3. 转换为响应结构
	deptList := make([]*types.DepartmentInfo, 0)
	for _, dept := range list {
		deptList = append(deptList, &types.DepartmentInfo{
			DepartmentID: dept.DepartmentID,
			Name:         dept.Name,
			ParentID:     dept.ParentID,
			Sort:         dept.Sort,
			CreatedAt:    dept.CreatedAt,
			UpdatedAt:    dept.UpdatedAt,
		})
	}

	// 4. 返回结果
	resp = &types.DepartmentListResp{
		Page: &types.PageResponse{
			Current:  req.Current,
			PageSize: req.PageSize,
			Total:    total,
		},
		List: deptList,
	}

	return resp, nil
}
