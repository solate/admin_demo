package organization

import (
	"context"

	"admin_backend/app/admin/internal/repository/organizationrepo"
	"admin_backend/app/admin/internal/svc"
	"admin_backend/app/admin/internal/types"
	"admin_backend/pkg/common/xerr"
	"admin_backend/pkg/ent/generated/position"
	"admin_backend/pkg/ent/generated/predicate"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListPositionLogic struct {
	logx.Logger
	ctx     context.Context
	svcCtx  *svc.ServiceContext
	orgRepo *organizationrepo.PositionRepo
}

// 获取岗位列表
func NewListPositionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListPositionLogic {
	return &ListPositionLogic{
		Logger:  logx.WithContext(ctx),
		ctx:     ctx,
		svcCtx:  svcCtx,
		orgRepo: organizationrepo.NewPositionRepo(svcCtx.DB),
	}
}

func (l *ListPositionLogic) ListPosition(req *types.PositionListReq) (resp *types.PositionListResp, err error) {
	// 1. 构建查询条件
	var predicates []predicate.Position
	if req.Name != "" {
		predicates = append(predicates, position.NameContains(req.Name))
	}
	if req.DepartmentID != "" {
		predicates = append(predicates, position.DepartmentID(req.DepartmentID))
	}

	// 2. 查询岗位列表
	list, total, err := l.orgRepo.PageList(l.ctx, req.Current, req.PageSize, predicates)
	if err != nil {
		l.Error("ListPosition PageList err:", err.Error())
		return nil, xerr.NewErrCodeMsg(xerr.DbError, "查询岗位列表失败")
	}

	// 3. 转换为响应结构
	posList := make([]*types.PositionInfo, 0)
	for _, pos := range list {
		posList = append(posList, &types.PositionInfo{
			PositionID:   pos.PositionID,
			Name:         pos.Name,
			DepartmentID: pos.DepartmentID,
			CreatedAt:    pos.CreatedAt,
			UpdatedAt:    pos.UpdatedAt,
		})
	}

	// 4. 返回结果
	resp = &types.PositionListResp{
		Page: &types.PageResponse{
			Current:  req.Current,
			PageSize: req.PageSize,
			Total:    total,
		},
		List: posList,
	}

	return resp, nil
}
