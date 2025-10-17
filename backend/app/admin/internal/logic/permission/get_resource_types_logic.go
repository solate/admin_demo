package permission

import (
	"context"

	"admin_backend/app/admin/internal/svc"
	"admin_backend/app/admin/internal/types"
	"admin_backend/pkg/common/casbin"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetResourceTypesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取资源类型列表
func NewGetResourceTypesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetResourceTypesLogic {
	return &GetResourceTypesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetResourceTypesLogic) GetResourceTypes() (resp *types.GetResourceTypesResp, err error) {

	var list []*types.ResourceTypeInfo
	// 将 casbin.ResourceTypes 转换为 API 响应格式
	for _, rt := range casbin.ResourceTypes {
		list = append(list, &types.ResourceTypeInfo{
			Type:      rt.Type,
			Actions:   rt.Actions,
			DataRules: rt.DataRules,
		})
	}

	resp = &types.GetResourceTypesResp{
		List: list,
	}

	return resp, nil
}
