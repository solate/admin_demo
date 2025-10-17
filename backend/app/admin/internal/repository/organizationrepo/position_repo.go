package organizationrepo

import (
	"context"
	"time"

	"admin_backend/pkg/common"
	"admin_backend/pkg/common/contextutil"
	"admin_backend/pkg/ent"
	"admin_backend/pkg/ent/generated"
	"admin_backend/pkg/ent/generated/position"
	"admin_backend/pkg/ent/generated/predicate"
)

type PositionRepo struct {
	db *ent.Client
}

// NewPositionRepo 创建岗位仓储实例
func NewPositionRepo(db *ent.Client) *PositionRepo {
	return &PositionRepo{db: db}
}

func (r *PositionRepo) Create(ctx context.Context, pos *generated.Position) (*generated.Position, error) {
	now := time.Now().UnixMilli()
	return r.db.Position.Create().
		SetCreatedAt(now).
		SetUpdatedAt(now).
		SetTenantCode(pos.TenantCode).
		SetPositionID(pos.PositionID).
		SetName(pos.Name).
		SetDepartmentID(pos.DepartmentID).
		SetDescription(pos.Description).
		Save(ctx)
}

func (r *PositionRepo) Update(ctx context.Context, update *generated.Position) (int, error) {
	now := time.Now().UnixMilli()
	update.UpdatedAt = now
	return r.db.Position.Update().
		SetUpdatedAt(now).
		SetName(update.Name).
		SetDepartmentID(update.DepartmentID).
		SetDescription(update.Description).
		Where(position.PositionID(update.PositionID)).Save(ctx)
}

func (r *PositionRepo) GetByPositionID(ctx context.Context, positionID string) (*generated.Position, error) {
	return r.Get(ctx, []predicate.Position{position.PositionID(positionID)})
}

func (r *PositionRepo) GetByName(ctx context.Context, name string) (*generated.Position, error) {
	return r.Get(ctx, []predicate.Position{position.Name(name)})
}

// defaultQuery 默认查询条件
func (r *PositionRepo) defaultQuery(ctx context.Context, where []predicate.Position) []predicate.Position {
	where = append(where, position.DeletedAtIsNil())
	where = append(where, position.TenantCode(contextutil.GetTenantCodeFromCtx(ctx)))
	return where
}

func (r *PositionRepo) Get(ctx context.Context, where []predicate.Position) (*generated.Position, error) {
	where = r.defaultQuery(ctx, where)
	return r.db.Position.Query().Where(where...).Only(ctx)
}

func (r *PositionRepo) PageList(ctx context.Context, current, limit int, where []predicate.Position) ([]*generated.Position, int, error) {
	where = r.defaultQuery(ctx, where)

	offset := common.Offset(current, limit)
	query := r.db.Position.Query().Where(where...).Order(generated.Desc(position.FieldCreatedAt))

	// 查询总数
	total, err := query.Count(ctx)
	if err != nil || total == 0 {
		return nil, 0, err
	}

	// 分页查询
	list, err := query.Offset(offset).Limit(limit).All(ctx)
	return list, total, err
}

// Delete 根据岗位ID删除岗位，软删除
func (r *PositionRepo) Delete(ctx context.Context, positionID string) (int, error) {
	now := time.Now().UnixMilli()
	return r.db.Position.Update().
		SetDeletedAt(now).
		Where(position.PositionID(positionID)).Save(ctx)
}
