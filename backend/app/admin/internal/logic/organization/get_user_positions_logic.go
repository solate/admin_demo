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

type GetUserPositionsLogic struct {
	logx.Logger
	ctx              context.Context
	svcCtx           *svc.ServiceContext
	userPositionRepo *organizationrepo.UserPositionRepo
}

// 获取用户的岗位列表
func NewGetUserPositionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserPositionsLogic {
	return &GetUserPositionsLogic{
		Logger:           logx.WithContext(ctx),
		ctx:              ctx,
		svcCtx:           svcCtx,
		userPositionRepo: organizationrepo.NewUserPositionRepo(svcCtx.DB),
	}
}

func (l *GetUserPositionsLogic) GetUserPositions(req *types.GetUserPositionsReq) (resp *types.GetUserPositionsResp, err error) {
	// 1. 查询用户的岗位关联信息
	userPositions, err := l.userPositionRepo.GetUserPositions(l.ctx, req.UserID)
	if err != nil {
		if generated.IsNotFound(err) {
			return nil, xerr.NewErrCodeMsg(xerr.RecordNotFound, "用户没有分配岗位")
		}
		l.Error("GetUserPositions l.userPositionRepo.GetUserPositions err: ", err.Error())
		return nil, xerr.NewErrCodeMsg(xerr.DbError, "数据库查询错误")
	}

	// 2. 转换为响应结构
	positionList := make([]*types.PositionInfo, 0)
	for _, up := range userPositions {
		positionList = append(positionList, &types.PositionInfo{
			PositionID: up.PositionID,
		})
	}

	resp = &types.GetUserPositionsResp{
		List: positionList,
	}

	return resp, nil
}
