package dlock

import (
	"log"
	"testing"
	"time"

	"github.com/go-redis/redis"
)

func Test_Locker(t *testing.T) {
	l := NewDistributeLock("abcd")
	l.Lock()

	if l.TryLock() {
		t.Fatal("TryLock should return false")
	}

	t1 := time.Now()
	if l.LockWithTimeout(time.Second) {
		t.Fatal("LockWithTimeout should return false")
	}
	if time.Now().Before(t1.Add(time.Second)) {
		t.Fatal("Lock not block for 1 second")
	}
	l.Unlock()
}

func Test_InvalidLocker(t *testing.T) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:        "9.9.9.9:2222",
		Password:    "", // no password set
		DB:          0,  // use default DB
		DialTimeout: time.Second,
	})

	l := NewDistributeLock("abcd").SetLocker(
		NewRedisLocker(redisClient),
	)

	// lock should always success
	l.Lock()

	log.Println("TryLock")
	if !l.TryLock() {
		t.Fatal("TryLock should return true")
	}

	log.Println("LockWithTimeout")
	t1 := time.Now()
	if !l.LockWithTimeout(5 * time.Second) {
		t.Fatal("LockWithTimeout should return true")
	}
	if time.Now().After(t1.Add(5 * time.Second)) {
		t.Fatal("Lock shouldn't block")
	}

	// lock even success
	log.Println("Lock")
	l.Lock()
	log.Println("Lock")
	l.Lock()
	log.Println("Lock")
	l.Lock()

	l.Unlock()
}
