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

type RemoveUserPositionLogic struct {
	logx.Logger
	ctx              context.Context
	svcCtx           *svc.ServiceContext
	userPositionRepo *organizationrepo.UserPositionRepo
}

// 移除用户岗位
func NewRemoveUserPositionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveUserPositionLogic {
	return &RemoveUserPositionLogic{
		Logger:           logx.WithContext(ctx),
		ctx:              ctx,
		svcCtx:           svcCtx,
		userPositionRepo: organizationrepo.NewUserPositionRepo(svcCtx.DB),
	}
}

func (l *RemoveUserPositionLogic) RemoveUserPosition(req *types.RemoveUserPositionReq) (resp bool, err error) {
	// 1. 参数校验
	if req.UserID == "" {
		return false, xerr.NewErrCodeMsg(xerr.ParamError, "用户ID不能为空")
	}
	if req.PositionID == "" {
		return false, xerr.NewErrCodeMsg(xerr.ParamError, "岗位ID不能为空")
	}

	// 2. 删除用户岗位关联
	_, err = l.userPositionRepo.Delete(l.ctx, req.UserID, req.PositionID)
	if err != nil {
		if generated.IsNotFound(err) {
			return false, xerr.NewErrCodeMsg(xerr.RecordNotFound, "用户岗位关联不存在")
		}
		l.Error("RemoveUserPosition Delete err:", err.Error())
		return false, xerr.NewErrCodeMsg(xerr.ServerError, "移除用户岗位失败")
	}

	return true, nil
}
