// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package sys_user

import (
	"context"

	"admin_backend/app/admin/internal/repository/sysuserrepo"
	"admin_backend/app/admin/internal/svc"
	"admin_backend/app/admin/internal/types"
	"admin_backend/pkg/common/contextutil"
	"admin_backend/pkg/common/xerr"
	"admin_backend/pkg/ent/generated"
	"admin_backend/pkg/utils/idgen"
	"admin_backend/pkg/utils/passwordgen"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateUserLogic struct {
	logx.Logger
	ctx      context.Context
	svcCtx   *svc.ServiceContext
	userRepo *sysuserrepo.SysUserRepo
}

// 创建用户
func NewCreateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateUserLogic {
	return &CreateUserLogic{
		Logger:   logx.WithContext(ctx),
		ctx:      ctx,
		svcCtx:   svcCtx,
		userRepo: sysuserrepo.NewSysUserRepo(svcCtx.DB),
	}
}

func (l *CreateUserLogic) CreateUser(req *types.CreateUserReq) (resp *types.CreateUserResp, err error) {
	// 1. 检查用户名是否已存在
	user, err := l.userRepo.GetByUserName(l.ctx, req.UserName)
	if err != nil && !generated.IsNotFound(err) {
		l.Error("CreateUser GetByUserName err: ", err.Error())
		return nil, xerr.NewErrCodeMsg(xerr.DbError, "数据库查询错误")
	}
	if user != nil {
		return nil, xerr.NewErrCodeMsg(xerr.DbRecordExist, "用户名已存在")
	}

	// 2. 检查手机号是否已存在
	user, err = l.userRepo.GetByPhone(l.ctx, req.Phone)
	if err != nil && !generated.IsNotFound(err) {
		l.Error("CreateUser GetByPhone err: ", err.Error())
		return nil, xerr.NewErrCodeMsg(xerr.DbError, "数据库查询错误")
	}
	if user != nil {
		return nil, xerr.NewErrCodeMsg(xerr.DbRecordExist, "手机号已存在")
	}

	// 3. 生成密码盐和加密密码
	salt, err := passwordgen.GenerateSalt()
	if err != nil {
		l.Error("CreateUser GenerateSalt err: ", err.Error())
		return nil, xerr.NewErrCodeMsg(xerr.ServerError, "生成密码盐失败")
	}
	hashedPassword, err := passwordgen.Argon2Hash(req.Password, salt)
	if err != nil {
		l.Error("CreateUser Argon2Hash err: ", err.Error())
		return nil, xerr.NewErrCodeMsg(xerr.ServerError, "密码加密失败")
	}

	// 4. 生成用户ID
	userID, err := idgen.GenerateUUID()
	if err != nil {
		l.Error("CreateUser GenerateUUID err: ", err.Error())
		return nil, xerr.NewErrCodeMsg(xerr.ServerError, "生成用户ID失败")
	}

	// 5. 获取租户代码
	tenantCode := contextutil.GetTenantCodeFromCtx(l.ctx)

	// 6. 创建用户
	newUser, err := l.userRepo.Create(l.ctx, &generated.SysUser{
		TenantCode: tenantCode,
		UserID:     userID,
		Phone:      req.Phone,
		UserName:   req.UserName,
		Name:       req.Name,
		Email:      req.Email,
		Sex:        req.Sex,
		Avatar:     req.Avatar,
		Status:     req.Status,
		PwdHashed:  hashedPassword,
		PwdSalt:    salt,
	})
	if err != nil {
		l.Error("CreateUser Create err:", err.Error())
		return nil, xerr.NewErrCodeMsg(xerr.DbError, "创建用户失败")
	}

	// 7. 返回结果
	return &types.CreateUserResp{
		UserID: newUser.UserID,
	}, nil
}
