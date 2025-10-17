package auth

import (
	"admin_backend/app/admin/internal/repository/sysuserrepo"
	"admin_backend/app/admin/internal/svc"
	"admin_backend/app/admin/internal/types"
	"admin_backend/pkg/common/xerr"
	"admin_backend/pkg/constants"
	"admin_backend/pkg/ent/generated"
	"admin_backend/pkg/utils/idgen"
	"admin_backend/pkg/utils/passwordgen"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx      context.Context
	svcCtx   *svc.ServiceContext
	userRepo *sysuserrepo.SysUserRepo
}

// 用户注册
func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger:   logx.WithContext(ctx),
		ctx:      ctx,
		svcCtx:   svcCtx,
		userRepo: sysuserrepo.NewSysUserRepo(svcCtx.DB),
	}
}

func (l *RegisterLogic) Register(req *types.RegisterReq) (resp *types.RegisterResp, err error) {

	// 1. 检查用户名是否已存在
	user, err := l.userRepo.GetByUserName(l.ctx, req.UserName)
	if err != nil && !generated.IsNotFound(err) {
		l.Error("GetUser l.userRepo.GetByUserName err: ", err.Error())
		return nil, xerr.NewErrCodeMsg(xerr.DbError, "数据库查询错误")
	}
	if user != nil {
		return nil, xerr.NewErrCodeMsg(xerr.DbRecordExist, "用户名已存在")
	}

	// 检查电话号是否已存在
	user, err = l.userRepo.GetByPhone(l.ctx, req.Phone)
	if err != nil && !generated.IsNotFound(err) {
		l.Error("GetUser l.userRepo.GetByPhone err: ", err.Error())
		return nil, xerr.NewErrCodeMsg(xerr.DbError, "数据库查询错误")
	}
	if user != nil {
		return nil, xerr.NewErrCodeMsg(xerr.DbRecordExist, "手机号已存在")
	}

	// 2. 生成密码盐和加密密码
	salt, err := passwordgen.GenerateSalt()
	if err != nil {
		return nil, xerr.NewErrCodeMsg(xerr.ServerError, "生成密码盐失败")
	}
	hashedPassword, err := passwordgen.Argon2Hash(req.Password, salt)
	if err != nil {
		return nil, xerr.NewErrCodeMsg(xerr.ServerError, "密码加密失败")
	}

	// 3. 生成用户ID
	userID, err := idgen.GenerateUUID()
	if err != nil {
		return nil, xerr.NewErrCodeMsg(xerr.ServerError, "生成用户ID失败")
	}

	// 4. 创建用户
	newUser, err := l.userRepo.Create(l.ctx, &generated.SysUser{
		TenantCode: constants.TenantCodeDefault,
		UserID:     userID,
		Phone:      req.Phone,
		UserName:   req.UserName,
		Name:       req.Name,
		Email:      req.Email,
		Sex:        req.Sex,
		Status:     1,
		PwdHashed:  hashedPassword,
		PwdSalt:    salt,
	})
	if err != nil {
		l.Error("CreateUser Create err:", err.Error())
		return nil, xerr.NewErrCodeMsg(xerr.DbError, "创建用户失败")
	}

	// 5. 返回结果
	return &types.RegisterResp{
		UserID: newUser.UserID,
	}, nil
}
