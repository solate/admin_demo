package organization

import (
	"context"

	"admin_backend/app/admin/internal/repository/organizationrepo"
	"admin_backend/app/admin/internal/svc"
	"admin_backend/app/admin/internal/types"
	"admin_backend/pkg/common/xerr"
	"admin_backend/pkg/ent/generated"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPositionUsersLogic struct {
	logx.Logger
	ctx              context.Context
	svcCtx           *svc.ServiceContext
	userPositionRepo *organizationrepo.UserPositionRepo
}

// 获取岗位下的用户列表
func NewGetPositionUsersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPositionUsersLogic {
	return &GetPositionUsersLogic{
		Logger:           logx.WithContext(ctx),
		ctx:              ctx,
		svcCtx:           svcCtx,
		userPositionRepo: organizationrepo.NewUserPositionRepo(svcCtx.DB),
	}
}

func (l *GetPositionUsersLogic) GetPositionUsers(req *types.GetPositionUsersReq) (resp *types.GetPositionUsersResp, err error) {
	// 1. 查询岗位下的用户关联信息
	userPositions, err := l.userPositionRepo.GetPositionUsers(l.ctx, req.PositionID)
	if err != nil {
		if generated.IsNotFound(err) {
			return nil, xerr.NewErrCodeMsg(xerr.RecordNotFound, "岗位下没有用户")
		}
		l.Error("GetPositionUsers l.userPositionRepo.GetByPositionID err: ", err.Error())
		return nil, xerr.NewErrCodeMsg(xerr.DbError, "数据库查询错误")
	}

	// 2. 转换为响应结构
	userList := make([]*types.UserPositionInfo, 0)
	for _, up := range userPositions {
		userList = append(userList, &types.UserPositionInfo{
			UserID: up.UserID,
		})
	}

	resp = &types.GetPositionUsersResp{
		List: userList,
	}

	return resp, nil
}
