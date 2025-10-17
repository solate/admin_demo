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

type UpdateDictTypeLogic struct {
	logx.Logger
	ctx      context.Context
	svcCtx   *svc.ServiceContext
	dictRepo *dictrepo.DictTypeRepo
}

// 更新字典类型
func NewUpdateDictTypeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateDictTypeLogic {
	return &UpdateDictTypeLogic{
		Logger:   logx.WithContext(ctx),
		ctx:      ctx,
		svcCtx:   svcCtx,
		dictRepo: dictrepo.NewDictTypeRepo(svcCtx.DB),
	}
}

func (l *UpdateDictTypeLogic) UpdateDictType(req *types.UpdateDictTypeReq) (resp bool, err error) {
	// 1. 检查字典类型是否存在
	dictType, err := l.dictRepo.GetByTypeID(l.ctx, req.TypeID)
	if err != nil {
		if generated.IsNotFound(err) {
			return false, xerr.NewErrCodeMsg(xerr.DbError, "字典类型不存在")
		}
		l.Error("GetDictType l.dictRepo.GetByTypeID err: ", err.Error())
		return false, xerr.NewErrCodeMsg(xerr.DbError, "数据库查询错误")
	}

	// 2. 更新字典类型
	dictType.Name = req.Name
	dictType.Description = req.Description
	dictType.Status = req.Status

	_, err = l.dictRepo.Update(l.ctx, dictType)
	if err != nil {
		l.Error("UpdateDictType l.dictRepo.UpdateType err:", err.Error())
		return false, xerr.NewErrCodeMsg(xerr.ServerError, "更新字典类型失败")
	}

	return true, nil
}
