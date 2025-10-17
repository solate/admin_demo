package role

import (
	"context"

	"admin_backend/app/admin/internal/repository/rolerepo"
	"admin_backend/app/admin/internal/svc"
	"admin_backend/app/admin/internal/types"
	"admin_backend/pkg/common/xerr"
	"admin_backend/pkg/ent/generated"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateRoleLogic struct {
	logx.Logger
	ctx      context.Context
	svcCtx   *svc.ServiceContext
	roleRepo *rolerepo.RoleRepo
}

// 更新角色
func NewUpdateRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateRoleLogic {
	return &UpdateRoleLogic{
		Logger:   logx.WithContext(ctx),
		ctx:      ctx,
		svcCtx:   svcCtx,
		roleRepo: rolerepo.NewRoleRepo(svcCtx.DB),
	}
}

func (l *UpdateRoleLogic) UpdateRole(req *types.UpdateRoleReq) (resp bool, err error) {
	// 1. 检查角色是否存在
	role, err := l.roleRepo.GetByRoleID(l.ctx, req.RoleID)
	if err != nil {
		if generated.IsNotFound(err) {
			return false, xerr.NewErrMsg("角色不存在")
		}
		l.Error("UpdateRole GetByRoleID err:", err.Error())
		return false, xerr.NewErrCodeMsg(xerr.DbError, "查询角色失败")
	}

	// 3. 更新角色信息
	role.Name = req.Name
	role.Description = req.Description
	role.Status = req.Status
	role.Sort = req.Sort

	_, err = l.roleRepo.Update(l.ctx, role)
	if err != nil {
		l.Error("UpdateRole Update err:", err.Error())
		return false, xerr.NewErrMsg("更新角色失败")
	}

	return true, nil
}
