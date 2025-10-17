package rolerepo

import (
	"context"
	"time"

	"admin_backend/pkg/common"
	"admin_backend/pkg/common/contextutil"
	"admin_backend/pkg/ent"
	"admin_backend/pkg/ent/generated"
	"admin_backend/pkg/ent/generated/predicate"
	"admin_backend/pkg/ent/generated/role"
)

type RoleRepo struct {
	db *ent.Client
}

func NewRoleRepo(db *ent.Client) *RoleRepo {
	return &RoleRepo{db: db}
}

// Create 创建角色
func (r *RoleRepo) Create(ctx context.Context, role *generated.Role) (*generated.Role, error) {
	now := time.Now().UnixMilli()
	return r.db.Role.Create().
		SetCreatedAt(now).
		SetUpdatedAt(now).
		SetTenantCode(contextutil.GetTenantCodeFromCtx(ctx)).
		SetRoleID(role.RoleID).
		SetName(role.Name).
		SetCode(role.Code).
		SetDescription(role.Description).
		SetStatus(role.Status).
		SetSort(role.Sort).
		Save(ctx)
}

// Delete 删除角色（软删除）
func (r *RoleRepo) Delete(ctx context.Context, roleID string) (int, error) {
	now := time.Now().UnixMilli()
	return r.db.Role.Update().
		SetDeletedAt(now).
		Where(role.RoleID(roleID)).
		Save(ctx)
}

// Update 更新角色
func (r *RoleRepo) Update(ctx context.Context, update *generated.Role) (int, error) {
	now := time.Now().UnixMilli()
	update.UpdatedAt = now
	return r.db.Role.Update().
		SetUpdatedAt(now).
		SetName(update.Name).
		SetCode(update.Code).
		SetDescription(update.Description).
		SetStatus(update.Status).
		SetSort(update.Sort).
		Where(role.RoleID(update.RoleID)).
		Save(ctx)
}

// defaultQuery 默认查询条件
func (r *RoleRepo) defaultQuery(ctx context.Context, where []predicate.Role) []predicate.Role {
	where = append(where, role.DeletedAtIsNil())
	where = append(where, role.TenantCode(contextutil.GetTenantCodeFromCtx(ctx)))
	return where
}

// Get 获取角色
func (r *RoleRepo) Get(ctx context.Context, where []predicate.Role) (*generated.Role, error) {
	where = r.defaultQuery(ctx, where)
	return r.db.Role.Query().Where(where...).Only(ctx)
}

// GetByRoleID 根据角色ID获取角色
func (r *RoleRepo) GetByRoleID(ctx context.Context, roleID string) (*generated.Role, error) {
	return r.Get(ctx, []predicate.Role{role.RoleID(roleID)})
}

// GetByCode 根据编码获取角色
func (r *RoleRepo) GetByCode(ctx context.Context, code string) (*generated.Role, error) {
	return r.Get(ctx, []predicate.Role{role.Code(code)})
}

// List 获取角色列表
func (r *RoleRepo) PageList(ctx context.Context, current, limit int, name, code string, status int) ([]*generated.Role, int, error) {
	where := []predicate.Role{}

	if name != "" {
		where = append(where, role.NameContains(name))
	}
	if code != "" {
		where = append(where, role.CodeContains(code))
	}
	if status != 0 {
		where = append(where, role.Status(status))
	}

	where = r.defaultQuery(ctx, where)

	offset := common.Offset(current, limit)
	query := r.db.Role.Query().Where(where...).Order(generated.Desc(role.FieldCreatedAt))

	// 查询总数
	total, err := query.Count(ctx)
	if err != nil || total == 0 {
		return nil, 0, err
	}

	// 分页查询
	roles, err := query.Offset(offset).Limit(limit).All(ctx)
	return roles, total, err
}

// List 获取角色列表（不分页）
func (r *RoleRepo) List(ctx context.Context, where []predicate.Role) ([]*generated.Role, error) {
	where = r.defaultQuery(ctx, where)
	return r.db.Role.Query().Where(where...).Order(generated.Desc(role.FieldCreatedAt)).All(ctx)
}
