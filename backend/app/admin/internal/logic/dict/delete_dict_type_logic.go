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

type DeleteDictTypeLogic struct {
	logx.Logger
	ctx      context.Context
	svcCtx   *svc.ServiceContext
	dictRepo *dictrepo.DictTypeRepo
}

// 删除字典类型
func NewDeleteDictTypeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteDictTypeLogic {
	return &DeleteDictTypeLogic{
		Logger:   logx.WithContext(ctx),
		ctx:      ctx,
		svcCtx:   svcCtx,
		dictRepo: dictrepo.NewDictTypeRepo(svcCtx.DB),
	}
}

func (l *DeleteDictTypeLogic) DeleteDictType(req *types.DeleteDictTypeReq) (resp bool, err error) {
	// 1. 检查字典类型是否存在
	dictType, err := l.dictRepo.GetByTypeID(l.ctx, req.TypeID)
	if err != nil {
		l.Error("DeleteDictType dictRepo.GetByTypeID err:", err.Error())
		if generated.IsNotFound(err) {
			return false, xerr.NewErrMsg("字典类型不存在")
		}
		return false, xerr.NewErrCodeMsg(xerr.DbError, "查询字典类型失败")
	}

	// 2. 软删除字典类型
	_, err = l.dictRepo.Delete(l.ctx, dictType.TypeID)
	if err != nil {
		l.Error("DeleteDictType dictRepo.DeleteByTypeID err:", err.Error())
		return false, xerr.NewErrCodeMsg(xerr.DbError, "删除字典类型失败")
	}

	return true, nil
}
