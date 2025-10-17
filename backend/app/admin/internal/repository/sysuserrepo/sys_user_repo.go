package sysuserrepo

import (
	"context"
	"time"

	"admin_backend/pkg/common"
	"admin_backend/pkg/common/contextutil"
	"admin_backend/pkg/ent"
	"admin_backend/pkg/ent/generated"
	"admin_backend/pkg/ent/generated/predicate"
	"admin_backend/pkg/ent/generated/sysuser"
)

type SysUserRepo struct {
	db *ent.Client
}

// NewSysUserRepo 创建用户仓储实例
func NewSysUserRepo(db *ent.Client) *SysUserRepo {
	return &SysUserRepo{db: db}
}

func (r *SysUserRepo) Create(ctx context.Context, sysUser *generated.SysUser) (*generated.SysUser, error) {
	now := time.Now().UnixMilli()
	return r.db.SysUser.Create().
		SetCreatedAt(now).
		SetUpdatedAt(now).
		SetTenantCode(sysUser.TenantCode).
		SetUserID(sysUser.UserID).
		SetPhone(sysUser.Phone).
		SetUserName(sysUser.UserName).
		SetPwdHashed(sysUser.PwdHashed).
		SetPwdSalt(sysUser.PwdSalt).
		SetStatus(sysUser.Status).
		SetName(sysUser.Name).
		SetEmail(sysUser.Email).
		SetSex(sysUser.Sex).
		Save(ctx)
}

func (r *SysUserRepo) Update(ctx context.Context, update *generated.SysUser) (int, error) {
	now := time.Now().UnixMilli()
	update.UpdatedAt = now
	return r.db.SysUser.Update().
		SetUpdatedAt(now).
		SetPhone(update.Phone).
		SetUserName(update.UserName).
		SetName(update.Name).
		SetEmail(update.Email).
		SetSex(update.Sex).
		SetStatus(update.Status).
		SetPwdHashed(update.PwdHashed).
		SetPwdSalt(update.PwdSalt).
		Where(sysuser.UserID(update.UserID)).Save(ctx)
}

// func (r *UserRepo) GetByID(ctx context.Context, id int) (*generated.User, error) {
// 	return r.db.User.Query().Where(user.ID(id)).Only(ctx)
// }

func (r *SysUserRepo) GetByUserID(ctx context.Context, userID string) (*generated.SysUser, error) {
	return r.Get(ctx, []predicate.SysUser{sysuser.UserID(userID)})
}

func (r *SysUserRepo) GetByPhone(ctx context.Context, phone string) (*generated.SysUser, error) {
	return r.Get(ctx, []predicate.SysUser{sysuser.Phone(phone)})
}

func (r *SysUserRepo) GetByUserName(ctx context.Context, userName string) (*generated.SysUser, error) {
	return r.Get(ctx, []predicate.SysUser{sysuser.UserName(userName)})
}

// defaultQuery 默认查询条件
func (r *SysUserRepo) defaultQuery(ctx context.Context, where []predicate.SysUser) []predicate.SysUser {
	where = append(where, sysuser.DeletedAtIsNil())
	where = append(where, sysuser.TenantCode(contextutil.GetTenantCodeFromCtx(ctx)))
	return where
}

func (r *SysUserRepo) Get(ctx context.Context, where []predicate.SysUser) (*generated.SysUser, error) {
	where = r.defaultQuery(ctx, where)
	return r.db.SysUser.Query().Where(where...).Only(ctx)
}

func (r *SysUserRepo) PageList(ctx context.Context, current, limit int, where []predicate.SysUser) ([]*generated.SysUser, int, error) {

	where = r.defaultQuery(ctx, where)

	offset := common.Offset(current, limit)
	query := r.db.SysUser.Query().Where(where...).Order(generated.Desc(sysuser.FieldCreatedAt))

	// 查询总数
	total, err := query.Count(ctx)
	if err != nil || total == 0 {
		return nil, 0, err
	}

	// 分页查询
	list, err := query.Offset(offset).Limit(limit).All(ctx)
	return list, total, err
}

// DeleteByUserID 根据用户ID删除用户，软删除
func (r *SysUserRepo) Delete(ctx context.Context, userID string) (int, error) {
	now := time.Now().UnixMilli()
	return r.db.SysUser.Update().
		SetDeletedAt(now).
		Where(sysuser.UserID(userID)).Save(ctx)
}

// 用户登出
func (r *SysUserRepo) Logout(ctx context.Context, userID string) (int, error) {
	return r.db.SysUser.Update().
		SetToken("").
		Where(sysuser.UserID(userID)).Save(ctx)
}

// 用户登录
func (r *SysUserRepo) UpdateToken(ctx context.Context, userID string, token string) (int, error) {
	return r.db.SysUser.Update().
		SetToken(token).
		SetUpdatedAt(time.Now().UnixMilli()).
		Where(sysuser.UserID(userID)).Save(ctx)
}
