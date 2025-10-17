package cache

import (
	"fmt"
	"strconv"
	"time"

	"github.com/zeromicro/go-zero/core/stores/redis"
)

// DownloadLimitManager 下载限制管理器
type DownloadLimitManager struct {
	redis *redis.Redis
}

// DownloadLimits 下载限制配置
type DownloadLimits struct {
	DailyPerResource  int // 每份资料每日最大下载次数
	HourlyPerResource int // 每份资料每小时最大下载次数
	DailyTotalLimit   int // 每日总下载资料份数限制
}

// DefaultDownloadLimits 默认下载限制
var DefaultDownloadLimits = DownloadLimits{
	DailyPerResource:  20,  // 每份资料每个自然日最多下载20次
	HourlyPerResource: 10,  // 每份资料每个自然时最多下载10次
	DailyTotalLimit:   100, // 每天下载超过100份资料需要身份验证
}

// CheckResult 检查结果
type CheckResult struct {
	Allowed       bool      // 是否允许下载
	ReasonCode    string    // 限制原因代码
	ReasonMsg     string    // 限制原因描述
	NextAllowTime time.Time // 下次允许下载的时间
}

// NewDownloadLimitManager 创建下载限制管理器
func NewDownloadLimitManager(redis *redis.Redis) *DownloadLimitManager {
	return &DownloadLimitManager{
		redis: redis,
	}
}

// CheckDownloadLimit 检查下载限制
func (d *DownloadLimitManager) CheckDownloadLimit(userID, resourceID string, limits DownloadLimits) (*CheckResult, error) {
	now := time.Now()

	// 1. 检查每份资料每小时下载次数
	hourlyKey := d.getHourlyKey(userID, resourceID, now)
	hourlyCount, err := d.getCounter(hourlyKey)
	if err != nil {
		return nil, err
	}

	if hourlyCount >= limits.HourlyPerResource {
		nextHour := now.Truncate(time.Hour).Add(time.Hour)
		return &CheckResult{
			Allowed:       false,
			ReasonCode:    "HOURLY_LIMIT_EXCEEDED",
			ReasonMsg:     fmt.Sprintf("该资料每小时最多下载%d次，请稍后重试", limits.HourlyPerResource),
			NextAllowTime: nextHour,
		}, nil
	}

	// 2. 检查每份资料每日下载次数
	dailyKey := d.getDailyKey(userID, resourceID, now)
	dailyCount, err := d.getCounter(dailyKey)
	if err != nil {
		return nil, err
	}

	if dailyCount >= limits.DailyPerResource {
		nextDay := now.Truncate(24 * time.Hour).Add(24 * time.Hour)
		return &CheckResult{
			Allowed:       false,
			ReasonCode:    "DAILY_LIMIT_EXCEEDED",
			ReasonMsg:     fmt.Sprintf("该资料每日最多下载%d次，请明天重试", limits.DailyPerResource),
			NextAllowTime: nextDay,
		}, nil
	}

	// 3. 检查每日总下载资料份数
	totalKey := d.getTotalKey(userID, now)
	totalCount, err := d.getTotalResourceCount(totalKey)
	if err != nil {
		return nil, err
	}

	if totalCount >= limits.DailyTotalLimit {
		nextDay := now.Truncate(24 * time.Hour).Add(24 * time.Hour)
		return &CheckResult{
			Allowed:       false,
			ReasonCode:    "DAILY_TOTAL_LIMIT_EXCEEDED",
			ReasonMsg:     fmt.Sprintf("每日最多下载%d份不同资料，请进行身份验证或明天重试", limits.DailyTotalLimit),
			NextAllowTime: nextDay,
		}, nil
	}

	return &CheckResult{
		Allowed:    true,
		ReasonCode: "",
		ReasonMsg:  "",
	}, nil
}

// IncrementDownloadCount 增加下载计数
func (d *DownloadLimitManager) IncrementDownloadCount(userID, resourceID string) error {
	now := time.Now()

	// 1. 增加小时计数器
	hourlyKey := d.getHourlyKey(userID, resourceID, now)
	hourlyExpire := d.getHourlyExpire(now)
	err := d.incrementCounter(hourlyKey, hourlyExpire)
	if err != nil {
		return err
	}

	// 2. 增加日计数器
	dailyKey := d.getDailyKey(userID, resourceID, now)
	dailyExpire := d.getDailyExpire(now)
	err = d.incrementCounter(dailyKey, dailyExpire)
	if err != nil {
		return err
	}

	// 3. 添加到总下载资料集合
	totalKey := d.getTotalKey(userID, now)
	err = d.addToResourceSet(totalKey, resourceID, dailyExpire)
	if err != nil {
		return err
	}

	return nil
}

// getHourlyKey 获取小时计数器键
func (d *DownloadLimitManager) getHourlyKey(userID, resourceID string, t time.Time) string {
	date := t.Format("2006-01-02")
	hour := t.Hour()
	return fmt.Sprintf("download:hourly:%s:%s:%s:%02d", userID, resourceID, date, hour)
}

// getDailyKey 获取日计数器键
func (d *DownloadLimitManager) getDailyKey(userID, resourceID string, t time.Time) string {
	date := t.Format("2006-01-02")
	return fmt.Sprintf("download:daily:%s:%s:%s", userID, resourceID, date)
}

// getTotalKey 获取总计数器键
func (d *DownloadLimitManager) getTotalKey(userID string, t time.Time) string {
	date := t.Format("2006-01-02")
	return fmt.Sprintf("download:total:%s:%s", userID, date)
}

// getHourlyExpire 获取小时过期时间
func (d *DownloadLimitManager) getHourlyExpire(t time.Time) time.Duration {
	nextHour := t.Truncate(time.Hour).Add(time.Hour)
	return time.Until(nextHour) + time.Minute // 多加1分钟确保不会太早过期
}

// getDailyExpire 获取日过期时间
func (d *DownloadLimitManager) getDailyExpire(t time.Time) time.Duration {
	nextDay := t.Truncate(24 * time.Hour).Add(24 * time.Hour)
	return time.Until(nextDay) + time.Minute // 多加1分钟确保不会太早过期
}

// getCounter 获取计数器值
func (d *DownloadLimitManager) getCounter(key string) (int, error) {
	val, err := d.redis.Get(key)
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		}
		return 0, err
	}

	// 处理空字符串的情况
	if val == "" {
		return 0, nil
	}

	count, err := strconv.Atoi(val)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// incrementCounter 增加计数器
func (d *DownloadLimitManager) incrementCounter(key string, expire time.Duration) error {
	// 使用Redis的INCR命令原子性增加计数
	_, err := d.redis.Incr(key)
	if err != nil {
		return err
	}

	// 设置过期时间
	return d.redis.Expire(key, int(expire.Seconds()))
}

// addToResourceSet 添加资料到集合
func (d *DownloadLimitManager) addToResourceSet(key, resourceID string, expire time.Duration) error {
	// 使用Redis的SADD命令添加到集合
	_, err := d.redis.Sadd(key, resourceID)
	if err != nil {
		return err
	}

	// 设置过期时间
	return d.redis.Expire(key, int(expire.Seconds()))
}

// getTotalResourceCount 获取总下载资料数量
func (d *DownloadLimitManager) getTotalResourceCount(key string) (int, error) {
	count, err := d.redis.Scard(key)
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		}
		return 0, err
	}

	return int(count), nil
}

// ResetUserDailyLimit 重置用户每日限制（用于身份验证后）
func (d *DownloadLimitManager) ResetUserDailyLimit(userID string) error {
	now := time.Now()
	totalKey := d.getTotalKey(userID, now)

	// 删除总下载资料集合，允许用户重新下载
	_, err := d.redis.Del(totalKey)
	return err
}

// GetUserDownloadStats 获取用户下载统计信息
func (d *DownloadLimitManager) GetUserDownloadStats(userID, resourceID string) (map[string]interface{}, error) {
	now := time.Now()

	// 获取小时计数
	hourlyKey := d.getHourlyKey(userID, resourceID, now)
	hourlyCount, err := d.getCounter(hourlyKey)
	if err != nil {
		return nil, err
	}

	// 获取日计数
	dailyKey := d.getDailyKey(userID, resourceID, now)
	dailyCount, err := d.getCounter(dailyKey)
	if err != nil {
		return nil, err
	}

	// 获取总资料数
	totalKey := d.getTotalKey(userID, now)
	totalCount, err := d.getTotalResourceCount(totalKey)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"hourly_count": hourlyCount,
		"daily_count":  dailyCount,
		"total_count":  totalCount,
		"limits": map[string]int{
			"hourly_per_resource": DefaultDownloadLimits.HourlyPerResource,
			"daily_per_resource":  DefaultDownloadLimits.DailyPerResource,
			"daily_total_limit":   DefaultDownloadLimits.DailyTotalLimit,
		},
	}, nil
}
