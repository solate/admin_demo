package dictrepo

import (
	"context"
	"time"

	"admin_backend/pkg/common"
	"admin_backend/pkg/common/contextutil"
	"admin_backend/pkg/ent"
	"admin_backend/pkg/ent/generated"
	"admin_backend/pkg/ent/generated/dicttype"
	"admin_backend/pkg/ent/generated/predicate"
)

type DictTypeRepo struct {
	db *ent.Client
}

// NewDictTypeRepo 创建字典类型仓储实例
func NewDictTypeRepo(db *ent.Client) *DictTypeRepo {
	return &DictTypeRepo{db: db}
}

func (r *DictTypeRepo) Create(ctx context.Context, dictType *generated.DictType) (*generated.DictType, error) {
	now := time.Now().UnixMilli()
	return r.db.DictType.Create().
		SetCreatedAt(now).
		SetUpdatedAt(now).
		SetTenantCode(dictType.TenantCode).
		SetTypeID(dictType.TypeID).
		SetCode(dictType.Code).
		SetName(dictType.Name).
		SetDescription(dictType.Description).
		SetStatus(dictType.Status).
		Save(ctx)
}

func (r *DictTypeRepo) Update(ctx context.Context, update *generated.DictType) (int, error) {
	now := time.Now().UnixMilli()
	update.UpdatedAt = now
	return r.db.DictType.Update().
		SetUpdatedAt(now).
		SetName(update.Name).
		SetDescription(update.Description).
		SetStatus(update.Status).
		Where(dicttype.TypeID(update.TypeID)).Save(ctx)
}

func (r *DictTypeRepo) GetByTypeID(ctx context.Context, typeID string) (*generated.DictType, error) {
	return r.Get(ctx, []predicate.DictType{dicttype.TypeID(typeID)})
}

func (r *DictTypeRepo) GetByCode(ctx context.Context, code string) (*generated.DictType, error) {
	return r.Get(ctx, []predicate.DictType{dicttype.Code(code)})
}

// defaultQuery 默认查询条件
func (r *DictTypeRepo) defaultQuery(ctx context.Context, where []predicate.DictType) []predicate.DictType {
	where = append(where, dicttype.DeletedAtIsNil())
	where = append(where, dicttype.TenantCode(contextutil.GetTenantCodeFromCtx(ctx)))
	return where
}

func (r *DictTypeRepo) Get(ctx context.Context, where []predicate.DictType) (*generated.DictType, error) {
	where = r.defaultQuery(ctx, where)
	return r.db.DictType.Query().Where(where...).Only(ctx)
}

func (r *DictTypeRepo) PageList(ctx context.Context, current, limit int, where []predicate.DictType) ([]*generated.DictType, int, error) {
	where = r.defaultQuery(ctx, where)

	offset := common.Offset(current, limit)
	query := r.db.DictType.Query().Where(where...).Order(generated.Desc(dicttype.FieldCreatedAt))

	// 查询总数
	total, err := query.Count(ctx)
	if err != nil || total == 0 {
		return nil, 0, err
	}

	// 分页查询
	list, err := query.Offset(offset).Limit(limit).All(ctx)
	return list, total, err
}

// DeleteByTypeID 根据字典类型ID删除字典类型，软删除
func (r *DictTypeRepo) Delete(ctx context.Context, typeID string) (int, error) {
	now := time.Now().UnixMilli()
	return r.db.DictType.Update().
		SetDeletedAt(now).
		Where(dicttype.TypeID(typeID)).Save(ctx)
}
