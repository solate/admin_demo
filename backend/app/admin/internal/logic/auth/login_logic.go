package auth

import (
	"context"
	"net/http"

	"admin_backend/app/admin/internal/repository/sysloginlogrepo"
	"admin_backend/app/admin/internal/repository/sysuserrepo"
	"admin_backend/app/admin/internal/svc"
	"admin_backend/app/admin/internal/types"
	"admin_backend/pkg/common/xerr"
	"admin_backend/pkg/constants"
	"admin_backend/pkg/ent/generated"
	"admin_backend/pkg/utils/passwordgen"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx          context.Context
	svcCtx       *svc.ServiceContext
	sysUserRepo  *sysuserrepo.SysUserRepo
	loginLogRepo *sysloginlogrepo.LoginLogRepo
	r            *http.Request
}

// 用户登录
func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext, r *http.Request) *LoginLogic {
	return &LoginLogic{
		Logger:       logx.WithContext(ctx),
		ctx:          ctx,
		svcCtx:       svcCtx,
		sysUserRepo:  sysuserrepo.NewSysUserRepo(svcCtx.DB),
		loginLogRepo: sysloginlogrepo.NewLoginLogRepo(svcCtx.DB),
		r:            r,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {

	// 1. 验证验证码
	if !l.svcCtx.CaptchaManager.Verify(req.CaptchaId, req.Captcha) {
		return nil, xerr.NewErrCodeMsg(xerr.ParamError, "验证码错误或已过期")
	}

	// 2. 查找用户
	sysUser, err := l.sysUserRepo.GetByUserName(l.ctx, req.UserName)
	if err != nil {
		l.Error("Login sysUserRepo.GetByUserName err:", err.Error())
		if generated.IsNotFound(err) {
			return nil, xerr.NewErrCodeMsg(xerr.RecordNotFound, "用户不存在")

		}
		return nil, xerr.NewErrCode(xerr.DbError)
	}

	// 2. 验证密码
	if !passwordgen.VerifyPassword(req.Password, sysUser.PwdHashed) {
		l.Error("Login passwordgen.VerifyPassword err")
		return nil, xerr.NewErrCodeMsg(xerr.ParamError, "密码错误")
	}

	// 3. 检查用户状态
	if sysUser.Status != 1 {
		l.Error("Login sysUser.Status != 1 err")
		return nil, xerr.NewErrCodeMsg(xerr.ParamError, "用户已被禁用")
	}

	// 4. 一行代码完成所有token操作
	result, err := l.svcCtx.JWTManager.Login(l.ctx, sysUser.UserID, sysUser.TenantCode, "", constants.SourcePC)
	if err != nil {
		l.Error("Login JWTManager.Login err:", err.Error())
		return nil, xerr.NewErrCode(xerr.ServerError)
	}

	// 5. 更新用户Token (保存access token)
	_, err = l.sysUserRepo.UpdateToken(l.ctx, sysUser.UserID, result.AccessToken)
	if err != nil {
		l.Error("Login sysUserRepo.UpdateToken err:", err.Error())
		return nil, xerr.NewErrCode(xerr.DbError)
	}

	// // 添加登录日志
	// err = l.loginLogRepo.AddLoginLog(l.ctx, l.r, sysUser, "登录成功")
	// if err != nil {
	// 	l.Error("Login addLoginLog err:", err.Error())
	// }

	// 6. 返回结果
	return &types.LoginResp{
		UserID:       sysUser.UserID,
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		UserName:     sysUser.UserName,
		Phone:        sysUser.Phone,
		Email:        sysUser.Email,
	}, nil
}
