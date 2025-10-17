package dict

import (
	"context"

	"admin_backend/app/admin/internal/repository/dictrepo"
	"admin_backend/app/admin/internal/svc"
	"admin_backend/app/admin/internal/types"
	"admin_backend/pkg/common/xerr"
	"admin_backend/pkg/ent/generated"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDictTypeLogic struct {
	logx.Logger
	ctx      context.Context
	svcCtx   *svc.ServiceContext
	dictRepo *dictrepo.DictTypeRepo
}

// 获取字典类型详情
func NewGetDictTypeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDictTypeLogic {
	return &GetDictTypeLogic{
		Logger:   logx.WithContext(ctx),
		ctx:      ctx,
		svcCtx:   svcCtx,
		dictRepo: dictrepo.NewDictTypeRepo(svcCtx.DB),
	}
}

func (l *GetDictTypeLogic) GetDictType(req *types.GetDictTypeReq) (resp *types.DictTypeInfo, err error) {
	// 1. 获取字典类型
	dictType, err := l.dictRepo.GetByTypeID(l.ctx, req.TypeID)
	if err != nil {
		if generated.IsNotFound(err) {
			return nil, xerr.NewErrCodeMsg(xerr.DbError, "字典类型不存在")
		}
		l.Error("GetDictType l.dictRepo.GetByTypeID err: ", err.Error())
		return nil, xerr.NewErrCodeMsg(xerr.DbError, "数据库查询错误")
	}

	// 2. 返回结果
	return &types.DictTypeInfo{
		TypeID:      dictType.TypeID,
		Name:        dictType.Name,
		Code:        dictType.Code,
		Description: dictType.Description,
		Status:      dictType.Status,
		CreatedAt:   dictType.CreatedAt,
	}, nil
}
