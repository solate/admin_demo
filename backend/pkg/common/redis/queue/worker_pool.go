package queue

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/zeromicro/go-zero/core/stores/redis"
)

// WorkerPool 工作池配置
type WorkerPool struct {
	workerNum     int                    // 工作协程数量
	taskChan      chan *DelayTask        // 任务通道
	timeout       time.Duration          // 任务执行超时时间
	processFunc   func(*DelayTask) error // 任务处理函数
	deadLetterKey string                 // 死信队列key
	rdb           *redis.Redis           // Redis客户端
	wg            sync.WaitGroup         // 等待所有worker完成
}

// NewWorkerPool 创建新的工作池
func NewWorkerPool(rdb *redis.Redis, workerNum int, timeout time.Duration, deadLetterKey string) *WorkerPool {
	return &WorkerPool{
		workerNum:     workerNum,
		taskChan:      make(chan *DelayTask, workerNum*2), // 任务通道容量为工作协程数的2倍
		timeout:       timeout,
		deadLetterKey: deadLetterKey,
		rdb:           rdb,
	}
}

// Start 启动工作池
func (wp *WorkerPool) Start(processFunc func(*DelayTask) error) {
	wp.processFunc = processFunc

	// 启动指定数量的worker
	for i := 0; i < wp.workerNum; i++ {
		wp.wg.Add(1)
		go wp.worker()
	}
}

// Stop 停止工作池
func (wp *WorkerPool) Stop() {
	close(wp.taskChan)
	wp.wg.Wait()
}

// AddTask 添加任务到工作池
func (wp *WorkerPool) AddTask(task *DelayTask) {
	wp.taskChan <- task
}

// worker 工作协程
func (wp *WorkerPool) worker() {
	defer wp.wg.Done()

	for task := range wp.taskChan {
		// 创建带超时的上下文
		ctx, cancel := context.WithTimeout(context.Background(), wp.timeout)

		// 使用channel来控制任务执行超时
		done := make(chan error, 1)
		go func() {
			done <- wp.processFunc(task)
		}()

		// 等待任务完成或超时
		select {
		case err := <-done:
			if err != nil {
				// 任务执行失败，加入死信队列
				wp.moveToDeadLetter(task)
			}
		case <-ctx.Done():
			// 任务执行超时，加入死信队列
			wp.moveToDeadLetter(task)
		}

		cancel() // 清理上下文
	}
}

// moveToDeadLetter 将任务移动到死信队列
func (wp *WorkerPool) moveToDeadLetter(task *DelayTask) {
	taskJSON, _ := json.Marshal(task)
	wp.rdb.Zadd(wp.deadLetterKey, time.Now().UnixMilli(), string(taskJSON))
}
