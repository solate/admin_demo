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

type GetUserLogic struct {
	logx.Logger
	ctx      context.Context
	svcCtx   *svc.ServiceContext
	userRepo *sysuserrepo.SysUserRepo
}

// 获取用户详情
func NewGetUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserLogic {
	return &GetUserLogic{
		Logger:   logx.WithContext(ctx),
		ctx:      ctx,
		svcCtx:   svcCtx,
		userRepo: sysuserrepo.NewSysUserRepo(svcCtx.DB),
	}
}

func (l *GetUserLogic) GetUser(req *types.GetUserReq) (resp *types.SysUserInfo, err error) {
	// 1. 查询用户信息
	user, err := l.userRepo.GetByUserID(l.ctx, req.UserID)
	if err != nil {
		if generated.IsNotFound(err) {
			return nil, xerr.NewErrCodeMsg(xerr.RecordNotFound, "用户不存在")
		}
		l.Error("GetUser GetByUserID err: ", err.Error())
		return nil, xerr.NewErrCodeMsg(xerr.DbError, "数据库查询错误")
	}

	// 2. 构建返回数据
	return &types.SysUserInfo{
		UserID:    user.UserID,
		UserName:  user.UserName,
		Name:      user.Name,
		Phone:     user.Phone,
		Email:     user.Email,
		Avatar:    user.Avatar,
		Status:    user.Status,
		CreatedAt: user.CreatedAt,
		RoleList:  []*types.RoleListInfo{}, // TODO: 查询用户角色列表
	}, nil
}
