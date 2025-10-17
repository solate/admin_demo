package queue

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/zeromicro/go-zero/core/stores/redis"
)

const (
	DefaultQueueKey      = "delay_queue"
	DefaultDeadLetterKey = "delay_queue_dead_letter"
	DefaultWorkerNum     = 3
	DefaultTimeout       = 30 * time.Second
	DefaultMaxRetries    = 3
	DefaultRetryDelay    = 5 * time.Minute
)

// DelayTask 延迟任务
type DelayTask struct {
	ID       string `json:"id"`
	Topic    string `json:"topic"`     // 任务类型（如：订单超时）
	Data     any    `json:"data"`      // 任务数据
	ExecTime int64  `json:"exec_time"` // 执行时间戳（毫秒）
	Retries  int    `json:"retries"`   // 重试次数
}

// DelayQueue 延时队列管理器
type DelayQueue struct {
	rdb        *redis.Redis
	workerPool *WorkerPool
	queueKey   string
	once       sync.Once
}

// NewDelayQueue 创建延时队列管理器
func NewDelayQueue(rdb *redis.Redis) *DelayQueue {
	return &DelayQueue{
		rdb:      rdb,
		queueKey: DefaultQueueKey,
	}
}

// AddDelayTask 添加任务到延时队列
func (dq *DelayQueue) AddDelayTask(task *DelayTask, delay time.Duration) (bool, error) {
	task.ExecTime = time.Now().Add(delay).UnixMilli() // 使用毫秒级时间戳
	taskJSON, _ := json.Marshal(task)
	return dq.rdb.Zadd(dq.queueKey, task.ExecTime, string(taskJSON))
}

// StartConsumer 启动消费者
func (dq *DelayQueue) StartConsumer(processFunc func(*DelayTask) error) {
	dq.once.Do(func() {
		// 初始化工作池
		dq.workerPool = NewWorkerPool(dq.rdb, DefaultWorkerNum, DefaultTimeout, DefaultDeadLetterKey)
		dq.workerPool.Start(processFunc)

		// 启动消费者循环
		go dq.consumeLoop()
	})
}

// consumeLoop 消费循环
func (dq *DelayQueue) consumeLoop() {
	luaScript := `
        local current_time = tonumber(ARGV[1])
        local tasks = redis.call('ZRANGEBYSCORE', KEYS[1], '-inf', current_time)
        if #tasks > 0 then
            redis.call('ZREM', KEYS[1], unpack(tasks))
        end
        return tasks
    `

	for {
		// 获取当前时间戳（毫秒）
		now := time.Now().UnixMilli()

		// 执行 Lua 脚本
		tasks, err := dq.rdb.Eval(luaScript, []string{dq.queueKey}, now)
		if err != nil {
			log.Println("Error fetching tasks:", err)
			continue
		}

		// 将任务分发到工作池
		for _, taskJSON := range tasks.([]any) {
			var task DelayTask
			if err := json.Unmarshal([]byte(taskJSON.(string)), &task); err != nil {
				continue
			}
			dq.workerPool.AddTask(&task)
		}

		time.Sleep(100 * time.Millisecond) // 缩短轮询间隔提高实时性
	}
}

// Stop 停止延时队列
func (dq *DelayQueue) Stop() {
	if dq.workerPool != nil {
		dq.workerPool.Stop()
	}
}
