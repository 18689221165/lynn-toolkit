package redis

import (
	"context"
	"time"
)

var ctx = context.TODO()

const unlock = "if redis.call('get', KEYS[1]) == ARGV[1] then redis.call('del', KEYS[1]); return 1; end return 0;"

// Lock 分布式锁加锁
func (cli *Client) Lock(key string, expiration time.Duration) bool {
	cli.mutex.Lock()
	defer cli.mutex.Unlock()

	// 已经被其他goroutine加锁了
	if cli.Get(ctx, key).Val() != "" {
		return false
	}

	// 尝试设置Redis锁
	ok, err := cli.SetNX(ctx, key, "1", expiration).Result()
	if err != nil {
		return false
	}
	return ok
}

// Unlock 分布式锁解锁
func (cli *Client) Unlock(key string) bool {
	i, err := cli.Eval(ctx, unlock, []string{key}, "1").Int()
	if err != nil {
		return false
	}
	return i > 0
}
