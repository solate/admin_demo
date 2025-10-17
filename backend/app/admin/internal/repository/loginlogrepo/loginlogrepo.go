package loginlogrepo

import (
	"context"
	"net/http"
	"time"

	"admin_backend/pkg/common"
	"admin_backend/pkg/common/contextutil"
	"admin_backend/pkg/ent"
	"admin_backend/pkg/ent/generated"
	"admin_backend/pkg/ent/generated/loginlog"
	"admin_backend/pkg/ent/generated/predicate"
	"admin_backend/pkg/utils/idgen"
	"admin_backend/pkg/utils/useragent"
)

type LoginLogRepo struct {
	db *ent.Client
}

func NewLoginLogRepo(db *ent.Client) *LoginLogRepo {
	return &LoginLogRepo{
		db: db,
	}
}

// Create 创建登录日志
func (r *LoginLogRepo) Create(ctx context.Context, log *generated.LoginLog) (*generated.LoginLog, error) {

	// 创建日志记录
	return r.db.LoginLog.Create().
		SetCreatedAt(time.Now().UnixMilli()).
		SetTenantCode(log.TenantCode).
		SetLogID(log.LogID).
		SetUserID(log.UserID).
		SetUserName(log.UserName).
		SetIP(log.IP).
		SetMessage(log.Message).
		SetBrowser(log.Browser).
		SetOs(log.Os).
		SetUserAgent(log.UserAgent).
		SetDevice(log.Device).
		SetLoginTime(log.LoginTime).
		Save(ctx)
}

// ListLoginLogs 查询登录日志列表
func (r *LoginLogRepo) PageList(ctx context.Context, current, limit int, where []predicate.LoginLog) ([]*generated.LoginLog, int, error) {
	offset := common.Offset(current, limit)

	// 构建查询条件
	query := r.db.LoginLog.Query().Where(where...).Order(generated.Desc(loginlog.FieldCreatedAt))

	// 获取总数
	total, err := query.Count(ctx)
	if err != nil || total == 0 {
		return nil, 0, err
	}

	// 分页查询
	list, err := query.Offset(offset).Limit(limit).All(ctx)
	return list, total, nil
}

// addLoginLog
func (l *LoginLogRepo) AddLoginLog(ctx context.Context, r *http.Request, user *generated.SysUser, message string) error {
	tenantCode := contextutil.GetTenantCodeFromCtx(ctx)

	id, err := idgen.GenerateUUID()
	if err != nil {
		return err
	}

	// 获取客户端信息
	clientInfo := useragent.GetClientInfo(r)
	log := &generated.LoginLog{
		TenantCode: tenantCode,
		LogID:      id,
		UserID:     user.UserID,
		UserName:   user.UserName,
		IP:         clientInfo.IP,
		Message:    message,
		UserAgent:  clientInfo.UserAgent,
		Browser:    clientInfo.Browser + " " + clientInfo.BrowserVer,
		Os:         clientInfo.OS,
		Device:     clientInfo.Device,
		LoginTime:  time.Now().UnixMilli(),
	}

	// 创建登录日志记录
	_, err = l.Create(ctx, log)
	return err
}
