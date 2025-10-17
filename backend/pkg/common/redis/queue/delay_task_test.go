package queue

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDelayQueue_AddDelayTask(t *testing.T) {
	rdb, cleanup := SetupRedis(t)
	defer cleanup()

	dq := NewDelayQueue(rdb)
	task := &DelayTask{
		ID:    "test-task-1",
		Topic: "test-topic",
		Data:  map[string]interface{}{"key": "value"},
	}

	// 记录当前时间
	now := time.Now()
	delayDuration := 5 * time.Second

	// 测试添加任务
	success, err := dq.AddDelayTask(task, delayDuration)
	assert.NoError(t, err)
	assert.True(t, success)

	// 验证任务是否正确添加到Redis
	tasks, err := rdb.Zrange(DefaultQueueKey, 0, -1)
	assert.NoError(t, err)
	assert.Len(t, tasks, 1)

	// 验证任务内容和执行时间
	var savedTask DelayTask
	err = json.Unmarshal([]byte(tasks[0]), &savedTask)
	assert.NoError(t, err)
	assert.Equal(t, task.ID, savedTask.ID)
	assert.Equal(t, task.Topic, savedTask.Topic)

	// 验证任务的执行时间是否正确设置为5秒后
	expectedExecTime := now.Add(delayDuration).UnixNano() / 1e6
	assert.InDelta(t, expectedExecTime, savedTask.ExecTime, float64(100), "任务执行时间应该设置为5秒后")

	// 启动消费者
	processed := make(chan *DelayTask, 1)
	dq.StartConsumer(func(t *DelayTask) error {
		processed <- t
		return nil
	})

	// 等待任务处理（应该在5秒后）
	select {
	case processedTask := <-processed:
		assert.Equal(t, task.ID, processedTask.ID)
		assert.Equal(t, task.Topic, processedTask.Topic)
		actualDelay := time.Since(now)
		assert.InDelta(t, delayDuration.Seconds(), actualDelay.Seconds(), 1, "任务应该在大约5秒后执行")
	case <-time.After(7 * time.Second):
		t.Fatal("任务没有在预期的延迟时间内执行")
	}

	dq.Stop()
}

func TestDelayQueue_ConsumerLoop(t *testing.T) {
	rdb, cleanup := SetupRedis(t)
	defer cleanup()

	dq := NewDelayQueue(rdb)

	// 添加一个立即执行的任务
	task := &DelayTask{
		ID:       "test-task-2",
		Topic:    "test-topic",
		Data:     "test-data",
		ExecTime: time.Now().UnixNano() / 1e6,
	}

	taskJSON, _ := json.Marshal(task)
	success, err := rdb.Zadd(DefaultQueueKey, task.ExecTime, string(taskJSON))
	assert.NoError(t, err)
	assert.True(t, success)

	// 创建任务处理通道
	processed := make(chan *DelayTask, 1)
	dq.StartConsumer(func(t *DelayTask) error {
		processed <- t
		return nil
	})

	// 等待任务处理
	select {
	case processedTask := <-processed:
		assert.Equal(t, task.ID, processedTask.ID)
		assert.Equal(t, task.Topic, processedTask.Topic)
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for task processing")
	}

	// 验证任务是否从队列中移除
	tasks, err := rdb.Zrange(DefaultQueueKey, 0, -1)
	assert.NoError(t, err)
	assert.Len(t, tasks, 0)

	dq.Stop()
}

func TestDelayQueue_TaskTimeout(t *testing.T) {
	rdb, cleanup := SetupRedis(t)
	defer cleanup()

	dq := NewDelayQueue(rdb)

	// 添加一个会超时的任务
	task := &DelayTask{
		ID:       "test-task-3",
		Topic:    "test-topic",
		Data:     "test-data",
		ExecTime: time.Now().UnixNano() / 1e6,
	}

	taskJSON, _ := json.Marshal(task)
	success, err := rdb.Zadd(DefaultQueueKey, task.ExecTime, string(taskJSON))
	assert.NoError(t, err)
	assert.True(t, success)

	// 启动消费者，使用会超时的处理函数
	dq.StartConsumer(func(t *DelayTask) error {
		time.Sleep(DefaultTimeout + time.Second)
		return nil
	})

	// 等待任务处理和超时
	time.Sleep(DefaultTimeout + 2*time.Second)

	// 验证任务是否被移到死信队列
	deadLetterTasks, err := rdb.Zrange(DefaultDeadLetterKey, 0, -1)
	assert.NoError(t, err)
	assert.Len(t, deadLetterTasks, 1)

	var deadTask DelayTask
	err = json.Unmarshal([]byte(deadLetterTasks[0]), &deadTask)
	assert.NoError(t, err)
	assert.Equal(t, task.ID, deadTask.ID)

	dq.Stop()
}
