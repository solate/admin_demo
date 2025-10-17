package organization

import (
	"context"

	"admin_backend/app/admin/internal/repository/organizationrepo"
	"admin_backend/app/admin/internal/svc"
	"admin_backend/app/admin/internal/types"
	"admin_backend/pkg/common/contextutil"
	"admin_backend/pkg/common/xerr"
	"admin_backend/pkg/ent/generated"
	"admin_backend/pkg/utils/idgen"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreatePositionLogic struct {
	logx.Logger
	ctx     context.Context
	svcCtx  *svc.ServiceContext
	orgRepo *organizationrepo.PositionRepo
}

// 创建岗位
func NewCreatePositionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreatePositionLogic {
	return &CreatePositionLogic{
		Logger:  logx.WithContext(ctx),
		ctx:     ctx,
		svcCtx:  svcCtx,
		orgRepo: organizationrepo.NewPositionRepo(svcCtx.DB),
	}
}

func (l *CreatePositionLogic) CreatePosition(req *types.CreatePositionReq) (resp *types.CreatePositionResp, err error) {
	// 1. 参数校验
	if req.Name == "" {
		return nil, xerr.NewErrCodeMsg(xerr.ParamError, "岗位名称不能为空")
	}

	// 2. 检查岗位名称是否已存在
	pos, err := l.orgRepo.GetByName(l.ctx, req.Name)
	if err != nil && !generated.IsNotFound(err) {
		l.Error("GetPosition l.orgRepo.GetByName err: ", err.Error())
		return nil, xerr.NewErrCodeMsg(xerr.DbError, "数据库查询错误")
	}
	if pos != nil {
		l.Errorf("岗位名称已存在: %s", req.Name)
		return nil, xerr.NewErrCodeMsg(xerr.DbRecordExist, "岗位名称已存在")
	}

	// 3. 生成岗位ID
	positionID, err := idgen.GenerateUUID()
	if err != nil {
		l.Error("CreatePosition GenerateUUID err:", err.Error())
		return nil, xerr.NewErrCodeMsg(xerr.ServerError, "生成岗位ID失败")
	}

	// 4. 设置默认值

	// 5. 创建岗位
	newPos := &generated.Position{
		TenantCode: contextutil.GetTenantCodeFromCtx(l.ctx),
		PositionID: positionID,
		Name:       req.Name,
	}

	_, err = l.orgRepo.Create(l.ctx, newPos)
	if err != nil {
		l.Error("CreatePosition Create err:", err.Error())
		return nil, xerr.NewErrCodeMsg(xerr.ServerError, "创建岗位失败")
	}

	// 6. 返回结果
	return &types.CreatePositionResp{
		PositionID: newPos.PositionID,
	}, nil
}
