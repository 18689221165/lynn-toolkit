package task

import (
	"github.com/18689221165/lynn-toolkit/redis"
	"github.com/hibiken/asynq"
	"go.uber.org/zap"
	"time"
)

type AsynqServerConfig struct {
	Concurrency    int            `yaml:"concurrency"`
	Queues         map[string]int `yaml:"queues"`
	StrictPriority bool           `yaml:"strictPriority"`
}

func redisConnOpe(conf redis.Conf) asynq.RedisClientOpt {
	return asynq.RedisClientOpt{
		Addr:         conf.Addrs[0],
		Password:     conf.Password,
		ReadTimeout:  time.Duration(conf.Timeout) * time.Second,
		WriteTimeout: time.Duration(conf.Timeout) * time.Second,
		DB:           conf.DB,
		PoolSize:     conf.PoolSize,
	}
}

func NewAsynqClient(conf redis.Conf) *asynq.Client {
	return asynq.NewClient(redisConnOpe(conf))
}

func NewAsynqServer(log *zap.SugaredLogger, conf redis.Conf, scfg AsynqServerConfig) *asynq.Server {
	return asynq.NewServer(
		redisConnOpe(conf),
		asynq.Config{
			// Specify how many concurrent workers to use
			Concurrency: scfg.Concurrency,
			// Optionally specify multiple queues with different priority.
			Queues:         scfg.Queues,
			StrictPriority: scfg.StrictPriority,
			Logger:         log,
			LogLevel:       asynq.InfoLevel,
		},
	)
}
