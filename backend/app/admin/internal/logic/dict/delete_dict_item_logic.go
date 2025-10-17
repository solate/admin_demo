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

type DeleteDictItemLogic struct {
	logx.Logger
	ctx      context.Context
	svcCtx   *svc.ServiceContext
	dictRepo *dictrepo.DictItemRepo
}

// 删除字典数据
func NewDeleteDictItemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteDictItemLogic {
	return &DeleteDictItemLogic{
		Logger:   logx.WithContext(ctx),
		ctx:      ctx,
		svcCtx:   svcCtx,
		dictRepo: dictrepo.NewDictItemRepo(svcCtx.DB),
	}
}

func (l *DeleteDictItemLogic) DeleteDictItem(req *types.DeleteDictItemReq) (resp bool, err error) {
	// 1. 检查字典数据项是否存在
	dictItem, err := l.dictRepo.GetByItemID(l.ctx, req.ItemID)
	if err != nil {
		l.Error("DeleteDictItem dictRepo.GetByItemID err:", err.Error())
		if generated.IsNotFound(err) {
			return false, xerr.NewErrMsg("字典数据项不存在")
		}
		return false, xerr.NewErrCodeMsg(xerr.DbError, "查询字典数据项失败")
	}

	// 2. 软删除字典数据项
	_, err = l.dictRepo.Delete(l.ctx, dictItem.ItemID)
	if err != nil {
		l.Error("DeleteDictItem dictRepo.Delete err:", err.Error())
		return false, xerr.NewErrCodeMsg(xerr.DbError, "删除字典数据项失败")
	}

	return true, nil
}
