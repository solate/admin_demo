package cache

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDownloadLimitManager_CheckDownloadLimit(t *testing.T) {
	// 使用测试Redis连接
	dlm := NewDownloadLimitManager(testRedis)
	userID := "test_user_123"
	resourceID := "test_resource_456"

	// 测试限制配置
	limits := DownloadLimits{
		DailyPerResource:  3, // 每日3次
		HourlyPerResource: 2, // 每小时2次
		DailyTotalLimit:   5, // 每日总共5份资料
	}

	// 清理测试数据
	defer func() {
		now := time.Now()
		hourlyKey := dlm.getHourlyKey(userID, resourceID, now)
		dailyKey := dlm.getDailyKey(userID, resourceID, now)
		totalKey := dlm.getTotalKey(userID, now)
		testRedis.Del(hourlyKey)
		testRedis.Del(dailyKey)
		testRedis.Del(totalKey)
	}()

	// 第一次检查 - 应该允许
	result, err := dlm.CheckDownloadLimit(userID, resourceID, limits)
	assert.NoError(t, err)
	assert.True(t, result.Allowed)

	// 增加第一次下载计数
	err = dlm.IncrementDownloadCount(userID, resourceID)
	assert.NoError(t, err)

	// 第二次检查 - 应该允许
	result, err = dlm.CheckDownloadLimit(userID, resourceID, limits)
	assert.NoError(t, err)
	assert.True(t, result.Allowed)

	// 增加第二次下载计数
	err = dlm.IncrementDownloadCount(userID, resourceID)
	assert.NoError(t, err)

	// 第三次检查 - 应该被小时限制阻止
	result, err = dlm.CheckDownloadLimit(userID, resourceID, limits)
	assert.NoError(t, err)
	assert.False(t, result.Allowed)
	assert.Equal(t, "HOURLY_LIMIT_EXCEEDED", result.ReasonCode)
}

func TestDownloadLimitManager_DailyLimit(t *testing.T) {
	dlm := NewDownloadLimitManager(testRedis)
	userID := "test_user_daily"
	resourceID := "test_resource_daily"

	limits := DownloadLimits{
		DailyPerResource:  2,  // 每日2次
		HourlyPerResource: 10, // 每小时10次（不会触发）
		DailyTotalLimit:   10, // 每日总共10份资料（不会触发）
	}

	// 清理测试数据
	defer func() {
		now := time.Now()
		hourlyKey := dlm.getHourlyKey(userID, resourceID, now)
		dailyKey := dlm.getDailyKey(userID, resourceID, now)
		totalKey := dlm.getTotalKey(userID, now)
		testRedis.Del(hourlyKey)
		testRedis.Del(dailyKey)
		testRedis.Del(totalKey)
	}()

	// 模拟跨小时的下载（设置上个小时的计数，这样不会触发小时限制）
	lastHour := time.Now().Add(-time.Hour)
	lastHourKey := dlm.getHourlyKey(userID, resourceID, lastHour)
	err := dlm.incrementCounter(lastHourKey, time.Hour)
	assert.NoError(t, err)

	// 增加当前日的计数
	err = dlm.IncrementDownloadCount(userID, resourceID)
	assert.NoError(t, err)
	err = dlm.IncrementDownloadCount(userID, resourceID)
	assert.NoError(t, err)

	// 现在应该被日限制阻止
	result, err := dlm.CheckDownloadLimit(userID, resourceID, limits)
	assert.NoError(t, err)
	assert.False(t, result.Allowed)
	assert.Equal(t, "DAILY_LIMIT_EXCEEDED", result.ReasonCode)
}

func TestDownloadLimitManager_TotalLimit(t *testing.T) {
	dlm := NewDownloadLimitManager(testRedis)
	userID := "test_user_total"

	limits := DownloadLimits{
		DailyPerResource:  10, // 每日10次（不会触发）
		HourlyPerResource: 10, // 每小时10次（不会触发）
		DailyTotalLimit:   2,  // 每日总共2份资料
	}

	// 清理测试数据
	defer func() {
		now := time.Now()
		totalKey := dlm.getTotalKey(userID, now)
		testRedis.Del(totalKey)

		// 清理所有测试资料
		for i := 1; i <= 3; i++ {
			resourceID := fmt.Sprintf("resource_%d", i)
			hourlyKey := dlm.getHourlyKey(userID, resourceID, now)
			dailyKey := dlm.getDailyKey(userID, resourceID, now)
			testRedis.Del(hourlyKey)
			testRedis.Del(dailyKey)
		}
	}()

	// 下载第一份资料
	err := dlm.IncrementDownloadCount(userID, "resource_1")
	assert.NoError(t, err)

	// 下载第二份资料
	err = dlm.IncrementDownloadCount(userID, "resource_2")
	assert.NoError(t, err)

	// 检查第三份资料 - 应该被总限制阻止
	result, err := dlm.CheckDownloadLimit(userID, "resource_3", limits)
	assert.NoError(t, err)
	assert.False(t, result.Allowed)
	assert.Equal(t, "DAILY_TOTAL_LIMIT_EXCEEDED", result.ReasonCode)
}

func TestDownloadLimitManager_GetUserDownloadStats(t *testing.T) {
	dlm := NewDownloadLimitManager(testRedis)
	userID := "test_user_stats"
	resourceID := "test_resource_stats"

	// 清理测试数据
	defer func() {
		now := time.Now()
		hourlyKey := dlm.getHourlyKey(userID, resourceID, now)
		dailyKey := dlm.getDailyKey(userID, resourceID, now)
		totalKey := dlm.getTotalKey(userID, now)
		testRedis.Del(hourlyKey)
		testRedis.Del(dailyKey)
		testRedis.Del(totalKey)
	}()

	// 增加下载计数
	err := dlm.IncrementDownloadCount(userID, resourceID)
	assert.NoError(t, err)

	// 获取统计信息
	stats, err := dlm.GetUserDownloadStats(userID, resourceID)
	assert.NoError(t, err)
	assert.Equal(t, 1, stats["hourly_count"])
	assert.Equal(t, 1, stats["daily_count"])
	assert.Equal(t, 1, stats["total_count"])
}

func TestDownloadLimitManager_ResetUserDailyLimit(t *testing.T) {
	dlm := NewDownloadLimitManager(testRedis)
	userID := "test_user_reset"

	// 清理测试数据
	defer func() {
		now := time.Now()
		totalKey := dlm.getTotalKey(userID, now)
		testRedis.Del(totalKey)
	}()

	// 添加一些下载记录
	err := dlm.IncrementDownloadCount(userID, "resource_1")
	assert.NoError(t, err)
	err = dlm.IncrementDownloadCount(userID, "resource_2")
	assert.NoError(t, err)

	// 重置每日限制
	err = dlm.ResetUserDailyLimit(userID)
	assert.NoError(t, err)

	// 验证总计数被重置
	now := time.Now()
	totalKey := dlm.getTotalKey(userID, now)
	count, err := dlm.getTotalResourceCount(totalKey)
	assert.NoError(t, err)
	assert.Equal(t, 0, count)
}
