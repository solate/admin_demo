package auth

import (
	"context"

	"admin_backend/app/admin/internal/repository/loginlogrepo"
	"admin_backend/app/admin/internal/repository/sysuserrepo"
	"admin_backend/app/admin/internal/svc"
	"admin_backend/pkg/common/contextutil"
	"admin_backend/pkg/common/xerr"
	"admin_backend/pkg/ent/generated"

	"github.com/zeromicro/go-zero/core/logx"
)

type LogoutAllLogic struct {
	logx.Logger
	ctx          context.Context
	svcCtx       *svc.ServiceContext
	userRepo     *sysuserrepo.SysUserRepo
	loginLogRepo *loginlogrepo.LoginLogRepo
}

// 用户登出（所有设备）
func NewLogoutAllLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogoutAllLogic {
	return &LogoutAllLogic{
		Logger:       logx.WithContext(ctx),
		ctx:          ctx,
		svcCtx:       svcCtx,
		userRepo:     sysuserrepo.NewSysUserRepo(svcCtx.DB),
		loginLogRepo: loginlogrepo.NewLoginLogRepo(svcCtx.DB),
	}
}

func (l *LogoutAllLogic) LogoutAll() (resp bool, err error) {
	// 1. 获取当前用户ID
	userID := contextutil.GetUserIDFromCtx(l.ctx)

	// 2. 查找用户
	user, err := l.userRepo.GetByUserID(l.ctx, userID)
	if err != nil {
		l.Error("LogoutAll User.Query Error:", err.Error())
		if generated.IsNotFound(err) {
			return false, xerr.NewErrCodeMsg(xerr.DbError, "用户不存在")
		}
		return false, xerr.NewErrCodeMsg(xerr.DbError, err.Error())
	}

	// 3. 使用JWT Manager登出所有设备
	err = l.svcCtx.JWTManager.LogoutAllDevices(l.ctx, userID)
	if err != nil {
		l.Error("LogoutAll LogoutAllDevices Error:", err.Error())
		return false, xerr.NewErrCodeMsg(xerr.ServerError, "登出所有设备失败")
	}

	// 4. 清空数据库中的token字段
	_, err = l.userRepo.Logout(l.ctx, userID)
	if err != nil {
		l.Error("LogoutAll User.Logout Error:", err.Error())
		// 数据库更新失败不影响主流程，因为JWT Manager已经处理了
		l.Error("数据库token清空失败，但JWT已失效: userID=%s", userID)
	}

	// 5. 添加登出日志（可选）
	// 注意：这里无法获取request对象，需要修改接口或忽略request参数
	// err = l.loginLogRepo.AddLoginLog(l.ctx, nil, user, "登出所有设备成功")
	// if err != nil {
	// 	l.Error("LogoutAll addLoginLog err:", err.Error())
	// }

	l.Infof("LogoutAll success: userID=%s, username=%s", user.UserID, user.UserName)
	return true, nil
}
