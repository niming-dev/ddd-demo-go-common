package redis

import (
	"github.com/go-redis/redis/v8"
)

type Option interface {
	apply(opts *redis.ClusterOptions)
}

type funcOption struct {
	f func(opts *redis.ClusterOptions)
}

func (a *funcOption) apply(opts *redis.ClusterOptions) {
	a.f(opts)
}

func newFuncOption(f func(opts *redis.ClusterOptions)) Option {
	return &funcOption{
		f: f,
	}
}
