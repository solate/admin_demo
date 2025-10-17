package auth

import (
	"admin_backend/app/admin/internal/repository/loginlogrepo"
	"admin_backend/app/admin/internal/repository/sysuserrepo"
	"admin_backend/app/admin/internal/svc"
	"admin_backend/pkg/common/contextutil"
	"admin_backend/pkg/common/xerr"
	"admin_backend/pkg/ent/generated"
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type LogoutLogic struct {
	logx.Logger
	ctx          context.Context
	svcCtx       *svc.ServiceContext
	userRepo     *sysuserrepo.SysUserRepo
	loginLogRepo *loginlogrepo.LoginLogRepo
	r            *http.Request
}

// 用户登出
func NewLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext, r *http.Request) *LogoutLogic {
	return &LogoutLogic{
		Logger:       logx.WithContext(ctx),
		ctx:          ctx,
		svcCtx:       svcCtx,
		userRepo:     sysuserrepo.NewSysUserRepo(svcCtx.DB),
		loginLogRepo: loginlogrepo.NewLoginLogRepo(svcCtx.DB),
		r:            r,
	}
}

func (l *LogoutLogic) Logout() (resp bool, err error) {

	userID := contextutil.GetUserIDFromCtx(l.ctx)

	// 1. 查找用户
	user, err := l.userRepo.GetByUserID(l.ctx, userID)
	if err != nil {
		l.Error("Logout User.Query Error:", err.Error())
		if generated.IsNotFound(err) {
			return false, xerr.NewErrCodeMsg(xerr.DbError, "用户不存在")
		}
		return false, xerr.NewErrCodeMsg(xerr.DbError, err.Error())
	}

	// 更新user
	_, err = l.userRepo.Logout(l.ctx, userID)
	if err != nil {
		l.Error("Logout User.Logout Error:", err.Error())
		return false, xerr.NewErrCodeMsg(xerr.DbError, err.Error())
	}

	// 添加登出日志
	err = l.loginLogRepo.AddLoginLog(l.ctx, l.r, user, "登出成功")
	if err != nil {
		l.Error("Logout addLoginLog err:", err.Error())
	}

	return true, nil
}
