// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package sys_user

import (
	"context"

	"admin_backend/app/admin/internal/repository/loginlogrepo"
	"admin_backend/app/admin/internal/svc"
	"admin_backend/app/admin/internal/types"
	"admin_backend/pkg/common/xerr"
	"admin_backend/pkg/ent/generated/loginlog"
	"admin_backend/pkg/ent/generated/predicate"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListLoginLogLogic struct {
	logx.Logger
	ctx          context.Context
	svcCtx       *svc.ServiceContext
	loginLogRepo *loginlogrepo.LoginLogRepo
}

// 查询登录记录
func NewListLoginLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLoginLogLogic {
	return &ListLoginLogLogic{
		Logger:       logx.WithContext(ctx),
		ctx:          ctx,
		svcCtx:       svcCtx,
		loginLogRepo: loginlogrepo.NewLoginLogRepo(svcCtx.DB),
	}
}

func (l *ListLoginLogLogic) ListLoginLog(req *types.LoginLogListReq) (resp *types.LoginLogListResp, err error) {
	// 1. 构建查询条件
	var where []predicate.LoginLog
	if req.UserName != "" {
		where = append(where, loginlog.UserNameContains(req.UserName))
	}
	if req.IP != "" {
		where = append(where, loginlog.IPContains(req.IP))
	}
	if req.Status != 0 {
		// Note: LoginLog 表中没有 Status 字段，这里暂时注释
		// where = append(where, loginlog.Status(req.Status))
	}
	// TODO: 处理时间范围查询
	// if req.StartTime != "" {
	// 	startTime, _ := time.Parse("2006-01-02 15:04:05", req.StartTime)
	// 	where = append(where, loginlog.LoginTimeGTE(startTime.UnixMilli()))
	// }
	// if req.EndTime != "" {
	// 	endTime, _ := time.Parse("2006-01-02 15:04:05", req.EndTime)
	// 	where = append(where, loginlog.LoginTimeLTE(endTime.UnixMilli()))
	// }

	// 2. 分页查询
	list, total, err := l.loginLogRepo.PageList(l.ctx, req.Current, req.PageSize, where)
	if err != nil {
		l.Error("ListLoginLog PageList err: ", err.Error())
		return nil, xerr.NewErrCodeMsg(xerr.DbError, "数据库查询错误")
	}

	// 3. 构建返回数据
	logList := make([]*types.LoginLogInfo, 0, len(list))
	for _, log := range list {
		logList = append(logList, &types.LoginLogInfo{
			LogID:     log.LogID,
			UserID:    log.UserID,
			UserName:  log.UserName,
			IP:        log.IP,
			UserAgent: log.UserAgent,
			Status:    0, // TODO: 添加状态字段
			Message:   log.Message,
			CreatedAt: log.CreatedAt,
		})
	}

	return &types.LoginLogListResp{
		Page: &types.PageResponse{
			Total:           total,
			PageSize:        len(list),
			Current:         req.Current,
			RequestPageSize: req.PageSize,
		},
		List: logList,
	}, nil
}
