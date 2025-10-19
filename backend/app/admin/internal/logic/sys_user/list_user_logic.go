// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package sys_user

import (
	"context"

	"admin_backend/app/admin/internal/repository/sysuserrepo"
	"admin_backend/app/admin/internal/svc"
	"admin_backend/app/admin/internal/types"
	"admin_backend/pkg/common/xerr"
	"admin_backend/pkg/ent/generated/predicate"
	"admin_backend/pkg/ent/generated/sysuser"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListUserLogic struct {
	logx.Logger
	ctx      context.Context
	svcCtx   *svc.ServiceContext
	userRepo *sysuserrepo.SysUserRepo
}

// 获取用户列表
func NewListUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListUserLogic {
	return &ListUserLogic{
		Logger:   logx.WithContext(ctx),
		ctx:      ctx,
		svcCtx:   svcCtx,
		userRepo: sysuserrepo.NewSysUserRepo(svcCtx.DB),
	}
}

func (l *ListUserLogic) ListUser(req *types.UserListReq) (resp *types.UserListResp, err error) {
	// 1. 构建查询条件
	var where []predicate.SysUser
	if req.Name != "" {
		where = append(where, sysuser.NameContains(req.Name))
	}
	if req.Phone != "" {
		where = append(where, sysuser.PhoneContains(req.Phone))
	}
	if req.Status != 0 {
		where = append(where, sysuser.Status(req.Status))
	}

	// 2. 分页查询
	list, total, err := l.userRepo.PageList(l.ctx, req.Current, req.PageSize, where)
	if err != nil {
		l.Error("ListUser PageList err: ", err.Error())
		return nil, xerr.NewErrCodeMsg(xerr.DbError, "数据库查询错误")
	}

	// 3. 构建返回数据
	userList := make([]*types.SysUserInfo, 0, len(list))
	for _, user := range list {
		userList = append(userList, &types.SysUserInfo{
			UserID:    user.UserID,
			UserName:  user.UserName,
			Name:      user.Name,
			Phone:     user.Phone,
			Email:     user.Email,
			Avatar:    user.Avatar,
			Status:    user.Status,
			CreatedAt: user.CreatedAt,
			RoleList:  []*types.RoleListInfo{}, // TODO: 查询用户角色列表
		})
	}

	return &types.UserListResp{
		Page: &types.PageResponse{
			Total:           total,
			PageSize:        len(list),
			Current:         req.Current,
			RequestPageSize: req.PageSize,
		},
		List: userList,
	}, nil
}
