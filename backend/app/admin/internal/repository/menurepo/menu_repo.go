package menurepo

import (
	"context"
	"time"

	"admin_backend/pkg/common"
	"admin_backend/pkg/common/contextutil"
	"admin_backend/pkg/ent"
	"admin_backend/pkg/ent/generated"
	"admin_backend/pkg/ent/generated/menu"
	"admin_backend/pkg/ent/generated/predicate"
)

type MenuRepo struct {
	db *ent.Client
}

func NewMenuRepo(db *ent.Client) *MenuRepo {
	return &MenuRepo{db: db}
}

// GetByMenuCode 根据菜单编码获取菜单
func (r *MenuRepo) GetByMenuCode(ctx context.Context, menuCode string) (*generated.Menu, error) {
	return r.Get(ctx, []predicate.Menu{menu.Code(menuCode)})
}

// Create 创建菜单
func (r *MenuRepo) Create(ctx context.Context, menu *generated.Menu) (*generated.Menu, error) {
	now := time.Now().UnixMilli()
	return r.db.Menu.Create().
		SetCreatedAt(now).
		SetUpdatedAt(now).
		SetTenantCode(menu.TenantCode).
		SetMenuID(menu.MenuID).
		SetCode(menu.Code).
		SetParentID(menu.ParentID).
		SetName(menu.Name).
		SetPath(menu.Path).
		SetComponent(menu.Component).
		SetRedirect(menu.Redirect).
		SetIcon(menu.Icon).
		SetSort(menu.Sort).
		SetType(menu.Type).
		SetStatus(menu.Status).
		Save(ctx)
}

// Update 更新菜单
func (r *MenuRepo) Update(ctx context.Context, update *generated.Menu) (int, error) {
	now := time.Now().UnixMilli()
	update.UpdatedAt = now
	return r.db.Menu.Update().
		SetUpdatedAt(now).
		SetParentID(update.ParentID).
		SetName(update.Name).
		SetPath(update.Path).
		SetComponent(update.Component).
		SetRedirect(update.Redirect).
		SetIcon(update.Icon).
		SetSort(update.Sort).
		SetType(update.Type).
		SetStatus(update.Status).
		Where(menu.MenuID(update.MenuID)).Save(ctx)
}

// GetByMenuID 根据菜单ID获取菜单
func (r *MenuRepo) GetByMenuID(ctx context.Context, menuID string) (*generated.Menu, error) {
	return r.Get(ctx, []predicate.Menu{menu.MenuID(menuID)})
}

// defaultQuery 默认查询条件
func (r *MenuRepo) defaultQuery(ctx context.Context, where []predicate.Menu) []predicate.Menu {
	where = append(where, menu.DeletedAtIsNil())
	where = append(where, menu.TenantCode(contextutil.GetTenantCodeFromCtx(ctx)))
	return where
}

// Get 获取单个菜单
func (r *MenuRepo) Get(ctx context.Context, where []predicate.Menu) (*generated.Menu, error) {
	where = r.defaultQuery(ctx, where)
	return r.db.Menu.Query().Where(where...).Only(ctx)
}

// PageList 分页查询菜单列表
func (r *MenuRepo) PageList(ctx context.Context, current, limit int, where []predicate.Menu) ([]*generated.Menu, int, error) {
	where = r.defaultQuery(ctx, where)

	offset := common.Offset(current, limit)
	query := r.db.Menu.Query().Where(where...).Order(generated.Desc(menu.FieldCreatedAt))

	// 查询总数
	total, err := query.Count(ctx)
	if err != nil || total == 0 {
		return nil, 0, err
	}

	// 分页查询
	list, err := query.Offset(offset).Limit(limit).All(ctx)
	return list, total, err
}

// List 获取菜单列表（不分页）
func (r *MenuRepo) List(ctx context.Context, where []predicate.Menu) ([]*generated.Menu, error) {
	where = r.defaultQuery(ctx, where)
	return r.db.Menu.Query().Where(where...).Order(generated.Desc(menu.FieldCreatedAt)).All(ctx)
}

// Delete 删除菜单（软删除）
func (r *MenuRepo) Delete(ctx context.Context, menuID string) (int, error) {
	now := time.Now().UnixMilli()
	return r.db.Menu.Update().
		SetDeletedAt(now).
		Where(menu.MenuID(menuID)).Save(ctx)
}
