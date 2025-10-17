package role

import (
	"context"

	"admin_backend/app/admin/internal/repository/rolerepo"
	"admin_backend/app/admin/internal/svc"
	"admin_backend/app/admin/internal/types"
	"admin_backend/pkg/common/xerr"
	"admin_backend/pkg/ent/generated"
	"admin_backend/pkg/utils/idgen"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateRoleLogic struct {
	logx.Logger
	ctx      context.Context
	svcCtx   *svc.ServiceContext
	roleRepo *rolerepo.RoleRepo
}

// 创建角色
func NewCreateRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateRoleLogic {
	return &CreateRoleLogic{
		Logger:   logx.WithContext(ctx),
		ctx:      ctx,
		svcCtx:   svcCtx,
		roleRepo: rolerepo.NewRoleRepo(svcCtx.DB),
	}
}

func (l *CreateRoleLogic) CreateRole(req *types.CreateRoleReq) (resp *types.CreateRoleResp, err error) {
	// 1. 检查角色编码是否已存在
	role, err := l.roleRepo.GetByCode(l.ctx, req.Code)
	if err != nil && !generated.IsNotFound(err) {
		l.Error("GetRole l.roleRepo.GetByCode err: ", err.Error())
		return nil, xerr.NewErrCodeMsg(xerr.DbError, "数据库查询错误")
	}
	if role != nil {
		return nil, xerr.NewErrCodeMsg(xerr.DbRecordExist, "角色编码已存在")
	}

	// 2. 生成角色ID
	roleID, err := idgen.GenerateUUID()
	if err != nil {
		l.Error("CreateRole GenerateUUID err:", err.Error())
		return nil, xerr.NewErrCodeMsg(xerr.ServerError, "生成角色ID失败")
	}

	// 3. 创建角色
	newRole := &generated.Role{
		RoleID:      roleID,
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		Status:      req.Status,
		Sort:        req.Sort,
	}
	role, err = l.roleRepo.Create(l.ctx, newRole)
	if err != nil {
		l.Error("CreateRole Create err:", err.Error())
		return nil, xerr.NewErrCodeMsg(xerr.ServerError, "创建角色失败")
	}

	// 4. 返回结果
	return &types.CreateRoleResp{
		RoleID: newRole.RoleID,
	}, nil
}
