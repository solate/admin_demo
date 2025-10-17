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

type CreateDictItemLogic struct {
	logx.Logger
	ctx      context.Context
	svcCtx   *svc.ServiceContext
	dictRepo *dictrepo.DictItemRepo
}

// 创建字典数据
func NewCreateDictItemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateDictItemLogic {
	return &CreateDictItemLogic{
		Logger:   logx.WithContext(ctx),
		ctx:      ctx,
		svcCtx:   svcCtx,
		dictRepo: dictrepo.NewDictItemRepo(svcCtx.DB),
	}
}

func (l *CreateDictItemLogic) CreateDictItem(req *types.CreateDictItemReq) (resp *types.CreateDictItemResp, err error) {
	// 1. 检查字典类型是否存在

	// 2. 生成字典数据项ID
	itemID, err := idgen.GenerateUUID()
	if err != nil {
		l.Error("CreateDictItem GenerateUUID err:", err.Error())
		return nil, xerr.NewErrCodeMsg(xerr.ServerError, "生成字典数据项ID失败")
	}

	// 3. 设置默认值
	if req.Status == 0 {
		req.Status = 1 // 默认启用
	}

	// 4. 创建字典数据项
	newDictItem := &generated.DictItem{
		TenantCode:  contextutil.GetTenantCodeFromCtx(l.ctx),
		ItemID:      itemID,
		TypeCode:    req.TypeCode,
		Label:       req.Label,
		Value:       req.Value,
		Sort:        req.Sort,
		Description: req.Description,
		Status:      req.Status,
	}

	_, err = l.dictRepo.Create(l.ctx, newDictItem)
	if err != nil {
		l.Error("CreateDictItem Create err:", err.Error())
		return nil, xerr.NewErrCodeMsg(xerr.ServerError, "创建字典数据项失败")
	}

	// 5. 返回结果
	return &types.CreateDictItemResp{
		ItemID: newDictItem.ItemID,
	}, nil
}
