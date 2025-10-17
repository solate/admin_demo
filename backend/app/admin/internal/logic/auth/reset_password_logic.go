package auth

import (
	"admin_backend/app/admin/internal/repository/sysuserrepo"
	"admin_backend/app/admin/internal/svc"
	"admin_backend/app/admin/internal/types"
	"admin_backend/pkg/common/xerr"
	"admin_backend/pkg/ent/generated"
	"admin_backend/pkg/utils/passwordgen"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type ResetPasswordLogic struct {
	logx.Logger
	ctx      context.Context
	svcCtx   *svc.ServiceContext
	userRepo *sysuserrepo.SysUserRepo
}

// 重置密码
func NewResetPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetPasswordLogic {
	return &ResetPasswordLogic{
		Logger:   logx.WithContext(ctx),
		ctx:      ctx,
		svcCtx:   svcCtx,
		userRepo: sysuserrepo.NewSysUserRepo(svcCtx.DB),
	}
}

func (l *ResetPasswordLogic) ResetPassword(req *types.ResetPasswordReq) (resp bool, err error) {
	// 1. 查找用户
	user, err := l.userRepo.GetByUserID(l.ctx, req.UserID)
	if err != nil {
		l.Error("ResetPassword userRepo.GetByUserID err:", err.Error())
		if generated.IsNotFound(err) {
			return false, xerr.NewErrMsg("用户不存在")
		}
		return false, xerr.NewErrCode(xerr.DbError)
	}

	// 2. 生成新的密码盐和加密密码
	salt, err := passwordgen.GenerateSalt()
	if err != nil {
		l.Error("ResetPassword GenerateSalt err:", err.Error())
		return false, xerr.NewErrCodeMsg(xerr.ServerError, "生成密码盐失败")
	}
	hashedPassword, err := passwordgen.Argon2Hash(req.NewPassword, salt)
	if err != nil {
		l.Error("ResetPassword Argon2Hash err:", err.Error())
		return false, xerr.NewErrCodeMsg(xerr.ServerError, "密码加密失败")
	}

	// 3. 更新用户密码
	user.PwdSalt = salt
	user.PwdHashed = hashedPassword
	_, err = l.userRepo.Update(l.ctx, user)
	if err != nil {
		l.Error("ResetPassword Update err:", err.Error())
		return false, xerr.NewErrCodeMsg(xerr.DbError, "更新密码失败")
	}

	return true, nil
}
