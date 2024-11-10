package dlock

import (
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis"
)

const (
	REDIS_SERVER   = "REDIS_SERVER"
	REDIS_PASSWORD = "REDIS_PASSWORD"
)

var (
	globalRedisClient RedisClient
	one               sync.Once
)

type RedisClient interface {
	SetNX(key string, value interface{}, expiration time.Duration) *redis.BoolCmd
	Eval(script string, keys []string, args ...interface{}) *redis.Cmd
	Publish(channel string, message interface{}) *redis.IntCmd
}

type RedisLocker struct {
	cli RedisClient
	m   sync.Mutex
}

func _get_globalRedisClient() RedisClient {
	if nil == globalRedisClient {
		one.Do(func() {
			redisServer := os.Getenv(REDIS_SERVER)
			redisPassword := os.Getenv(REDIS_PASSWORD)
			servers := strings.Split(redisServer, ",")
			if len(servers) <= 1 {
				globalRedisClient = redis.NewClient(&redis.Options{
					Addr:     redisServer,
					Password: redisPassword, // no password set
					DB:       0,             // use default DB
				})
			} else {
				globalRedisClient = redis.NewClusterClient(&redis.ClusterOptions{
					Addrs:        servers,
					Password:     redisPassword,
					MaxRetries:   3,
					PoolSize:     50,
					MinIdleConns: 5,
					MaxConnAge:   time.Minute * 30,
				})
			}
		})
	}
	return globalRedisClient
}

// 新建一个RedisLocker，如果没有指定RedisClient则使用全局共用的
func NewRedisLocker(cli RedisClient) *RedisLocker {
	ret := &RedisLocker{cli: cli}
	if cli == nil {
		ret.cli = _get_globalRedisClient()
	}
	return ret
}

func SetGlobalRedisClient(cli RedisClient) {
	globalRedisClient = cli
}

func (l *RedisLocker) SetRedisClient(cli RedisClient) *RedisLocker {
	l.m.Lock()
	defer l.m.Unlock()

	l.cli = cli
	return l
}

func (l *RedisLocker) checkClient() {
	l.m.Lock()
	defer l.m.Unlock()
	if nil == l.cli {
		l.cli = _get_globalRedisClient()
	}
}

func (l *RedisLocker) Lock(key, value string, expire time.Duration) {
	l.checkClient()

	for !l.TryLock(key, value, expire) {
		// 每隔10毫秒去尝试获取锁
		time.Sleep(time.Millisecond * 10)
	}
}

func (l *RedisLocker) LockWithTimeout(key, value string, expire, timeout time.Duration) bool {
	l.checkClient()

	timeBegin := time.Now()
	timeCurrent := timeBegin
	timeoutTime := timeBegin.Add(timeout)
	ret := true
	for !l.TryLock(key, value, expire) {
		// 每隔10毫秒去尝试获取锁
		time.Sleep(time.Millisecond * 10)
		timeCurrent = timeCurrent.Add(time.Millisecond * 10)
		if timeCurrent.After(timeoutTime) {
			ret = false
			break
		}
	}
	return ret
}

func (l *RedisLocker) Unlock(key, value string) {
	l.checkClient()
	if nil == l.cli {
		// 由于redis配置错误或网络错误，无法获取锁，直接返回成功
		return
	}

	script :=
		`if ( redis.call('get', KEYS[1]) == ARGV[1] ) then 
			return redis.call('del', KEYS[1]) 
		end`

	cmd := l.cli.Eval(script, []string{key}, []string{value})
	intr, err := cmd.Result()
	if nil != err {
		log.Printf("redis unlock got error: %v, %v\n", intr, err)
	}
}

func (l *RedisLocker) TryLock(key, value string, expire time.Duration) bool {
	l.checkClient()
	if nil == l.cli {
		// 由于redis配置错误或网络错误，无法获取锁，直接返回成功
		return true
	}

	res := l.cli.SetNX(key, value, expire)
	lockSuccess, err := res.Result()
	if nil != err {
		// 由于redis配置错误或网络错误，无法获取锁，直接返回成功
		return true
	}
	if !lockSuccess {
		// 未取得锁的时候返回失败
		return false
	}
	return true
}
