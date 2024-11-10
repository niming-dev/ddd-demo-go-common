package dlock

import (
	"reflect"
	"time"

	"github.com/google/uuid"
)

var (
	tryLockType         = reflect.TypeOf((*TryLocker)(nil)).Elem()
	lockWithTimeoutType = reflect.TypeOf((*LockWithTimeout)(nil)).Elem()

	defaultLocker Locker = NewRedisLocker(nil)
)

type Locker interface {
	Lock(key, value string, expire time.Duration)
	Unlock(key, value string)
}

type TryLocker interface {
	TryLock(key, value string, expire time.Duration) bool
}

type LockWithTimeout interface {
	LockWithTimeout(key, value string, expire, timeout time.Duration) bool
}

// 分布式锁
type DistributeLock struct {
	// 加锁的key
	key string
	// 加锁的值，用来唯一标识一把锁，解锁时需要value一致才能解锁
	value string
	// 锁的失效时间，默认是60秒
	expireTime time.Duration
	// 真正的锁对象
	locker Locker
}

// 新建一个分布式锁对象
func NewDistributeLock(key string) *DistributeLock {
	return &DistributeLock{key: key, expireTime: 60 * time.Second, locker: defaultLocker, value: uuid.New().String()}
}

// 设置锁的失效时间
func (dl *DistributeLock) SetExpireTime(d time.Duration) *DistributeLock {
	dl.expireTime = d
	return dl
}

// 设置锁接口
func (dl *DistributeLock) SetLocker(l Locker) *DistributeLock {
	dl.locker = l
	return dl
}

// 加锁不等待，返回true表示拿到了锁
func (dl *DistributeLock) TryLock() bool {
	typ := reflect.TypeOf(dl.locker)
	if typ.Implements(tryLockType) {
		return dl.locker.(TryLocker).TryLock(dl.key, dl.value, dl.expireTime)
	} else {
		return false
	}
}

// 加锁阻塞等待
func (dl *DistributeLock) Lock() {
	dl.locker.Lock(dl.key, dl.value, dl.expireTime)
}

// 加锁等待一个超时时间，返回true表示拿到了锁
func (dl *DistributeLock) LockWithTimeout(t time.Duration) bool {
	typ := reflect.TypeOf(dl.locker)
	if typ.Implements(lockWithTimeoutType) {
		return dl.locker.(LockWithTimeout).LockWithTimeout(dl.key, dl.value, dl.expireTime, t)
	} else {
		return false
	}
}

// 解锁
func (dl *DistributeLock) Unlock() {
	dl.locker.Unlock(dl.key, dl.value)
}
