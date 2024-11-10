package storage

import (
	"context"
	"encoding"
	"errors"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/json-iterator/go"
	"github.com/niming-dev/ddd-demo/go-common/strsconv"
	"github.com/niming-dev/ddd-demo/go-common/uftelemetry"
	"go.opentelemetry.io/otel/trace"
)

// RedisStorage redis存储器
type RedisStorage struct {
	cmd redis.Cmdable
}

// NewRedisStorage 新建redis存储器
func NewRedisStorage(rc redis.Cmdable) *RedisStorage {
	return &RedisStorage{
		cmd: rc,
	}
}

// Get 查询数据
func (r RedisStorage) Get(ctx context.Context, name string) (result interface{}, ok bool, err error) {
	span := r.startSpan(ctx, "Get", name)
	defer func() {
		r.logResp(span, result, err)
		r.checkErr(span, err)
		span.End()
	}()

	result, err = r.cmd.Get(ctx, name).Result()
	if err != nil {
		return result, false, err
	}

	return result, true, nil
}

func (r *RedisStorage) load(ctx context.Context, name string, dest interface{}) (err error) {
	span := r.startSpan(ctx, "Get", name)
	defer func() {
		r.checkErr(span, err)
		span.End()
	}()

	cmd := r.cmd.Get(ctx, name)
	err = cmd.Err()
	r.logResp(span, cmd.Val(), err)
	if err != nil {
		return
	}

	if cmd.Val() == "" {
		v := reflect.ValueOf(dest)
		v.Elem().Set(reflect.Zero(v.Elem().Type()))
		return nil
	}

	switch dest.(type) {
	case nil, *string, *[]byte, *int, *int8, *int16, *int32, *int64, *uint, *uint8, *uint16, *uint32, *uint64, *float32,
		*float64, *bool, encoding.BinaryUnmarshaler:
		err = cmd.Scan(dest)
		return
	default:
		var bs []byte
		bs, err = cmd.Bytes()
		if err != nil {
			return
		}
		err = jsoniter.Unmarshal(bs, dest)
		return
	}
}

// Load 读取数据, dest必须是指针类型
func (r *RedisStorage) Load(ctx context.Context, name string, dest interface{}) error {
	err := r.load(ctx, name, dest)
	if err == redis.Nil {
		return NotFound
	}

	return err
}

func (r *RedisStorage) store(ctx context.Context, name string, value interface{}, expire time.Duration) (err error) {
	span := r.startSpan(ctx, "Set", name, strsconv.Any2String(value), duration2Str(expire))
	defer func() {
		r.checkErr(span, err)
		span.End()
	}()

	switch value.(type) {
	case nil, string, []byte, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64,
		encoding.BinaryMarshaler, bool:
	default:
		value, err = jsoniter.Marshal(value)
		if err != nil {
			return
		}
		// 把 value=nil 产生的 null 转为空字符串
		if value == "null" {
			value = ""
		}
	}

	result, err := r.cmd.Set(ctx, name, value, expire).Result()
	r.logResp(span, result, err)
	if err != nil {
		return
	}

	if result != "OK" {
		err = errors.New("set result is not OK: " + result)
		return
	}

	return nil
}

// Store 存储数据
func (r *RedisStorage) Store(ctx context.Context, name string, value interface{}, expire time.Duration) error {
	return r.store(ctx, name, value, expire)
}

// LoadStore 如果存在则读取数据，如果不存在则存储数据, dest必须是指针类型
func (r *RedisStorage) LoadStore(ctx context.Context, name string, dest interface{}, read ReadFunc, expire time.Duration) error {
	span := r.startSpan(ctx, "LoadStore", name, duration2Str(expire))
	defer func() {
		span.End()
	}()

	err := r.load(span.Context, name, dest)
	if err == redis.Nil {
		value, err := read(span.Context)
		if err != nil {
			span.Error(err)
			return err
		}

		err = scan(value, dest)
		if err != nil {
			span.Error(err)
			return err
		}

		return r.store(span.Context, name, value, expire)
	}

	return err
}

// Has 判断数据是否存在
func (r *RedisStorage) Has(ctx context.Context, name string) (exists bool, err error) {
	span := r.startSpan(ctx, "Exists", name)
	defer func() {
		r.checkErr(span, err)
		span.End()
	}()

	result, err := r.cmd.Exists(ctx, name).Result()
	r.logResp(span, result, err)
	if err != nil {
		return false, err
	}

	return result == 1, nil
}

// Del 删除数据
func (r *RedisStorage) Del(ctx context.Context, name string) (err error) {
	span := r.startSpan(ctx, "Del", name)
	defer func() {
		r.checkErr(span, err)
		span.End()
	}()

	result, err := r.cmd.Del(ctx, name).Result()
	r.logResp(span, result, err)
	if err != nil {
		return err
	}

	if result > 0 {
		return nil
	}

	return NotFound
}

// Expire 从当前时间重新设置有效期
func (r *RedisStorage) Expire(ctx context.Context, name string, duration time.Duration) (ok bool, err error) {
	span := r.startSpan(ctx, "Expire", name, duration2Str(duration))
	defer func() {
		r.logResp(span, ok, err)
		r.checkErr(span, err)
		span.End()
	}()

	ok, err = r.cmd.Expire(ctx, name, duration).Result()
	return
}

// ExpireAt 从当前时间重新设置有效期
func (r *RedisStorage) ExpireAt(ctx context.Context, name string, expiration time.Time) (ok bool, err error) {
	span := r.startSpan(ctx, "ExpireAt", name, strconv.FormatInt(expiration.Unix(), 10))
	defer func() {
		r.logResp(span, ok, err)
		r.checkErr(span, err)
		span.End()
	}()

	ok, err = r.cmd.ExpireAt(ctx, name, expiration).Result()
	return
}

func (r RedisStorage) startSpan(ctx context.Context, cmd string, args ...string) *uftelemetry.Span {
	span := uftelemetry.StartChildSpan(ctx, "cache.redis", cmd)
	span.SetAttributes(uftelemetry.Any2Attr("component", "redis"))
	span.SetAttributes(uftelemetry.Any2Attr("span.kind", "client"))
	span.SetAttributes(uftelemetry.Any2Attr("db.type", "redis"))
	span.SetAttributes(uftelemetry.Any2Attr("db.statement", cmd+" "+strings.Join(args, " ")))

	return span
}

func (r RedisStorage) logResp(span *uftelemetry.Span, resp interface{}, err error) {
	message := resp
	if err == redis.Nil {
		message = "<redis.Nil>"
	}
	span.AddEvent("redis.response.content", trace.WithAttributes(uftelemetry.Any2Attr("message", message)))
}

func (r RedisStorage) checkErr(span *uftelemetry.Span, err error) {
	if err != nil && err != redis.Nil {
		span.CheckError(err)
	}
}
