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

type UpdateDictItemLogic struct {
	logx.Logger
	ctx      context.Context
	svcCtx   *svc.ServiceContext
	dictRepo *dictrepo.DictItemRepo
}

// 更新字典数据
func NewUpdateDictItemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateDictItemLogic {
	return &UpdateDictItemLogic{
		Logger:   logx.WithContext(ctx),
		ctx:      ctx,
		svcCtx:   svcCtx,
		dictRepo: dictrepo.NewDictItemRepo(svcCtx.DB),
	}
}

func (l *UpdateDictItemLogic) UpdateDictItem(req *types.UpdateDictItemReq) (resp bool, err error) {
	// 1. 检查字典数据项是否存在
	dictItem, err := l.dictRepo.GetByItemID(l.ctx, req.ItemID)
	if err != nil {
		if generated.IsNotFound(err) {
			return false, xerr.NewErrCodeMsg(xerr.DbError, "字典数据项不存在")
		}
		l.Error("GetDictItem l.dictRepo.GetByItemID err: ", err.Error())
		return false, xerr.NewErrCodeMsg(xerr.DbError, "数据库查询错误")
	}

	// 2. 更新字典数据项
	dictItem.Label = req.Label
	dictItem.Value = req.Value
	dictItem.Sort = req.Sort
	dictItem.Description = req.Description
	dictItem.Status = req.Status

	_, err = l.dictRepo.Update(l.ctx, dictItem)
	if err != nil {
		l.Error("UpdateDictItem l.dictRepo.Update err:", err.Error())
		return false, xerr.NewErrCodeMsg(xerr.ServerError, "更新字典数据项失败")
	}

	return true, nil
}
