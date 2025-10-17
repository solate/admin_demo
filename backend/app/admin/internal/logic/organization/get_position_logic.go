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

type GetPositionLogic struct {
	logx.Logger
	ctx     context.Context
	svcCtx  *svc.ServiceContext
	orgRepo *organizationrepo.PositionRepo
}

// 获取岗位详情
func NewGetPositionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPositionLogic {
	return &GetPositionLogic{
		Logger:  logx.WithContext(ctx),
		ctx:     ctx,
		svcCtx:  svcCtx,
		orgRepo: organizationrepo.NewPositionRepo(svcCtx.DB),
	}
}

func (l *GetPositionLogic) GetPosition(req *types.GetPositionReq) (resp *types.PositionInfo, err error) {
	// 1. 查询岗位信息
	pos, err := l.orgRepo.GetByPositionID(l.ctx, req.PositionID)
	if err != nil {
		if generated.IsNotFound(err) {
			return nil, xerr.NewErrCodeMsg(xerr.RecordNotFound, "岗位不存在")
		}
		l.Error("GetPosition l.orgRepo.GetByPositionID err: ", err.Error())
		return nil, xerr.NewErrCodeMsg(xerr.DbError, "数据库查询错误")
	}

	// 2. 转换为响应结构
	resp = &types.PositionInfo{
		PositionID:   pos.PositionID,
		Name:         pos.Name,
		DepartmentID: pos.DepartmentID,
		CreatedAt:    pos.CreatedAt,
		UpdatedAt:    pos.UpdatedAt,
	}

	return resp, nil
}
