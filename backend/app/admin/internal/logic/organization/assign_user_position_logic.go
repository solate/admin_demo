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

type AssignUserPositionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	userPositionRepo *organizationrepo.UserPositionRepo
}

// 分配用户岗位
func NewAssignUserPositionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AssignUserPositionLogic {
	return &AssignUserPositionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		userPositionRepo: organizationrepo.NewUserPositionRepo(svcCtx.DB),
	}
}

func (l *AssignUserPositionLogic) AssignUserPosition(req *types.AssignUserPositionReq) (resp bool, err error) {
	// 1. 参数校验
	if req.UserID == "" {
		return false, xerr.NewErrCodeMsg(xerr.ParamError, "用户ID不能为空")
	}
	if req.PositionID == "" {
		return false, xerr.NewErrCodeMsg(xerr.ParamError, "岗位ID不能为空")
	}

	// 2. 创建用户岗位关联
	newUserPosition := &generated.UserPosition{
		UserID:     req.UserID,
		PositionID: req.PositionID,
	}

	_, err = l.userPositionRepo.Create(l.ctx, newUserPosition)
	if err != nil {
		l.Error("AssignUserPosition Create err:", err.Error())
		return false, xerr.NewErrCodeMsg(xerr.ServerError, "分配用户岗位失败")
	}

	return true, nil
}
