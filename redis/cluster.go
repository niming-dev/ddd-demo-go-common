package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/niming-dev/ddd-demo/go-common/log"
)

type ClusterClient interface {
	redis.Cmdable
	Subscribe(context.Context, ...string) *redis.PubSub
}

type hook struct {
}

func (a *hook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	return ctx, nil
}
func (a *hook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	logCmd(ctx, cmd)
	return nil
}
func (a *hook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	return ctx, nil
}
func (a *hook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	return nil
}

func logCmd(ctx context.Context, cmd redis.Cmder) {
	err := cmd.Err()
	if err == nil {
		log.Infof(ctx, "[REDIS] %+v", cmd.Args())
	} else {
		log.Errorf(ctx, "[REDIS] %+v err=%v", cmd.Args(), err)
	}
}

func NewClusterClient(addrs []string, password string) ClusterClient {
	opt := &redis.ClusterOptions{
		Addrs:              addrs,
		Password:           password,
		MaxRedirects:       3,
		ReadOnly:           true,
		RouteRandomly:      true,
		MaxRetries:         3,
		DialTimeout:        1 * time.Second,
		ReadTimeout:        3 * time.Second,
		WriteTimeout:       3 * time.Second,
		PoolSize:           100,
		MinIdleConns:       3,
		MaxConnAge:         3 * time.Hour,
		IdleTimeout:        10 * time.Minute,
		IdleCheckFrequency: 1 * time.Minute,
		OnConnect: func(ctx context.Context, cn *redis.Conn) error {
			log.Infof(ctx, "[REDIS] connect %+v", cn)
			return nil
		},
	}
	redisCluster := redis.NewClusterClient(opt)
	redisCluster.AddHook(&hook{})
	return redisCluster
}
