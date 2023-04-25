package redis

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

// Conf Redis相关配置
type Conf struct {
	DB          int      `yaml:"db"`
	Namespace   string   `yaml:"namespace"`
	Type        RdbType  `yaml:"type"`
	MasterName  string   `yaml:"masterName"`
	Addrs       []string `yaml:"addrs"`
	Password    string   `yaml:"password"`
	PoolSize    int      `yaml:"poolSize"`
	MaxIdleConn int      `yaml:"maxIdleConn"`
	Timeout     int      `yaml:"timeout"`
}

// RdbType Redis类型
type RdbType string

const (
	// RdbCluster Redis集群
	RdbCluster RdbType = "cluster"
	// RdbSentinel Redis哨兵
	RdbSentinel RdbType = "sentinel"
)

type rdb interface {
	redis.Cmdable
	Do(ctx context.Context, args ...interface{}) *redis.Cmd
	Close() error
}

type Client struct {
	mutex   sync.Mutex
	rdbType RdbType
	rdb
	namespace string // key 的命名空间
}

// Destroy 销毁数据库客户端
func (cli *Client) Destroy() {
	_ = cli.Close()
}

// WrapKey 使用配置的命名空间包装Key，返回一个包装过的key
func (cli *Client) WrapKey(subKey string) string {
	return fmt.Sprintf("%s:%s", cli.namespace, subKey)
}

// NewRedisClient 初始化Redis连接池
func NewRedisClient(conf Conf) *Client {
	timeout := 3 * time.Second
	if conf.Timeout > 0 {
		timeout = time.Duration(conf.Timeout) * time.Second
	}

	var rdb *Client
	if conf.Type == RdbCluster {
		rdb = &Client{rdb: newClusterClient(conf, timeout), rdbType: conf.Type}
	} else if conf.Type == RdbSentinel {
		rdb = &Client{rdb: newFailoverClient(conf, timeout), rdbType: conf.Type}
	} else {
		rdb = &Client{rdb: newSingleClient(conf, timeout), rdbType: conf.Type}
	}

	if _, err := rdb.Ping(context.TODO()).Result(); err != nil {
		log.Fatalf("redis ping fail: %v", err)
	}
	rdb.namespace = conf.Namespace
	return rdb
}

// newSingleClient 单机模式客户端
func newSingleClient(conf Conf, timeout time.Duration) *redis.Client {
	return redis.NewClient(&redis.Options{
		DB:           conf.DB,
		Addr:         conf.Addrs[0],
		Password:     conf.Password,
		PoolSize:     conf.PoolSize,
		MinIdleConns: conf.MaxIdleConn,
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
	})
}

// newFailoverClient 哨兵模式客户端
func newFailoverClient(conf Conf, timeout time.Duration) *redis.Client {
	return redis.NewFailoverClient(&redis.FailoverOptions{
		DB:            conf.DB,
		MasterName:    conf.MasterName,
		SentinelAddrs: conf.Addrs,
		Password:      conf.Password,
		PoolSize:      conf.PoolSize,
		MinIdleConns:  conf.MaxIdleConn,
		ReadTimeout:   timeout,
		WriteTimeout:  timeout,
	})
}

// newClusterClient 集群模式客户端
func newClusterClient(conf Conf, timeout time.Duration) *redis.ClusterClient {
	return redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:        conf.Addrs,
		Password:     conf.Password,
		PoolSize:     conf.PoolSize,
		MinIdleConns: conf.MaxIdleConn,
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
	})
}
