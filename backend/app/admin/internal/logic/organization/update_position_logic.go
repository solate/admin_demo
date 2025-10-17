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

type UpdatePositionLogic struct {
	logx.Logger
	ctx     context.Context
	svcCtx  *svc.ServiceContext
	orgRepo *organizationrepo.PositionRepo
}

// 更新岗位
func NewUpdatePositionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatePositionLogic {
	return &UpdatePositionLogic{
		Logger:  logx.WithContext(ctx),
		ctx:     ctx,
		svcCtx:  svcCtx,
		orgRepo: organizationrepo.NewPositionRepo(svcCtx.DB),
	}
}

func (l *UpdatePositionLogic) UpdatePosition(req *types.UpdatePositionReq) (resp bool, err error) {
	// 1. 检查岗位是否存在
	pos, err := l.orgRepo.GetByPositionID(l.ctx, req.PositionID)
	if err != nil {
		if generated.IsNotFound(err) {
			return false, xerr.NewErrCodeMsg(xerr.RecordNotFound, "岗位不存在")
		}
		l.Error("GetPosition l.orgRepo.GetByPositionID err: ", err.Error())
		return false, xerr.NewErrCodeMsg(xerr.DbError, "数据库查询错误")
	}

	// 2. 如果要更新岗位名称，检查新名称是否已存在
	if req.Name != "" && req.Name != pos.Name {
		existPos, err := l.orgRepo.GetByName(l.ctx, req.Name)
		if err != nil && !generated.IsNotFound(err) {
			l.Error("GetPosition l.orgRepo.GetByName err: ", err.Error())
			return false, xerr.NewErrCodeMsg(xerr.DbError, "数据库查询错误")
		}
		if existPos != nil {
			return false, xerr.NewErrCodeMsg(xerr.DbRecordExist, "岗位名称已存在")
		}
		pos.Name = req.Name
	}

	// 3. 更新其他字段
	if req.DepartmentID != "" {
		pos.DepartmentID = req.DepartmentID
	}

	// 4. 执行更新
	_, err = l.orgRepo.Update(l.ctx, pos)
	if err != nil {
		l.Error("UpdatePosition Update err:", err.Error())
		return false, xerr.NewErrCodeMsg(xerr.ServerError, "更新岗位失败")
	}

	return true, nil
}
