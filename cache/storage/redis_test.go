package storage

import (
	"context"
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis/v8"
	json "github.com/json-iterator/go"
	"github.com/spf13/cast"
	"github.com/stretchr/testify/assert"
)

type person struct {
	Name   string
	Age    int
	Height int
}

var redisTestAddr = net.JoinHostPort("127.0.0.1", "6379")

func initMiniRedis() *miniredis.Miniredis {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	redisTestAddr = s.Addr()
	return s
}

func TestRedisStorage_Store_int(t *testing.T) {
	s := initMiniRedis()
	defer s.Close()

	rs := NewRedisStorage(redis.NewClient(&redis.Options{
		Addr: redisTestAddr,
	}))

	testTables := []struct {
		Name   string
		Key    string
		Value  int
		Expire time.Duration
		Assert func(err error)
	}{
		{
			Name:   "normal 1",
			Key:    "store_int",
			Value:  1,
			Expire: time.Second * 5,
			Assert: func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			Name:   "normal 2",
			Key:    "store_int",
			Value:  3,
			Expire: time.Second * 1,
			Assert: func(err error) {
				assert.NoError(t, err)
			},
		},
	}

	for _, v := range testTables {
		t.Run(t.Name(), func(t *testing.T) {
			err := rs.Store(context.Background(), v.Key, v.Value, v.Expire)
			v.Assert(err)
			err = rs.Load(context.Background(), v.Key, &v.Value)
			v.Assert(err)
		})
	}
}

func TestRedisStorage_Store_bytes(t *testing.T) {
	s := initMiniRedis()
	defer s.Close()
	rs := NewRedisStorage(redis.NewClient(&redis.Options{
		Addr: redisTestAddr,
	}))

	test := []byte(`bytes`)
	err := rs.Store(context.Background(), "store_bytes", test, time.Second*5)
	assert.NoError(t, err)

	var test1 []byte
	err = rs.Load(context.Background(), "store_bytes", &test1)
	assert.NoError(t, err)
	assert.Equal(t, test, test1)
}

func TestRedisStorage_Store_struct(t *testing.T) {
	s := initMiniRedis()
	defer s.Close()

	test := &person{
		Name:   "boy",
		Age:    20,
		Height: 200,
	}

	rs := NewRedisStorage(redis.NewClient(&redis.Options{
		Addr: redisTestAddr,
	}))
	err := rs.Store(context.Background(), "store_struct", test, time.Second*5)

	assert.NoError(t, err)
}

func TestRedisStorage_Store_slice(t *testing.T) {
	s := initMiniRedis()
	defer s.Close()

	data := []person{
		{
			Name:   "boy",
			Age:    20,
			Height: 200,
		},
		{
			Name:   "girl",
			Age:    18,
			Height: 170,
		},
	}

	rs := NewRedisStorage(redis.NewClient(&redis.Options{
		Addr: redisTestAddr,
	}))

	err := rs.Store(context.Background(), "store_slice", data, time.Second*5)
	assert.NoError(t, err)

	var data1 []person
	err = rs.Load(context.Background(), "store_slice", &data1)
	assert.NoError(t, err)
	assert.Equal(t, data, data1)
}

func TestRedisStorage_Get_int(t *testing.T) {
	s := initMiniRedis()
	defer s.Close()

	rs := NewRedisStorage(redis.NewClient(&redis.Options{
		Addr: redisTestAddr,
	}))
	_ = rs.Store(context.Background(), "get_int", 2, time.Second*5)

	v, ok, err := rs.Get(context.Background(), "get_int")
	assert.NoError(t, err)
	assert.True(t, ok)

	vi, err := strconv.Atoi(v.(string))
	assert.NoError(t, err)
	assert.Equal(t, 2, vi)
}

func TestRedisStorage_Get_struct(t *testing.T) {
	s := initMiniRedis()
	defer s.Close()

	test := person{
		Name:   "boy",
		Age:    20,
		Height: 200,
	}

	rs := NewRedisStorage(redis.NewClient(&redis.Options{
		Addr: redisTestAddr,
	}))
	_ = rs.Store(context.Background(), "load_struct", test, time.Second*5)

	v, ok, err := rs.Get(context.Background(), "load_struct")
	assert.NoError(t, err)
	assert.True(t, ok)

	var vs person
	err = json.Unmarshal([]byte(v.(string)), &vs)
	assert.NoError(t, err)
	assert.Equal(t, test, vs)
}

func TestRedisStorage_Load_int(t *testing.T) {
	s := initMiniRedis()
	defer s.Close()

	rs := NewRedisStorage(redis.NewClient(&redis.Options{
		Addr: redisTestAddr,
	}))
	_ = rs.Store(context.Background(), "load_int", 1, time.Second*5)

	var v int
	err := rs.Load(context.Background(), "load_int", &v)

	assert.Nil(t, err)
	assert.Equal(t, 1, cast.ToInt(v))
}

func TestRedisStorage_Load_struct(t *testing.T) {
	s := initMiniRedis()
	defer s.Close()

	test := person{
		Name:   "boy",
		Age:    20,
		Height: 200,
	}

	rs := NewRedisStorage(redis.NewClient(&redis.Options{
		Addr: redisTestAddr,
	}))
	_ = rs.Store(context.Background(), "load_struct", test, time.Second*5)

	var v person
	err := rs.Load(context.Background(), "load_struct", &v)

	assert.Nil(t, err)
	assert.Equal(t, test, v)
}

func TestRedisStorage_Load_slice(t *testing.T) {
	s := initMiniRedis()
	defer s.Close()

	data := []person{
		{
			Name:   "boy",
			Age:    20,
			Height: 200,
		},
		{
			Name:   "girl",
			Age:    18,
			Height: 170,
		},
	}

	rs := NewRedisStorage(redis.NewClient(&redis.Options{
		Addr: redisTestAddr,
	}))
	_ = rs.Store(context.Background(), "load_struct", data, time.Second*5)

	var v []person
	err := rs.Load(context.Background(), "load_struct", &v)

	assert.Nil(t, err)
	assert.Equal(t, data, v)
}

func TestRedisStorage_LoadStore_TypeNotMatch(t *testing.T) {
	s := initMiniRedis()
	defer s.Close()

	rs := NewRedisStorage(redis.NewClient(&redis.Options{
		Addr: redisTestAddr,
	}))

	var v int
	err := rs.LoadStore(context.Background(), "load_store", &v, func(ctx context.Context) (interface{}, error) {
		return "1", nil
	}, time.Second*5)

	assert.EqualError(t, err, TypeMismatch.Error())
}

func TestRedisStorage_LoadStore(t *testing.T) {
	s := initMiniRedis()
	defer s.Close()

	rs := NewRedisStorage(redis.NewClient(&redis.Options{
		Addr: redisTestAddr,
	}))

	var v int
	err := rs.LoadStore(context.Background(), "load_store", &v, func(ctx context.Context) (interface{}, error) {
		return 1, nil
	}, time.Second*5)

	assert.NoError(t, err)
	assert.Equal(t, 1, v)
}

func TestRedisStorage_LoadStore_Struct_TypeMismatch(t *testing.T) {
	s := initMiniRedis()
	defer s.Close()

	rs := NewRedisStorage(redis.NewClient(&redis.Options{
		Addr: redisTestAddr,
	}))

	result := person{
		Name:   "boy",
		Age:    20,
		Height: 200,
	}
	var v *person
	err := rs.LoadStore(context.Background(), "load_store", &v, func(ctx context.Context) (interface{}, error) {
		return result, nil
	}, time.Second*5)

	assert.Equal(t, TypeMismatch, err)
}

func TestRedisStorage_LoadStore_Struct_ptr(t *testing.T) {
	s := initMiniRedis()
	defer s.Close()

	rs := NewRedisStorage(redis.NewClient(&redis.Options{
		Addr: redisTestAddr,
	}))

	result := &person{
		Name:   "boy",
		Age:    20,
		Height: 200,
	}
	var v *person
	err := rs.LoadStore(context.Background(), "load_store", &v, func(ctx context.Context) (interface{}, error) {
		return result, nil
	}, time.Second*5)

	assert.NoError(t, err)
	assert.Equal(t, result, v)
}

func TestRedisStorage_LoadStore_Struct(t *testing.T) {
	s := initMiniRedis()
	defer s.Close()

	rs := NewRedisStorage(redis.NewClient(&redis.Options{
		Addr: redisTestAddr,
	}))

	result := &person{
		Name:   "boy",
		Age:    20,
		Height: 200,
	}
	var v person
	err := rs.LoadStore(context.Background(), "load_store", &v, func(ctx context.Context) (interface{}, error) {
		return result, nil
	}, time.Second*5)

	assert.NoError(t, err)
	assert.Equal(t, result, &v)
}

func TestRedisStorage_LoadStore_Struct_nil(t *testing.T) {
	s := initMiniRedis()
	defer s.Close()

	rs := NewRedisStorage(redis.NewClient(&redis.Options{
		Addr: redisTestAddr,
	}))

	var v *person
	err := rs.LoadStore(context.Background(), "load_store", &v, func(ctx context.Context) (interface{}, error) {
		return nil, nil
	}, time.Second*5)

	assert.NoError(t, err)
	assert.Equal(t, (*person)(nil), v)
}

func TestRedisStorage_LoadStore_slice(t *testing.T) {
	s := initMiniRedis()
	defer s.Close()

	data := []person{
		{
			Name:   "boy",
			Age:    20,
			Height: 200,
		},
		{
			Name:   "girl",
			Age:    18,
			Height: 170,
		},
	}

	rs := NewRedisStorage(redis.NewClient(&redis.Options{
		Addr: redisTestAddr,
	}))

	var v []person
	err := rs.LoadStore(context.Background(), "load_store_slice", &v, func(ctx context.Context) (interface{}, error) {
		return data, nil
	}, time.Second*5)

	assert.NoError(t, err)
	assert.Equal(t, data, v)
}

func TestRedisStorage_LoadStore_bool(t *testing.T) {
	s := initMiniRedis()
	defer s.Close()

	rs := NewRedisStorage(redis.NewClient(&redis.Options{
		Addr: redisTestAddr,
	}))

	var v bool
	err := rs.LoadStore(context.Background(), "load_store_bool", &v, func(ctx context.Context) (interface{}, error) {
		return true, nil
	}, time.Second*5)

	assert.NoError(t, err)
	assert.Equal(t, true, v)

	err = rs.Load(context.Background(), "load_store_bool", &v)
	assert.NoError(t, err)
}

func TestRedisStorage_Del(t *testing.T) {
	s := initMiniRedis()
	defer s.Close()

	rs := NewRedisStorage(redis.NewClient(&redis.Options{
		Addr: redisTestAddr,
	}))
	err := rs.Store(context.Background(), "del", 1, time.Second*5)
	assert.NoError(t, err)

	err = rs.Del(context.Background(), "del")
	assert.NoError(t, err)

	var v int
	err = rs.Load(context.Background(), "del", &v)
	assert.EqualError(t, err, NotFound.Error())
}

func TestRedisStorage_Has(t *testing.T) {
	s := initMiniRedis()
	defer s.Close()

	rs := NewRedisStorage(redis.NewClient(&redis.Options{
		Addr: redisTestAddr,
	}))
	err := rs.Store(context.Background(), "has", 1, time.Second*5)
	assert.NoError(t, err)

	exists, err := rs.Has(context.Background(), "has")
	assert.NoError(t, err)
	assert.True(t, exists)

	err = rs.Del(context.Background(), "has")
	assert.NoError(t, err)

	exists, err = rs.Has(context.Background(), "has")
	assert.NoError(t, err)
	assert.False(t, exists)
}

func TestRedisStorage_Expire(t *testing.T) {
	s := initMiniRedis()
	defer s.Close()
	rs := NewRedisStorage(redis.NewClient(&redis.Options{
		Addr: redisTestAddr,
	}))

	err := rs.Store(context.Background(), "key", 1, time.Second)
	assert.NoError(t, err)

	_, _ = rs.Expire(context.Background(), "key", time.Second*2)
	s.FastForward(time.Second)
	exists, err := rs.Has(context.Background(), "key")
	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestRedisStorage_ExpireAt(t *testing.T) {
	s := initMiniRedis()
	defer s.Close()
	rs := NewRedisStorage(redis.NewClient(&redis.Options{
		Addr: redisTestAddr,
	}))

	err := rs.Store(context.Background(), "key", 1, time.Second)
	assert.NoError(t, err)

	_, _ = rs.ExpireAt(context.Background(), "key", time.Now().Add(time.Second*2))
	s.FastForward(time.Second)
	exists, err := rs.Has(context.Background(), "key")
	assert.NoError(t, err)
	assert.True(t, exists)
}

func BenchmarkRedisStorage_LoadStore(b *testing.B) {
	s := initMiniRedis()
	defer s.Close()

	data := []person{
		{
			Name:   "boy",
			Age:    20,
			Height: 200,
		},
		{
			Name:   "girl",
			Age:    18,
			Height: 170,
		},
	}

	rs := NewRedisStorage(redis.NewClient(&redis.Options{
		Addr: redisTestAddr,
	}))

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		var v []person
		_ = rs.LoadStore(context.Background(), "load_store", &v, func(ctx context.Context) (interface{}, error) {
			return data, nil
		}, time.Second*5)
	}
}
