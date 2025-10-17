package queue

import (
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWorkerPool_Basic(t *testing.T) {
	rdb, cleanup := SetupRedis(t)
	defer cleanup()

	workerNum := 3
	timeout := 1 * time.Second
	deadLetterKey := "test_dead_letter"

	// 创建工作池
	pool := NewWorkerPool(rdb, workerNum, timeout, deadLetterKey)

	// 记录处理的任务数
	var processedCount int32
	pool.Start(func(task *DelayTask) error {
		atomic.AddInt32(&processedCount, 1)
		return nil
	})

	// 添加任务
	taskCount := 5
	for i := 0; i < taskCount; i++ {
		pool.AddTask(&DelayTask{
			ID:    string(rune(i)),
			Topic: "test",
			Data:  i,
		})
	}

	// 等待任务处理完成
	time.Sleep(2 * time.Second)
	pool.Stop()

	// 验证所有任务都被处理
	assert.Equal(t, int32(taskCount), atomic.LoadInt32(&processedCount))
}

func TestWorkerPool_TaskError(t *testing.T) {
	rdb, cleanup := SetupRedis(t)
	defer cleanup()

	workerNum := 1
	timeout := 1 * time.Second
	deadLetterKey := "test_dead_letter"

	// 创建工作池
	pool := NewWorkerPool(rdb, workerNum, timeout, deadLetterKey)

	// 设置任务处理函数返回错误
	pool.Start(func(task *DelayTask) error {
		return errors.New("task processing error")
	})

	// 添加任务
	task := &DelayTask{
		ID:    "error-task",
		Topic: "test",
		Data:  "error data",
	}
	pool.AddTask(task)

	// 等待任务处理
	time.Sleep(2 * time.Second)

	// 验证任务是否被移到死信队列
	deadLetterTasks, err := rdb.Zrange(deadLetterKey, 0, -1)
	assert.NoError(t, err)
	assert.Len(t, deadLetterTasks, 1)

	pool.Stop()
}

func TestWorkerPool_TaskTimeout(t *testing.T) {
	rdb, cleanup := SetupRedis(t)
	defer cleanup()

	workerNum := 1
	timeout := 100 * time.Millisecond
	deadLetterKey := "test_dead_letter"

	// 清理死信队列
	rdb.Del(deadLetterKey)

	// 创建工作池
	pool := NewWorkerPool(rdb, workerNum, timeout, deadLetterKey)

	// 设置任务处理函数超时
	pool.Start(func(task *DelayTask) error {
		time.Sleep(timeout * 2)
		return nil
	})

	// 添加任务
	task := &DelayTask{
		ID:    "timeout-task",
		Topic: "test",
		Data:  "timeout data",
	}
	pool.AddTask(task)

	// 等待任务超时处理
	time.Sleep(timeout * 3)

	// 验证任务是否被移到死信队列
	deadLetterTasks, err := rdb.Zrange(deadLetterKey, 0, -1)
	assert.NoError(t, err)
	assert.Len(t, deadLetterTasks, 1)

	// 清理死信队列
	rdb.Del(deadLetterKey)

	pool.Stop()
}

func TestWorkerPool_Concurrency(t *testing.T) {
	rdb, cleanup := SetupRedis(t)
	defer cleanup()

	workerNum := 5
	timeout := 1 * time.Second
	deadLetterKey := "test_dead_letter"

	// 创建工作池
	pool := NewWorkerPool(rdb, workerNum, timeout, deadLetterKey)

	// 记录并发执行的任务数
	var concurrentTasks int32
	var maxConcurrentTasks int32

	pool.Start(func(task *DelayTask) error {
		// 增加当前执行的任务数
		current := atomic.AddInt32(&concurrentTasks, 1)
		for {
			max := atomic.LoadInt32(&maxConcurrentTasks)
			if current <= max {
				break
			}
			if atomic.CompareAndSwapInt32(&maxConcurrentTasks, max, current) {
				break
			}
		}

		// 模拟任务处理
		time.Sleep(100 * time.Millisecond)

		// 减少当前执行的任务数
		atomic.AddInt32(&concurrentTasks, -1)
		return nil
	})

	// 添加大量任务
	taskCount := 20
	for i := 0; i < taskCount; i++ {
		pool.AddTask(&DelayTask{
			ID:    string(rune(i)),
			Topic: "test",
			Data:  i,
		})
	}

	// 等待所有任务完成
	time.Sleep(3 * time.Second)
	pool.Stop()

	// 验证最大并发数不超过工作池大小
	assert.LessOrEqual(t, atomic.LoadInt32(&maxConcurrentTasks), int32(workerNum))
}
