package organizationrepo

import (
	"context"
	"time"

	"admin_backend/pkg/common"
	"admin_backend/pkg/common/contextutil"
	"admin_backend/pkg/ent"
	"admin_backend/pkg/ent/generated"
	"admin_backend/pkg/ent/generated/department"
	"admin_backend/pkg/ent/generated/predicate"
)

type DepartmentRepo struct {
	db *ent.Client
}

// NewDepartmentRepo 创建部门仓储实例
func NewDepartmentRepo(db *ent.Client) *DepartmentRepo {
	return &DepartmentRepo{db: db}
}

func (r *DepartmentRepo) Create(ctx context.Context, dept *generated.Department) (*generated.Department, error) {
	now := time.Now().UnixMilli()
	return r.db.Department.Create().
		SetCreatedAt(now).
		SetUpdatedAt(now).
		SetTenantCode(dept.TenantCode).
		SetDepartmentID(dept.DepartmentID).
		SetName(dept.Name).
		SetParentID(dept.ParentID).
		SetSort(dept.Sort).
		Save(ctx)
}

func (r *DepartmentRepo) Update(ctx context.Context, update *generated.Department) (int, error) {
	now := time.Now().UnixMilli()
	update.UpdatedAt = now
	return r.db.Department.Update().
		SetUpdatedAt(now).
		SetName(update.Name).
		SetParentID(update.ParentID).
		SetSort(update.Sort).
		Where(department.DepartmentID(update.DepartmentID)).Save(ctx)
}

func (r *DepartmentRepo) GetByDepartmentID(ctx context.Context, departmentID string) (*generated.Department, error) {
	return r.Get(ctx, []predicate.Department{department.DepartmentID(departmentID)})
}

func (r *DepartmentRepo) GetByName(ctx context.Context, name string) (*generated.Department, error) {
	return r.Get(ctx, []predicate.Department{department.Name(name)})
}

// defaultQuery 默认查询条件
func (r *DepartmentRepo) defaultQuery(ctx context.Context, where []predicate.Department) []predicate.Department {
	where = append(where, department.DeletedAtIsNil())
	where = append(where, department.TenantCode(contextutil.GetTenantCodeFromCtx(ctx)))
	return where
}

func (r *DepartmentRepo) Get(ctx context.Context, where []predicate.Department) (*generated.Department, error) {
	where = r.defaultQuery(ctx, where)
	return r.db.Department.Query().Where(where...).Only(ctx)
}

func (r *DepartmentRepo) PageList(ctx context.Context, current, limit int, where []predicate.Department) ([]*generated.Department, int, error) {
	where = r.defaultQuery(ctx, where)

	offset := common.Offset(current, limit)
	query := r.db.Department.Query().Where(where...).Order(generated.Desc(department.FieldCreatedAt))

	// 查询总数
	total, err := query.Count(ctx)
	if err != nil || total == 0 {
		return nil, 0, err
	}

	// 分页查询
	list, err := query.Offset(offset).Limit(limit).All(ctx)
	return list, total, err
}

// Delete 根据部门ID删除部门，软删除
func (r *DepartmentRepo) Delete(ctx context.Context, departmentID string) (int, error) {
	now := time.Now().UnixMilli()
	return r.db.Department.Update().
		SetDeletedAt(now).
		Where(department.DepartmentID(departmentID)).Save(ctx)
}

// List 获取菜单列表（不分页）
func (r *DepartmentRepo) List(ctx context.Context, where []predicate.Department) ([]*generated.Department, error) {
	return r.db.Department.Query().Where(where...).Order(generated.Desc(department.FieldCreatedAt)).All(ctx)
}
