package permission

import (
	"context"

	"admin_backend/app/admin/internal/repository/permissionrepo"
	"admin_backend/app/admin/internal/svc"
	"admin_backend/app/admin/internal/types"
	"admin_backend/pkg/common/xerr"
	"admin_backend/pkg/ent/generated"
	"admin_backend/pkg/utils/idgen"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreatePermissionLogic struct {
	logx.Logger
	ctx            context.Context
	svcCtx         *svc.ServiceContext
	permissionRepo *permissionrepo.PermissionRepo
}

// 创建权限规则
func NewCreatePermissionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreatePermissionLogic {
	return &CreatePermissionLogic{
		Logger:         logx.WithContext(ctx),
		ctx:            ctx,
		svcCtx:         svcCtx,
		permissionRepo: permissionrepo.NewPermissionRepo(svcCtx.DB),
	}
}

func (l *CreatePermissionLogic) CreatePermission(req *types.CreatePermissionReq) (resp *types.CreatePermissionResp, err error) {
	// 1. 检查权限编码是否已存在
	permission, err := l.permissionRepo.GetByCode(l.ctx, req.Code)
	if err != nil && !generated.IsNotFound(err) {
		l.Error("GetPermission l.permissionRepo.GetByCode err: ", err.Error())
		return nil, xerr.NewErrCodeMsg(xerr.DbError, "数据库查询错误")
	}
	if permission != nil {
		return nil, xerr.NewErrCodeMsg(xerr.DbRecordExist, "权限编码已存在")
	}

	// 2. 生成权限ID
	permissionID, err := idgen.GenerateUUID()
	if err != nil {
		l.Error("CreatePermission GenerateUUID err:", err.Error())
		return nil, xerr.NewErrCodeMsg(xerr.ServerError, "生成权限ID失败")
	}

	// 3. 创建权限
	newPermission := &generated.Permission{
		PermissionID: permissionID,
		Name:         req.Name,
		Code:         req.Code,
		Type:         req.Type,
		Resource:     req.Resource,
		Action:       req.Action,
		ParentID:     req.ParentID,
		Description:  req.Description,
		Status:       req.Status,
		MenuID:       req.MenuID,
	}
	_, err = l.permissionRepo.Create(l.ctx, newPermission)
	if err != nil {
		l.Error("CreatePermission Create err:", err.Error())
		return nil, xerr.NewErrCodeMsg(xerr.ServerError, "创建权限失败")
	}

	// 4. 返回结果
	return &types.CreatePermissionResp{
		PermissionID: newPermission.PermissionID,
	}, nil
}
