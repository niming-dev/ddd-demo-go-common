package cache

import (
	"github.com/go-redis/redis/v8"
	"github.com/niming-dev/ddd-demo/go-common/cache/storage"
)

// Client 缓存客户端
type Client struct {
	storage.Storage
}

// NewRedis 生成一个以Redis为存储容器的缓存实例
func NewRedis(rc redis.Cmdable) *Client {
	return &Client{
		storage.NewRedisStorage(rc),
	}
}

// NewMemory 生成一个以内存为存储容器的缓存实例
func NewMemory(size int) *Client {
	return &Client{
		storage.NewMemoryStorage(size),
	}
}
