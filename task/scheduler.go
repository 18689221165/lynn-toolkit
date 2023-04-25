package task

import (
	"encoding/json"
	"github.com/hibiken/asynq"
	"go.uber.org/zap"
	"time"
)

type TaskType string // 异步任务类型

type TaskQueue string // 队列名

const (
	AQ_CRITICAL TaskQueue = "critical" // 优先级最高
	AQ_DEFAULT  TaskQueue = "default"  // 默认
	AQ_LOW      TaskQueue = "low"      // 最低的
)

type AsyncTask struct {
	TaskType  TaskType      // 任务类型
	DelaySec  time.Duration // 延时多少秒开始执行
	ProcessAt time.Time     // 指定时间执行
	MaxRetry  int           // 重试次数
	Payload   interface{}   // 作务载体
	TaskId    string        // 任务唯一ID
	UniqueTTL time.Duration // 保持任务唯一的时间段
	QueueName TaskQueue     // 队列名，优先级不同
}

// AsynqTaskScheduler https://github.com/hibiken/asynq/wiki/Getting-Started
type AsynqTaskScheduler struct {
	log    *zap.SugaredLogger
	client *asynq.Client
}

func NewAsynqTaskScheduler(log *zap.SugaredLogger, client *asynq.Client) *AsynqTaskScheduler {
	return &AsynqTaskScheduler{log, client}
}

func (scheduler AsynqTaskScheduler) Schedule(t *AsyncTask) error {
	payload, _ := json.Marshal(t.Payload)

	var options []asynq.Option
	if t.DelaySec > 0 {
		options = append(options, asynq.ProcessIn(t.DelaySec))
	}
	if !t.ProcessAt.IsZero() {
		options = append(options, asynq.ProcessAt(t.ProcessAt))
	}
	if t.MaxRetry > 0 {
		options = append(options, asynq.MaxRetry(t.MaxRetry))
	}
	if t.QueueName != "" {
		options = append(options, asynq.Queue(string(t.QueueName)))
	}
	if t.TaskId != "" {
		options = append(options, asynq.TaskID(t.TaskId))
	}
	if t.TaskId == "" && t.UniqueTTL > 0 {
		options = append(options, asynq.Unique(t.UniqueTTL))
	}
	task := asynq.NewTask(string(t.TaskType), payload, options...)

	_, err := scheduler.client.Enqueue(task)
	if err != nil {
		return err

	}
	scheduler.log.Infof("<<< Asynq任务入队成功,队列:%s,消息:%+v", t.TaskType, t)
	return nil
}
