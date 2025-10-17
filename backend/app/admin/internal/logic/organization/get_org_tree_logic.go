package organization

import (
	"context"

	"admin_backend/app/admin/internal/repository/organizationrepo"
	"admin_backend/app/admin/internal/svc"
	"admin_backend/app/admin/internal/types"
	"admin_backend/pkg/common/contextutil"
	"admin_backend/pkg/common/xerr"
	"admin_backend/pkg/ent/generated/department"
	"admin_backend/pkg/ent/generated/predicate"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOrgTreeLogic struct {
	logx.Logger
	ctx     context.Context
	svcCtx  *svc.ServiceContext
	orgRepo *organizationrepo.DepartmentRepo
}

// 获取组织架构树
func NewGetOrgTreeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOrgTreeLogic {
	return &GetOrgTreeLogic{
		Logger:  logx.WithContext(ctx),
		ctx:     ctx,
		svcCtx:  svcCtx,
		orgRepo: organizationrepo.NewDepartmentRepo(svcCtx.DB),
	}
}

func (l *GetOrgTreeLogic) GetOrgTree(req *types.GetOrgTreeReq) (resp *types.GetOrgTreeResp, err error) {
	// 1. 获取租户编码
	tenantCode := contextutil.GetTenantCodeFromCtx(l.ctx)
	// 2. 获取用户角色编码
	roleCode := contextutil.GetRoleCodeFromCtx(l.ctx)

	// 3. 构建查询条件
	var predicates []predicate.Department
	// 添加租户条件
	predicates = append(predicates, department.TenantCode(tenantCode))

	// 4. 获取角色权限
	pm := l.svcCtx.CasbinManager
	permissions, err := pm.GetRolePermissions(roleCode, tenantCode)
	if err != nil {
		l.Errorf("获取角色权限失败: %v", err)
		return nil, xerr.NewErrCodeMsg(xerr.ServerError, "获取角色权限失败")
	}

	// 5. 获取有权限的部门ID列表
	var authorizedDeptIDs []string
	for _, permission := range permissions {
		resource := permission[2]
		// 假设资源格式为 "dept:view:deptID"
		// 这里可以根据实际的权限格式进行解析
		authorizedDeptIDs = append(authorizedDeptIDs, resource)
	}

	// 6. 查询所有部门
	depts, err := l.orgRepo.List(l.ctx, predicates)
	if err != nil {
		l.Error("GetOrgTree List err:", err.Error())
		return nil, xerr.NewErrCodeMsg(xerr.DbError, "查询部门列表失败")
	}

	// 7. 构建部门树
	deptMap := make(map[string]*types.OrgTreeNode)
	for _, dept := range depts {
		// 检查是否有权限访问该部门
		hasPermission := false
		for _, authorizedID := range authorizedDeptIDs {
			if authorizedID == dept.DepartmentID {
				hasPermission = true
				break
			}
		}
		if !hasPermission {
			continue
		}

		deptMap[dept.DepartmentID] = &types.OrgTreeNode{
			DepartmentID: dept.DepartmentID,
			Name:         dept.Name,
			ParentID:     dept.ParentID,
			Sort:         dept.Sort,
			Children:     make([]*types.OrgTreeNode, 0),
		}
	}

	// 8. 构建树形结构
	tree := buildDepartmentTree(deptMap, "")

	// 9. 返回结果
	resp = &types.GetOrgTreeResp{
		Tree: tree,
	}

	return resp, nil
}

// buildDepartmentTree 构建部门树
func buildDepartmentTree(deptMap map[string]*types.OrgTreeNode, parentID string) []*types.OrgTreeNode {
	var tree []*types.OrgTreeNode
	for _, dept := range deptMap {
		if dept.ParentID == parentID {
			// 递归获取子部门
			dept.Children = buildDepartmentTree(deptMap, dept.DepartmentID)
			tree = append(tree, dept)
		}
	}
	return tree
}
