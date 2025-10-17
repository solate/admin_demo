package queue

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

func SetupRedis(t *testing.T) (*redis.Redis, func()) {
	rdb := redis.New("127.0.0.1:6379", redis.WithPass("123456"))
	return rdb, func() {
		// 清理所有测试相关的key
		rdb.Del(DefaultQueueKey)
		rdb.Del(DefaultDeadLetterKey)
	}
}

func TestDeadLetterQueue_AddAndRetry(t *testing.T) {
	rdb, cleanup := SetupRedis(t)
	defer cleanup()

	maxRetries := 3
	retryDelay := time.Second

	// 创建死信队列
	dlq := NewDeadLetterQueue(rdb, DefaultQueueKey, DefaultDeadLetterKey, maxRetries, retryDelay)

	// 添加任务到死信队列
	task := &DelayTask{
		ID:       "dead-task",
		Topic:    "test",
		Data:     "test data",
		ExecTime: time.Now().UnixNano() / 1e6,
		Retries:  0,
	}

	err := dlq.AddToDeadLetter(task)
	assert.NoError(t, err)

	// 验证任务是否在死信队列中
	deadTasks, err := rdb.Zrange(DefaultDeadLetterKey, 0, -1)
	assert.NoError(t, err)
	assert.Len(t, deadTasks, 1)

	// 重试死信队列中的任务
	dlq.RetryDeadLetters()

	// 验证任务是否被移回主队列
	queueTasks, err := rdb.Zrange(DefaultQueueKey, 0, -1)
	assert.NoError(t, err)
	assert.Len(t, queueTasks, 1)

	// 验证死信队列是否为空
	deadTasks, err = rdb.Zrange(DefaultDeadLetterKey, 0, -1)
	assert.NoError(t, err)
	assert.Len(t, deadTasks, 0)
}

func TestDeadLetterQueue_MaxRetries(t *testing.T) {
	rdb, cleanup := SetupRedis(t)
	defer cleanup()

	maxRetries := 2
	retryDelay := time.Second

	dlq := NewDeadLetterQueue(rdb, DefaultQueueKey, DefaultDeadLetterKey, maxRetries, retryDelay)

	// 创建一个已达到最大重试次数的任务
	task := &DelayTask{
		ID:       "max-retries-task",
		Topic:    "test",
		Data:     "test data",
		ExecTime: time.Now().UnixNano() / 1e6,
		Retries:  maxRetries,
	}

	err := dlq.AddToDeadLetter(task)
	assert.NoError(t, err)

	// 尝试重试任务
	dlq.RetryDeadLetters()

	// 验证任务是否被从死信队列中删除
	deadTasks, err := rdb.Zrange(DefaultDeadLetterKey, 0, -1)
	assert.NoError(t, err)
	assert.Len(t, deadTasks, 0)

	// 验证任务没有被移回主队列
	queueTasks, err := rdb.Zrange(DefaultQueueKey, 0, -1)
	assert.NoError(t, err)
	assert.Len(t, queueTasks, 0)
}

func TestDeadLetterQueue_CleanupDeadLetters(t *testing.T) {
	rdb, cleanup := SetupRedis(t)
	defer cleanup()

	maxRetries := 3
	retryDelay := time.Second

	dlq := NewDeadLetterQueue(rdb, DefaultQueueKey, DefaultDeadLetterKey, maxRetries, retryDelay)

	// 添加一个过期的任务
	oldTask := &DelayTask{
		ID:       "old-task",
		Topic:    "test",
		Data:     "old data",
		ExecTime: time.Now().Add(-2*time.Hour).UnixNano() / 1e6,
	}

	// 添加一个新的任务
	newTask := &DelayTask{
		ID:       "new-task",
		Topic:    "test",
		Data:     "new data",
		ExecTime: time.Now().UnixNano() / 1e6,
	}

	// 将任务添加到死信队列
	err := dlq.AddToDeadLetter(oldTask)
	assert.NoError(t, err)
	err = dlq.AddToDeadLetter(newTask)
	assert.NoError(t, err)

	// 清理超过1小时的死信任务
	dlq.CleanupDeadLetters(time.Hour)

	// 验证只有新任务保留在死信队列中
	deadTasks, err := rdb.Zrange(DefaultDeadLetterKey, 0, -1)
	assert.NoError(t, err)
	assert.Len(t, deadTasks, 1)
}
