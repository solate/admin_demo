package permissionrepo

import (
	"context"
	"time"

	"admin_backend/pkg/common"
	"admin_backend/pkg/common/contextutil"
	"admin_backend/pkg/ent"
	"admin_backend/pkg/ent/generated"
	"admin_backend/pkg/ent/generated/permission"
	"admin_backend/pkg/ent/generated/predicate"
)

type PermissionRepo struct {
	db *ent.Client
}

func NewPermissionRepo(db *ent.Client) *PermissionRepo {
	return &PermissionRepo{db: db}
}

// Create 创建权限
func (r *PermissionRepo) Create(ctx context.Context, permission *generated.Permission) (*generated.Permission, error) {
	now := time.Now().UnixMilli()
	return r.db.Permission.Create().
		SetCreatedAt(now).
		SetUpdatedAt(now).
		SetTenantCode(contextutil.GetTenantCodeFromCtx(ctx)).
		SetPermissionID(permission.PermissionID).
		SetName(permission.Name).
		SetCode(permission.Code).
		SetType(permission.Type).
		SetResource(permission.Resource).
		SetAction(permission.Action).
		SetParentID(permission.ParentID).
		SetDescription(permission.Description).
		SetStatus(permission.Status).
		SetMenuID(permission.MenuID).
		Save(ctx)
}

// Delete 删除权限（软删除）
func (r *PermissionRepo) Delete(ctx context.Context, permissionID string) (int, error) {
	now := time.Now().UnixMilli()
	return r.db.Permission.Update().
		SetDeletedAt(now).
		Where(permission.PermissionID(permissionID)).
		Save(ctx)
}

// Update 更新权限
func (r *PermissionRepo) Update(ctx context.Context, update *generated.Permission) (int, error) {
	now := time.Now().UnixMilli()
	update.UpdatedAt = now
	return r.db.Permission.Update().
		SetUpdatedAt(now).
		SetName(update.Name).
		SetCode(update.Code).
		SetType(update.Type).
		SetResource(update.Resource).
		SetAction(update.Action).
		SetParentID(update.ParentID).
		SetDescription(update.Description).
		SetStatus(update.Status).
		SetMenuID(update.MenuID).
		Where(permission.PermissionID(update.PermissionID)).
		Save(ctx)
}

// defaultQuery 默认查询条件
func (r *PermissionRepo) defaultQuery(ctx context.Context, where []predicate.Permission) []predicate.Permission {
	where = append(where, permission.DeletedAtIsNil())
	where = append(where, permission.TenantCode(contextutil.GetTenantCodeFromCtx(ctx)))
	return where
}

// Get 获取权限
func (r *PermissionRepo) Get(ctx context.Context, where []predicate.Permission) (*generated.Permission, error) {
	where = r.defaultQuery(ctx, where)
	return r.db.Permission.Query().Where(where...).Only(ctx)
}

// GetByPermissionID 根据权限ID获取权限
func (r *PermissionRepo) GetByPermissionID(ctx context.Context, permissionID string) (*generated.Permission, error) {
	return r.Get(ctx, []predicate.Permission{permission.PermissionID(permissionID)})
}

// GetByCode 根据编码获取权限
func (r *PermissionRepo) GetByCode(ctx context.Context, code string) (*generated.Permission, error) {
	return r.Get(ctx, []predicate.Permission{permission.Code(code)})
}

// PageList 获取权限列表（分页）
func (r *PermissionRepo) PageList(ctx context.Context, current, limit int, name, code string, status int) ([]*generated.Permission, int, error) {
	where := []predicate.Permission{}

	if name != "" {
		where = append(where, permission.NameContains(name))
	}
	if code != "" {
		where = append(where, permission.CodeContains(code))
	}
	if status != 0 {
		where = append(where, permission.Status(status))
	}

	where = r.defaultQuery(ctx, where)

	offset := common.Offset(current, limit)
	query := r.db.Permission.Query().Where(where...).Order(generated.Desc(permission.FieldCreatedAt))

	// 查询总数
	total, err := query.Count(ctx)
	if err != nil || total == 0 {
		return nil, 0, err
	}

	// 分页查询
	permissions, err := query.Offset(offset).Limit(limit).All(ctx)
	return permissions, total, err
}

// List 获取权限列表（不分页）
func (r *PermissionRepo) List(ctx context.Context, where []predicate.Permission) ([]*generated.Permission, error) {
	where = r.defaultQuery(ctx, where)
	return r.db.Permission.Query().Where(where...).Order(generated.Desc(permission.FieldCreatedAt)).All(ctx)
}
