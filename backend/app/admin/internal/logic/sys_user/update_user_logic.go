// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package sys_user

import (
	"context"

	"admin_backend/app/admin/internal/repository/sysuserrepo"
	"admin_backend/app/admin/internal/svc"
	"admin_backend/app/admin/internal/types"
	"admin_backend/pkg/common/xerr"
	"admin_backend/pkg/ent/generated"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserLogic struct {
	logx.Logger
	ctx      context.Context
	svcCtx   *svc.ServiceContext
	userRepo *sysuserrepo.SysUserRepo
}

// 更新用户
func NewUpdateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserLogic {
	return &UpdateUserLogic{
		Logger:   logx.WithContext(ctx),
		ctx:      ctx,
		svcCtx:   svcCtx,
		userRepo: sysuserrepo.NewSysUserRepo(svcCtx.DB),
	}
}

func (l *UpdateUserLogic) UpdateUser(req *types.UpdateUserReq) (resp bool, err error) {
	// 1. 检查用户是否存在
	user, err := l.userRepo.GetByUserID(l.ctx, req.UserID)
	if err != nil {
		if generated.IsNotFound(err) {
			return false, xerr.NewErrCodeMsg(xerr.RecordNotFound, "用户不存在")
		}
		l.Error("UpdateUser GetByUserID err: ", err.Error())
		return false, xerr.NewErrCodeMsg(xerr.DbError, "数据库查询错误")
	}

	// 2. 更新用户信息
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Status != 0 {
		user.Status = req.Status
	}
	if req.Sex != 0 {
		user.Sex = req.Sex
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}

	// 3. 保存更新
	_, err = l.userRepo.Update(l.ctx, user)
	if err != nil {
		l.Error("UpdateUser Update err:", err.Error())
		return false, xerr.NewErrCodeMsg(xerr.DbError, "更新用户失败")
	}

	return true, nil
}
