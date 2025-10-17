package auth

import (
	"admin_backend/app/admin/internal/repository/sysuserrepo"
	"admin_backend/app/admin/internal/svc"
	"admin_backend/app/admin/internal/types"
	"admin_backend/pkg/common/contextutil"
	"admin_backend/pkg/common/xerr"
	"admin_backend/pkg/ent/generated"
	"admin_backend/pkg/utils/passwordgen"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChangePasswordLogic struct {
	logx.Logger
	ctx      context.Context
	svcCtx   *svc.ServiceContext
	userRepo *sysuserrepo.SysUserRepo
}

// 修改密码
func NewChangePasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangePasswordLogic {
	return &ChangePasswordLogic{
		Logger:   logx.WithContext(ctx),
		ctx:      ctx,
		svcCtx:   svcCtx,
		userRepo: sysuserrepo.NewSysUserRepo(svcCtx.DB),
	}
}

func (l *ChangePasswordLogic) ChangePassword(req *types.ChangePasswordReq) (resp bool, err error) {
	// 1. 获取当前用户ID
	userID := contextutil.GetUserIDFromCtx(l.ctx)

	// 2. 查找用户
	user, err := l.userRepo.GetByUserID(l.ctx, userID)
	if err != nil {
		l.Error("ChangePassword userRepo.GetByUserID err:", err.Error())
		if generated.IsNotFound(err) {
			return false, xerr.NewErrMsg("用户不存在")
		}
		return false, xerr.NewErrCode(xerr.DbError)
	}

	// 3. 验证原密码
	if !passwordgen.VerifyPassword(req.OldPassword, user.PwdHashed) {
		return false, xerr.NewErrCodeMsg(xerr.ParamError, "原密码错误")
	}

	// 4. 生成新的密码盐和加密密码
	salt, err := passwordgen.GenerateSalt()
	if err != nil {
		l.Error("ChangePassword GenerateSalt err:", err.Error())
		return false, xerr.NewErrCodeMsg(xerr.ServerError, "生成密码盐失败")
	}
	hashedPassword, err := passwordgen.Argon2Hash(req.NewPassword, salt)
	if err != nil {
		l.Error("ChangePassword Argon2Hash err:", err.Error())
		return false, xerr.NewErrCodeMsg(xerr.ServerError, "密码加密失败")
	}

	// 5. 更新用户密码
	user.PwdSalt = salt
	user.PwdHashed = hashedPassword
	_, err = l.userRepo.Update(l.ctx, user)
	if err != nil {
		l.Error("ChangePassword Update err:", err.Error())
		return false, xerr.NewErrCodeMsg(xerr.DbError, "更新密码失败")
	}

	return true, nil
}
