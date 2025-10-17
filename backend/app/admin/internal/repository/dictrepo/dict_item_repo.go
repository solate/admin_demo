package dictrepo

import (
	"context"
	"time"

	"admin_backend/pkg/common"
	"admin_backend/pkg/common/contextutil"
	"admin_backend/pkg/ent"
	"admin_backend/pkg/ent/generated"
	"admin_backend/pkg/ent/generated/dictitem"
	"admin_backend/pkg/ent/generated/predicate"
)

type DictItemRepo struct {
	db *ent.Client
}

// NewDictItemRepo 创建字典数据项仓储实例
func NewDictItemRepo(db *ent.Client) *DictItemRepo {
	return &DictItemRepo{db: db}
}

func (r *DictItemRepo) Create(ctx context.Context, item *generated.DictItem) (*generated.DictItem, error) {
	now := time.Now().UnixMilli()
	return r.db.DictItem.Create().
		SetCreatedAt(now).
		SetUpdatedAt(now).
		SetTenantCode(item.TenantCode).
		SetItemID(item.ItemID).
		SetTypeCode(item.TypeCode).
		SetLabel(item.Label).
		SetValue(item.Value).
		SetDescription(item.Description).
		SetSort(item.Sort).
		SetStatus(item.Status).
		Save(ctx)
}

func (r *DictItemRepo) Update(ctx context.Context, update *generated.DictItem) (int, error) {
	now := time.Now().UnixMilli()
	update.UpdatedAt = now
	return r.db.DictItem.Update().
		SetUpdatedAt(now).
		SetLabel(update.Label).
		SetValue(update.Value).
		SetDescription(update.Description).
		SetSort(update.Sort).
		SetStatus(update.Status).
		Where(dictitem.ItemID(update.ItemID)).Save(ctx)
}

func (r *DictItemRepo) GetByItemID(ctx context.Context, itemID string) (*generated.DictItem, error) {
	return r.Get(ctx, []predicate.DictItem{dictitem.ItemID(itemID)})
}

func (r *DictItemRepo) GetByTypeCode(ctx context.Context, typeCode string) ([]*generated.DictItem, error) {
	where := []predicate.DictItem{dictitem.TypeCode(typeCode)}
	where = r.defaultQuery(ctx, where)
	return r.db.DictItem.Query().Where(where...).All(ctx)
}

// defaultQuery 默认查询条件
func (r *DictItemRepo) defaultQuery(ctx context.Context, where []predicate.DictItem) []predicate.DictItem {
	where = append(where, dictitem.DeletedAtIsNil())
	where = append(where, dictitem.TenantCode(contextutil.GetTenantCodeFromCtx(ctx)))
	return where
}

func (r *DictItemRepo) Get(ctx context.Context, where []predicate.DictItem) (*generated.DictItem, error) {
	where = r.defaultQuery(ctx, where)
	return r.db.DictItem.Query().Where(where...).Only(ctx)
}

func (r *DictItemRepo) PageList(ctx context.Context, current, limit int, where []predicate.DictItem) ([]*generated.DictItem, int, error) {
	where = r.defaultQuery(ctx, where)

	offset := common.Offset(current, limit)
	query := r.db.DictItem.Query().Where(where...).Order(generated.Desc(dictitem.FieldCreatedAt))

	// 查询总数
	total, err := query.Count(ctx)
	if err != nil || total == 0 {
		return nil, 0, err
	}

	// 分页查询
	list, err := query.Offset(offset).Limit(limit).All(ctx)
	return list, total, err
}

// DeleteByItemID 根据字典数据项ID删除字典数据项，软删除
func (r *DictItemRepo) Delete(ctx context.Context, itemID string) (int, error) {
	now := time.Now().UnixMilli()
	return r.db.DictItem.Update().
		SetDeletedAt(now).
		Where(dictitem.ItemID(itemID)).Save(ctx)
}
