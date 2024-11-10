package storage

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/spf13/cast"
	"github.com/stretchr/testify/assert"
)

func TestMemoryStorage_Store_int(t *testing.T) {
	ms := NewMemoryStorage(100)
	err := ms.Store(context.Background(), "store_int", 1, time.Second*5)
	assert.NoError(t, err)
}

func TestMemoryStorage_Store_outSize(t *testing.T) {
	ms := NewMemoryStorage(100)

	for i := 0; i < 10000; i++ {
		err := ms.Store(context.Background(), "store_int"+strconv.Itoa(i), 1, time.Minute)
		assert.NoError(t, err)
	}

	v, ok, err := ms.Get(context.Background(), "store_int1")
	assert.NoError(t, err)
	assert.False(t, ok)
	assert.Nil(t, v)

	v, ok, err = ms.Get(context.Background(), "store_int9999")
	assert.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, 1, v)
}

func TestMemoryStorage_Store_struct(t *testing.T) {
	test := &person{
		Name:   "boy",
		Age:    20,
		Height: 200,
	}

	ms := NewMemoryStorage(100)
	err := ms.Store(context.Background(), "store_struct", test, time.Second*5)
	assert.NoError(t, err)
}

func TestMemoryStorage_Store_slice(t *testing.T) {
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

	ms := NewMemoryStorage(100)
	err := ms.Store(context.Background(), "store_slice", data, time.Second*5)
	assert.NoError(t, err)
}

func TestMemoryStorage_Load_int(t *testing.T) {
	ms := NewMemoryStorage(100)
	err := ms.Store(context.Background(), "load_int", 1, time.Second*5)
	assert.NoError(t, err)

	var v int
	err = ms.Load(context.Background(), "load_int", &v)
	assert.NoError(t, err)
	assert.Equal(t, 1, cast.ToInt(v))
}

func TestMemoryStorage_Load_struct(t *testing.T) {
	test := person{
		Name:   "boy",
		Age:    20,
		Height: 200,
	}

	ms := NewMemoryStorage(100)
	err := ms.Store(context.Background(), "load_struct", test, time.Second*5)
	assert.NoError(t, err)

	var v person
	err = ms.Load(context.Background(), "load_struct", &v)
	assert.NoError(t, err)
	assert.Equal(t, test, v)
}

func TestMemoryStorage_Load_slice(t *testing.T) {
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

	ms := NewMemoryStorage(100)
	err := ms.Store(context.Background(), "load_struct", data, time.Second*5)
	assert.NoError(t, err)

	var v []person
	err = ms.Load(context.Background(), "load_struct", &v)
	assert.NoError(t, err)
	assert.Equal(t, data, v)
}

func TestMemoryStorage_LoadStore(t *testing.T) {
	ms := NewMemoryStorage(100)

	var v int
	err := ms.LoadStore(context.Background(), "load_store", &v, func(ctx context.Context) (interface{}, error) {
		return 1, nil
	}, time.Second*5)

	assert.NoError(t, err)
	assert.Equal(t, 1, v)
}

func TestMemoryStorage_LoadStore_slice(t *testing.T) {
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

	ms := NewMemoryStorage(100)

	var v []person
	err := ms.LoadStore(context.Background(), "load_store_slice", &v, func(ctx context.Context) (interface{}, error) {
		return data, nil
	}, time.Second*5)

	assert.NoError(t, err)
	assert.Equal(t, data, v)
}

func TestMemoryStorage_Del(t *testing.T) {
	ms := NewMemoryStorage(100)
	err := ms.Store(context.Background(), "del", 1, time.Second*5)
	assert.NoError(t, err)

	err = ms.Del(context.Background(), "del")
	assert.NoError(t, err)

	var v int
	err = ms.Load(context.Background(), "del", &v)
	assert.EqualError(t, err, NotFound.Error())
}

func TestMemoryStorage_Has(t *testing.T) {
	ms := NewMemoryStorage(100)
	err := ms.Store(context.Background(), "has", 1, time.Second*5)
	assert.NoError(t, err)

	exists, err := ms.Has(context.Background(), "has")
	assert.NoError(t, err)
	assert.True(t, exists)

	err = ms.Del(context.Background(), "has")
	assert.NoError(t, err)

	exists, err = ms.Has(context.Background(), "has")
	assert.NoError(t, err)
	assert.False(t, exists)
}

func TestMemoryStorage_Expire(t *testing.T) {
	ms := NewMemoryStorage(100)

	err := ms.Store(context.Background(), "key", 1, time.Second)
	assert.NoError(t, err)
	time.Sleep(time.Second)

	exists, err := ms.Has(context.Background(), "key")
	assert.NoError(t, err)
	assert.False(t, exists)

	err = ms.Store(context.Background(), "key", 1, time.Second)
	assert.NoError(t, err)
	_, err = ms.Expire(context.Background(), "key", time.Second*2)
	assert.NoError(t, err)
	time.Sleep(time.Second)

	exists, err = ms.Has(context.Background(), "key")
	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestMemoryStorage_ExpireAt(t *testing.T) {
	ms := NewMemoryStorage(100)

	err := ms.Store(context.Background(), "key", 1, time.Second)
	assert.NoError(t, err)
	time.Sleep(time.Second)
	exists, err := ms.Has(context.Background(), "key")
	assert.NoError(t, err)
	assert.False(t, exists)

	err = ms.Store(context.Background(), "key", 1, time.Second)
	assert.NoError(t, err)
	_, err = ms.ExpireAt(context.Background(), "key", time.Now().Add(time.Second*2))
	assert.NoError(t, err)
	time.Sleep(time.Second)
	exists, err = ms.Has(context.Background(), "key")
	assert.NoError(t, err)
	assert.True(t, exists)
}

func BenchmarkMemoryStorage_Get_Int_100(b *testing.B) {
	ms := NewMemoryStorage(100)
	for i := 0; i < 100; i++ {
		_ = ms.Store(context.Background(), strconv.Itoa(i), i, 0)
	}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		for i := 0; i < 100; i++ {
			if v, ok, _ := ms.Get(context.Background(), strconv.Itoa(i)); ok {
				_ = v.(int)
			}
		}
	}
}

func BenchmarkMemoryStorage_Get_Struct_100(b *testing.B) {
	data := person{
		Name:   "boy",
		Age:    20,
		Height: 200,
	}

	ms := NewMemoryStorage(100)
	for i := 0; i < 100; i++ {
		_ = ms.Store(context.Background(), strconv.Itoa(i), data, 0)
	}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		for i := 0; i < 100; i++ {
			if v, ok, _ := ms.Get(context.Background(), strconv.Itoa(i)); ok {
				_ = v.(person)
			}
		}
	}
}

func BenchmarkMemoryStorage_Load_Struct_100(b *testing.B) {
	data := person{
		Name:   "boy",
		Age:    20,
		Height: 200,
	}

	ms := NewMemoryStorage(100)
	for i := 0; i < 100; i++ {
		_ = ms.Store(context.Background(), strconv.Itoa(i), data, 0)
	}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		for i := 0; i < 100; i++ {
			var v person
			_ = ms.Load(context.Background(), strconv.Itoa(i), &v)
		}
	}
}

func BenchmarkMemoryStorage_LoadStore_Struct(b *testing.B) {
	data := person{
		Name:   "boy",
		Age:    20,
		Height: 200,
	}

	ms := NewMemoryStorage(100)

	var v person
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_ = ms.LoadStore(context.Background(), "load_store", &v, func(ctx context.Context) (interface{}, error) {
			return data, nil
		}, time.Second*5)
	}
}

func BenchmarkMemoryStorage_LoadStore_StructSlice(b *testing.B) {
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

	ms := NewMemoryStorage(100)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		var v []person
		_ = ms.LoadStore(context.Background(), "load_store", &v, func(ctx context.Context) (interface{}, error) {
			return data, nil
		}, time.Second*5)
	}
}

func BenchmarkMemoryStorage_LoadStore_MapSS(b *testing.B) {
	data := map[string]string{
		"Name":   "boy",
		"Age":    "20",
		"Height": "200",
	}

	ms := NewMemoryStorage(100)

	v := make(map[string]string)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_ = ms.LoadStore(context.Background(), "load_store", &v, func(ctx context.Context) (interface{}, error) {
			return data, nil
		}, time.Second*5)
	}
}

func BenchmarkMemoryStorage_LoadStore_MapSI(b *testing.B) {
	data := map[string]interface{}{
		"Name":   "boy",
		"Age":    20,
		"Height": 200,
	}

	ms := NewMemoryStorage(100)

	v := make(map[string]interface{})
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_ = ms.LoadStore(context.Background(), "load_store", &v, func(ctx context.Context) (interface{}, error) {
			return data, nil
		}, time.Second*5)
	}
}
