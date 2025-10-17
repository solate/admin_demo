package organizationrepo

import (
	"context"

	"admin_backend/pkg/ent"
	"admin_backend/pkg/ent/generated"
	"admin_backend/pkg/ent/generated/predicate"
	"admin_backend/pkg/ent/generated/userposition"
)

type UserPositionRepo struct {
	db *ent.Client
}

// NewUserPositionRepo 创建用户岗位关联仓储实例
func NewUserPositionRepo(db *ent.Client) *UserPositionRepo {
	return &UserPositionRepo{db: db}
}

// Create 创建用户岗位关联
func (r *UserPositionRepo) Create(ctx context.Context, userPosition *generated.UserPosition) (*generated.UserPosition, error) {
	return r.db.UserPosition.Create().
		SetUserID(userPosition.UserID).
		SetPositionID(userPosition.PositionID).
		Save(ctx)
}

// Delete 删除用户岗位关联
func (r *UserPositionRepo) Delete(ctx context.Context, userID, positionID string) (int, error) {
	return r.db.UserPosition.Update().
		Where(
			userposition.And(
				userposition.UserID(userID),
				userposition.PositionID(positionID),
			),
		).
		Save(ctx)
}

// GetUserPositions 获取用户的岗位列表
func (r *UserPositionRepo) GetUserPositions(ctx context.Context, userID string) ([]*generated.UserPosition, error) {
	where := []predicate.UserPosition{userposition.UserID(userID)}
	return r.db.UserPosition.Query().Where(where...).All(ctx)
}

// GetPositionUsers 获取岗位下的用户列表
func (r *UserPositionRepo) GetPositionUsers(ctx context.Context, positionID string) ([]*generated.UserPosition, error) {
	where := []predicate.UserPosition{userposition.PositionID(positionID)}
	return r.db.UserPosition.Query().Where(where...).All(ctx)
}
