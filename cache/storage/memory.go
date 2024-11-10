package storage

import (
	"context"
	"strconv"
	"strings"
	"time"

	lru "github.com/hashicorp/golang-lru"
	"github.com/niming-dev/ddd-demo/go-common/strsconv"
	"github.com/niming-dev/ddd-demo/go-common/uftelemetry"
	"go.opentelemetry.io/otel/trace"
)

// DefaultMemoryStorageSize 默认的存储大小
const DefaultMemoryStorageSize = 100

type lruItem struct {
	createTime time.Time
	expiration time.Duration

	v interface{}
}

func newLruItem(v interface{}, expiration time.Duration) *lruItem {
	return &lruItem{
		createTime: time.Now(),
		expiration: expiration,
		v:          v,
	}
}

// isExpired 返回是否已过期，当expiration==0时永不过期
func (i *lruItem) isExpired() bool {
	if i.expiration == 0 {
		return false
	}

	return i.createTime.Add(i.expiration).Before(time.Now())
}

// MemoryStorage 以内存为存储的缓存结构
// 缓存策略使用ARC+过期时间
type MemoryStorage struct {
	cache *lru.ARCCache
}

// NewMemoryStorage 新建内存存储器
func NewMemoryStorage(size int) *MemoryStorage {
	if size <= 0 {
		size = DefaultMemoryStorageSize
	}

	cache, _ := lru.NewARC(size)
	return &MemoryStorage{
		cache: cache,
	}
}

// Get 查询指定name的数据
func (m MemoryStorage) Get(ctx context.Context, name string) (value interface{}, ok bool, err error) {
	return m.get(ctx, name)
}

func (m MemoryStorage) get(ctx context.Context, name string) (v interface{}, ok bool, err error) {
	span := m.startSpan(ctx, "Get", name)
	defer func() {
		if ok {
			m.logResp(span, v)
		} else {
			m.logResp(span, "<memory.Nil>")
		}
		span.End()
	}()

	v, ok = m.cache.Get(name)
	if !ok {
		return nil, false, nil
	}

	item := v.(*lruItem)
	if item.isExpired() {
		m.cache.Remove(name)
		return nil, false, nil
	}

	return item.v, true, nil
}

// Load 把指定Name的数据存储到dest变量中，如果两者类型不同则返回 TypeMismatch
// 此方法与Get作用类似，但性能比Get差，尽量使用Get方法查询数据
func (m *MemoryStorage) Load(ctx context.Context, name string, dest interface{}) (err error) {
	span := m.startSpan(ctx, "Load", name)
	defer func() {
		span.CheckError(err)
		span.End()
	}()

	v, ok, _ := m.get(span.Context, name)
	if !ok {
		return NotFound
	}

	err = scan(v, dest)
	return
}

// Store 保存一个缓存
func (m *MemoryStorage) Store(ctx context.Context, name string, value interface{}, expire time.Duration) error {
	return m.store(ctx, name, value, expire)
}

func (m MemoryStorage) store(ctx context.Context, name string, value interface{}, expire time.Duration) error {
	span := m.startSpan(ctx, "Store", name, strsconv.Any2String(value), duration2Str(expire))
	m.cache.Add(name, newLruItem(value, expire))
	span.End()
	return nil
}

// LoadStore 如果缓存存在则读取缓存，否则使用read函数获取数据并存入缓存中
func (m *MemoryStorage) LoadStore(ctx context.Context, name string, dest interface{}, read ReadFunc, expire time.Duration) (err error) {
	span := m.startSpan(ctx, "LoadStore", name, duration2Str(expire))
	defer func() {
		span.CheckError(err)
		span.End()
	}()

	v, ok, _ := m.get(span.Context, name)
	if ok {
		err = scan(v, dest)
		return
	}

	v, err = read(span.Context)
	if err != nil {
		return
	}

	_ = m.store(span.Context, name, v, expire)
	err = scan(v, dest)
	return
}

// Has 判断缓存是否存在
func (m *MemoryStorage) Has(ctx context.Context, name string) (ok bool, err error) {
	span := m.startSpan(ctx, "Has", name)
	defer func() {
		m.logResp(span, ok)
		span.End()
	}()

	v, ok := m.cache.Peek(name)
	if !ok {
		return false, nil
	}

	item := v.(*lruItem)
	if item.isExpired() {
		m.cache.Remove(name)
		return false, nil
	}

	return true, nil
}

// Del 删除一个缓存
func (m *MemoryStorage) Del(ctx context.Context, name string) error {
	span := m.startSpan(ctx, "Del", name)
	defer func() {
		m.logResp(span, "true")
		span.End()
	}()

	m.cache.Remove(name)
	return nil
}

// Expire 从当前时间重新设置有效期
func (m *MemoryStorage) Expire(ctx context.Context, name string, duration time.Duration) (ok bool, err error) {
	span := m.startSpan(ctx, "Expire", name, duration2Str(duration))
	defer func() {
		m.logResp(span, ok)
		span.End()
	}()

	v, ok := m.cache.Peek(name)
	if !ok {
		return false, nil
	}

	item := v.(*lruItem)
	item.createTime = time.Now()
	item.expiration = duration

	return true, nil
}

// ExpireAt 从当前时间重新设置有效期
func (m *MemoryStorage) ExpireAt(ctx context.Context, name string, expiration time.Time) (ok bool, err error) {
	span := m.startSpan(ctx, "ExpireAt", name, strconv.FormatInt(expiration.Unix(), 10))
	defer func() {
		m.logResp(span, ok)
		span.End()
	}()

	v, ok := m.cache.Peek(name)
	if !ok {
		return false, nil
	}

	now := time.Now()
	item := v.(*lruItem)
	item.createTime = now
	item.expiration = expiration.Sub(now)

	return true, nil
}

func (m MemoryStorage) startSpan(ctx context.Context, cmd string, args ...string) *uftelemetry.Span {
	span := uftelemetry.StartChildSpan(ctx, "cache.memory", cmd)
	span.SetAttributes(uftelemetry.Any2Attr("component", "memory"))
	span.SetAttributes(uftelemetry.Any2Attr("span.kind", "client"))
	span.SetAttributes(uftelemetry.Any2Attr("db.type", "memory"))
	span.SetAttributes(uftelemetry.Any2Attr("db.statement", cmd+" "+strings.Join(args, " ")))
	return span
}

func (m MemoryStorage) logResp(span *uftelemetry.Span, resp interface{}) {
	span.AddEvent("memory.response.content", trace.WithAttributes(uftelemetry.Any2Attr("message", resp)))
}
