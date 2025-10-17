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

type DeletePositionLogic struct {
	logx.Logger
	ctx     context.Context
	svcCtx  *svc.ServiceContext
	orgRepo *organizationrepo.PositionRepo
}

// 删除岗位
func NewDeletePositionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeletePositionLogic {
	return &DeletePositionLogic{
		Logger:  logx.WithContext(ctx),
		ctx:     ctx,
		svcCtx:  svcCtx,
		orgRepo: organizationrepo.NewPositionRepo(svcCtx.DB),
	}
}

func (l *DeletePositionLogic) DeletePosition(req *types.DeletePositionReq) (resp bool, err error) {
	// 1. 检查岗位是否存在
	pos, err := l.orgRepo.GetByPositionID(l.ctx, req.PositionID)
	if err != nil {
		if generated.IsNotFound(err) {
			return false, xerr.NewErrCodeMsg(xerr.RecordNotFound, "岗位不存在")
		}
		l.Error("GetPosition l.orgRepo.GetByPositionID err: ", err.Error())
		return false, xerr.NewErrCodeMsg(xerr.DbError, "数据库查询错误")
	}

	// 2. 执行删除
	_, err = l.orgRepo.Delete(l.ctx, pos.PositionID)
	if err != nil {
		l.Error("DeletePosition Delete err:", err.Error())
		return false, xerr.NewErrCodeMsg(xerr.ServerError, "删除岗位失败")
	}

	return true, nil
}
