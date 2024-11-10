package storage

import (
	"context"
	"time"
)

// Storage 存储器接口
type Storage interface {
	// Load 读取数据, dest必须是指针类型
	Load(ctx context.Context, name string, dest interface{}) (err error)
	// Get 查询数据
	Get(ctx context.Context, name string) (value interface{}, ok bool, err error)
	// Store 存储数据
	Store(ctx context.Context, name string, value interface{}, duration time.Duration) error
	// LoadStore 如果存在则读取数据，如果不存在则存储数据, dest必须是指针类型
	LoadStore(ctx context.Context, name string, dest interface{}, read ReadFunc, duration time.Duration) error
	// Has 判断数据是否存在
	Has(ctx context.Context, name string) (bool, error)
	// Del 删除数据
	Del(ctx context.Context, name string) error
	// Expire 从当前时间重新设置有效期
	Expire(ctx context.Context, name string, duration time.Duration) (bool, error)
	// ExpireAt 设置缓存具体过期时间点
	ExpireAt(ctx context.Context, name string, expiration time.Time) (bool, error)
}

// ReadFunc LoadStore数据不存在时调用 读取函数 获取数据
type ReadFunc func(ctx context.Context) (interface{}, error)
