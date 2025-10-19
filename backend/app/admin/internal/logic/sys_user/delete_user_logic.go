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

type DeleteUserLogic struct {
	logx.Logger
	ctx      context.Context
	svcCtx   *svc.ServiceContext
	userRepo *sysuserrepo.SysUserRepo
}

// 删除用户
func NewDeleteUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteUserLogic {
	return &DeleteUserLogic{
		Logger:   logx.WithContext(ctx),
		ctx:      ctx,
		svcCtx:   svcCtx,
		userRepo: sysuserrepo.NewSysUserRepo(svcCtx.DB),
	}
}

func (l *DeleteUserLogic) DeleteUser(req *types.DeleteUserReq) (resp bool, err error) {
	// 1. 检查用户是否存在
	_, err = l.userRepo.GetByUserID(l.ctx, req.UserID)
	if err != nil {
		if generated.IsNotFound(err) {
			return false, xerr.NewErrCodeMsg(xerr.RecordNotFound, "用户不存在")
		}
		l.Error("DeleteUser GetByUserID err: ", err.Error())
		return false, xerr.NewErrCodeMsg(xerr.DbError, "数据库查询错误")
	}

	// 2. 软删除用户
	_, err = l.userRepo.Delete(l.ctx, req.UserID)
	if err != nil {
		l.Error("DeleteUser Delete err:", err.Error())
		return false, xerr.NewErrCodeMsg(xerr.DbError, "删除用户失败")
	}

	return true, nil
}
