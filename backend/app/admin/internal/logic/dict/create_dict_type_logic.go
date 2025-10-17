package dict

import (
	"context"

	"admin_backend/app/admin/internal/repository/dictrepo"
	"admin_backend/app/admin/internal/svc"
	"admin_backend/app/admin/internal/types"
	"admin_backend/pkg/common/contextutil"
	"admin_backend/pkg/common/xerr"
	"admin_backend/pkg/ent/generated"
	"admin_backend/pkg/utils/idgen"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateDictTypeLogic struct {
	logx.Logger
	ctx      context.Context
	svcCtx   *svc.ServiceContext
	dictRepo *dictrepo.DictTypeRepo
}

// 创建字典类型
func NewCreateDictTypeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateDictTypeLogic {
	return &CreateDictTypeLogic{
		Logger:   logx.WithContext(ctx),
		ctx:      ctx,
		svcCtx:   svcCtx,
		dictRepo: dictrepo.NewDictTypeRepo(svcCtx.DB),
	}
}

func (l *CreateDictTypeLogic) CreateDictType(req *types.CreateDictTypeReq) (resp *types.CreateDictTypeResp, err error) {
	// 1. 检查字典类型编码是否已存在
	dictType, err := l.dictRepo.GetByCode(l.ctx, req.Code)
	if err != nil && !generated.IsNotFound(err) {
		l.Error("GetDictType l.dictRepo.GetByCode err: ", err.Error())
		return nil, xerr.NewErrCodeMsg(xerr.DbError, "数据库查询错误")
	}
	if dictType != nil {
		return nil, xerr.NewErrCodeMsg(xerr.DbRecordExist, "字典类型编码已存在")
	}

	// 2. 生成字典类型ID
	typeID, err := idgen.GenerateUUID()
	if err != nil {
		l.Error("CreateDictType GenerateUUID err:", err.Error())
		return nil, xerr.NewErrCodeMsg(xerr.ServerError, "生成字典类型ID失败")
	}

	// 3. 设置默认值
	if req.Status == 0 {
		req.Status = 1 // 默认启用
	}

	// 4. 创建字典类型
	newDictType := &generated.DictType{
		TenantCode:  contextutil.GetTenantCodeFromCtx(l.ctx),
		TypeID:      typeID,
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
		Status:      req.Status,
	}

	_, err = l.dictRepo.Create(l.ctx, newDictType)
	if err != nil {
		l.Error("CreateDictType Create err:", err.Error())
		return nil, xerr.NewErrCodeMsg(xerr.ServerError, "创建字典类型失败")
	}

	// 5. 返回结果
	return &types.CreateDictTypeResp{
		TypeID: newDictType.TypeID,
	}, nil
}
