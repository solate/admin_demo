package queue

import (
	"encoding/json"
	"time"

	"github.com/zeromicro/go-zero/core/stores/redis"
)

// DeadLetterQueue 死信队列管理器
type DeadLetterQueue struct {
	rdb        *redis.Redis
	queueKey   string        // 主队列key
	deadKey    string        // 死信队列key
	maxRetries int           // 最大重试次数
	retryDelay time.Duration // 重试延迟时间
}

// NewDeadLetterQueue 创建死信队列管理器
func NewDeadLetterQueue(rdb *redis.Redis, queueKey, deadKey string, maxRetries int, retryDelay time.Duration) *DeadLetterQueue {
	return &DeadLetterQueue{
		rdb:        rdb,
		queueKey:   queueKey,
		deadKey:    deadKey,
		maxRetries: maxRetries,
		retryDelay: retryDelay,
	}
}

// AddToDeadLetter 添加任务到死信队列
func (dlq *DeadLetterQueue) AddToDeadLetter(task *DelayTask) error {
	task.Retries++
	taskJSON, err := json.Marshal(task)
	if err != nil {
		return err
	}

	// 使用任务的执行时间作为分数，方便后续按时间清理过期任务
	_, err = dlq.rdb.Zadd(dlq.deadKey, task.ExecTime, string(taskJSON))
	return err
}

// RetryDeadLetters 重试死信队列中的任务
func (dlq *DeadLetterQueue) RetryDeadLetters() {
	// 获取所有死信任务
	tasks, err := dlq.rdb.Zrange(dlq.deadKey, 0, -1)
	if err != nil {
		return
	}

	// Lua脚本：原子性地从死信队列移除任务并添加到主队列
	luaScript := `
		local taskJSON = ARGV[1]
		local newTaskJSON = ARGV[2]
		local execTime = tonumber(ARGV[3])
		local deadKey = KEYS[1]
		local queueKey = KEYS[2]

		-- 从死信队列中移除任务
		local removed = redis.call('ZREM', deadKey, taskJSON)
		if removed == 1 then
			-- 添加到主队列
			redis.call('ZADD', queueKey, execTime, newTaskJSON)
			return 1
		end
		return 0
	`

	for _, taskJSON := range tasks {
		var task DelayTask
		if err := json.Unmarshal([]byte(taskJSON), &task); err != nil {
			continue
		}

		// 检查重试次数
		if task.Retries >= dlq.maxRetries {
			// 超过最大重试次数，从死信队列中删除
			dlq.rdb.Zrem(dlq.deadKey, taskJSON)
			continue
		}

		// 重新计算执行时间并加入主队列
		task.ExecTime = time.Now().Add(dlq.retryDelay).UnixMilli()
		newTaskJSON, _ := json.Marshal(task)

		// 使用Lua脚本执行原子操作
		_, err = dlq.rdb.Eval(luaScript, []string{dlq.deadKey, dlq.queueKey}, taskJSON, string(newTaskJSON), task.ExecTime)
		if err != nil {
			continue
		}
	}
}

// CleanupDeadLetters 清理过期的死信任务
func (dlq *DeadLetterQueue) CleanupDeadLetters(maxAge time.Duration) {
	// 计算过期时间点（毫秒）
	expireTime := time.Now().Add(-maxAge).UnixMilli()

	// 删除过期的任务
	dlq.rdb.Zremrangebyscore(dlq.deadKey, 0, expireTime)
}
