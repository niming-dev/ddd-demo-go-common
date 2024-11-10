package client

import (
	"os"
	"testing"
	"time"

	"context"

	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	"github.com/niming-dev/ddd-demo/go-common/config-client/mocks"
	mvccpb "go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
)

var (
	mockController *gomock.Controller
)

/*
	Header *ResponseHeader `protobuf:"bytes,1,opt,name=header,proto3" json:"header,omitempty"`
	// kvs is the list of key-value pairs matched by the range request.
	// kvs is empty when count is requested.
	Kvs []*mvccpb.KeyValue `protobuf:"bytes,2,rep,name=kvs,proto3" json:"kvs,omitempty"`
	// more indicates if there are more keys to return in the requested range.
	More bool `protobuf:"varint,3,opt,name=more,proto3" json:"more,omitempty"`
	// count is set to the number of keys within the range when requested.
	Count                int64    `protobuf:"varint,4,opt,name=count,proto3" json:"count,omitempty"`
*/

type GetTestCase struct {
	Prefix       string
	InstanceName string
	InstanceId   string
	Key          string
	Expect       string
}

type WatchTestTrigger struct {
	Sleep          time.Duration
	Action         int
	CreateRevision int64
	ModRevision    int64
	Version        int64
	Key            string
	Value          string
}

const (
	Action_Create = 100
	Action_Put    = clientv3.EventTypePut
	Action_Delete = clientv3.EventTypeDelete
)

func (trigger *WatchTestTrigger) MakeResponse(key string) clientv3.WatchResponse {
	Created := false
	if trigger.Action == Action_Create {
		Created = true
	}
	Type := clientv3.EventTypePut
	if trigger.Action == int(Action_Delete) {
		Type = clientv3.EventTypeDelete
	}

	watchResp := clientv3.WatchResponse{
		Created: Created,
		Events: []*clientv3.Event{
			{
				Type: Type,
				Kv: &mvccpb.KeyValue{
					ModRevision: trigger.ModRevision,
				},
			},
		},
	}
	if Type == clientv3.EventTypePut {
		watchResp.Events[0].Kv.CreateRevision = trigger.CreateRevision
		watchResp.Events[0].Kv.Version = trigger.Version
		watchResp.Events[0].Kv.Key = []byte(key)
		watchResp.Events[0].Kv.Value = []byte(trigger.Value)
	}
	return watchResp
}

type WatchTestCase struct {
	Key string
}

var (
	mockEtcd = map[string]string{
		CONFIG_PREFIX + "/" + testKey:                       "value_at_/config",
		CONFIG_PREFIX + "/test-service/" + testKey:          "value_at_/config/test-service",
		CONFIG_PREFIX + "/test-service/inst0001/" + testKey: "value_at_/config/test-service/inst0001",
	}
	testKey = "LISTEN"

	getCases = []GetTestCase{
		// got root
		{Key: testKey, Expect: "value_at_/config"},
		// got service
		{Key: testKey, InstanceName: "test-service", Expect: "value_at_/config/test-service"},
		// got instance
		{Key: testKey, InstanceName: "test-service", InstanceId: "inst0001", Expect: "value_at_/config/test-service/inst0001"},
		// got service
		{Key: testKey, InstanceName: "test-service", InstanceId: "inst00012", Expect: "value_at_/config/test-service"},
		// got service
		{Prefix: "TEST", Key: testKey, InstanceName: "test-service", InstanceId: "inst00012", Expect: "value_at_/config/test-service"},
	}

	Trigers = []WatchTestTrigger{
		{
			Sleep:          time.Second * 2,
			Action:         Action_Create,
			CreateRevision: 1,
			ModRevision:    1,
			Version:        1,
			Value:          "1",
		},
		{
			Sleep:          time.Second * 2,
			Action:         int(Action_Put),
			CreateRevision: 1,
			ModRevision:    2,
			Version:        2,
			Value:          "2",
		},
		{
			Sleep:       time.Second * 2,
			Action:      int(Action_Delete),
			ModRevision: 3,
		},
	}
)

func NewMockEtcdClient(cfg clientv3.Config) (EtcdClient, error) {
	m := mocks.NewMockEtcdClient(mockController)
	m.EXPECT().Close().Return(nil).AnyTimes()

	// ctx, key, options
	m.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx interface{}, key string, option ...interface{}) (*clientv3.GetResponse, error) {
			r, ok := mockEtcd[key]
			Kvs := []*mvccpb.KeyValue{}
			if ok {
				Kvs = append(Kvs, &mvccpb.KeyValue{Key: []byte(key), Value: []byte(r)})
			}
			return &clientv3.GetResponse{
				Kvs: Kvs,
			}, nil
		},
	).AnyTimes()

	// ctx, key, options
	// 这里包含的是watch的基本测试
	m.EXPECT().Watch(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx interface{}, key string, option ...interface{}) (clientv3.WatchChan, error) {
			watchChan := make(chan clientv3.WatchResponse)
			go func() {
				for _, trigger := range Trigers {
					time.Sleep(trigger.Sleep)
					watchChan <- trigger.MakeResponse(key)
				}
			}()
			return watchChan, nil
		},
	).AnyTimes()
	return m, nil
}

func NewMockEtcdClient_forWatch(cfg clientv3.Config) (EtcdClient, error) {
	m := mocks.NewMockEtcdClient(mockController)
	m.EXPECT().Close().Return(nil).AnyTimes()

	m.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx interface{}, key string, option ...interface{}) (*clientv3.GetResponse, error) {
			r, ok := mockEtcd[key]
			Kvs := []*mvccpb.KeyValue{}
			if ok {
				Kvs = append(Kvs, &mvccpb.KeyValue{Key: []byte(key), Value: []byte(r)})
			}
			return &clientv3.GetResponse{
				Kvs: Kvs,
			}, nil
		},
	).AnyTimes()

	// ctx, key, options
	m.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx interface{}, key string, option ...interface{}) (*clientv3.GetResponse, error) {
			r, ok := mockEtcd[key]
			Kvs := []*mvccpb.KeyValue{}
			if ok {
				Kvs = append(Kvs, &mvccpb.KeyValue{Key: []byte(key), Value: []byte(r)})
			}
			return &clientv3.GetResponse{
				Kvs: Kvs,
			}, nil
		},
	).AnyTimes()

	// 初始化触发队列
	keys := []string{"/config/LISTEN", "/config/test-service/LISTEN", "/config/test-service/999999/LISTEN"}
	chs := map[string]chan clientv3.WatchResponse{}
	for _, k := range keys {
		chs[k] = make(chan clientv3.WatchResponse)

		// 针对不同的路径设置不同的Watch方法
		m.EXPECT().Watch(gomock.Any(), k, gomock.Any()).DoAndReturn(
			func(ctx interface{}, key string, option ...interface{}) (clientv3.WatchChan, error) {
				watchChan := make(chan clientv3.WatchResponse)
				go func() {
					for {
						// 等待总体调度给此路径Response消息，就发送给Watch者
						trig := <-chs[key]
						watchChan <- trig
					}
				}()
				return watchChan, nil
			},
		).AnyTimes()
	}

	// INSTANCE_NAME="test-service" INSTANCE_ID="999999"
	// 0. oldValue is value_at_/config/test-service
	// 1. root: set /config/LISTEN -> newvalue_at_/config
	//		父节点变化，不应该触配置项更新
	// 2. service: set /config/test-service/LISTEN -> value_at_/config/test-service
	//		值相同，不应该触发配置项更新
	// *3. service: set /config/test-service/LISTEN -> newvalue_at_/config/test-service
	//		触发值配置项更新			newvalue_at_/config/test-service
	// *4. instance: set /config/test-service/999999/LISTEN -> newvalue_at_/config/test-service/999999
	//		触发值配置项更新			newvalue_at_/config/test-service/999999
	// 5. service: set /config/test-service/LISTEN -> newvalue_at_/config/test-service_new
	//		当前在子节点，不应该触发配置项更新
	// *6. service: delete /config/test-service/999999/LISTEN
	//		应该触发配置项更新到上级	newvalue_at_/config/test-service_new
	// *7. instance: create /config/test-service/999999/LISTEN -> value_at_/config/test-service/999999
	//		触发值配置项更新			value_at_/config/test-service/999999
	// 8. service: delete /config/test-service/LISTEN
	//		不应该触发配置项更新
	// 9. root: delete /config/LISTEN
	//		不应该触发配置项更新
	// *10. instance: delete /config/test-service/999999/LISTEN
	//		触发值配置项更新			""
	// *11. service: create /config/test-service/LISTEN -> value_at_/config/test-service
	//		触发配置项更新 value_at_/config/test-service

	// 构造测试序列，并逐个触发
	triggers := []WatchTestTrigger{
		{ // 1. root: set /config/LISTEN -> newvalue_at_/config
			Action:         int(Action_Put),
			CreateRevision: 1, ModRevision: 2, Version: 2,
			Key:   "/config/LISTEN",
			Value: "newvalue_at_/config",
		},
		{ // 2. service: set /config/test-service/LISTEN -> value_at_/config/test-service
			Action:         int(Action_Put),
			CreateRevision: 1, ModRevision: 2, Version: 2,
			Key:   "/config/test-service/LISTEN",
			Value: "value_at_/config/test-service",
		},
		{ // 3. service: set /config/test-service/LISTEN -> newvalue_at_/config/test-service
			Action:         int(Action_Put),
			CreateRevision: 1, ModRevision: 2, Version: 2,
			Key:   "/config/test-service/LISTEN",
			Value: "newvalue_at_/config/test-service",
		},
		{ // 4. instance: set /config/test-service/999999/LISTEN -> newvalue_at_/config/test-service/999999
			Action:         int(Action_Put),
			CreateRevision: 1, ModRevision: 2, Version: 2,
			Key:   "/config/test-service/999999/LISTEN",
			Value: "newvalue_at_/config/test-service/999999",
		},
		{ // 5. service: set /config/test-service/LISTEN -> newvalue_at_/config/test-service_new
			Action:         int(Action_Put),
			CreateRevision: 1, ModRevision: 2, Version: 2,
			Key:   "/config/test-service/LISTEN",
			Value: "newvalue_at_/config/test-service_new",
		},
		{ // 6. service: delete /config/test-service/999999/LISTEN
			Action:         int(Action_Delete),
			CreateRevision: 1, ModRevision: 2, Version: 2,
			Key: "/config/test-service/999999/LISTEN",
		},
		{ // 7. instance: create /config/test-service/999999/LISTEN -> value_at_/config/test-service/999999
			Action:         int(Action_Create),
			CreateRevision: 1, ModRevision: 1, Version: 1,
			Key:   "/config/test-service/999999/LISTEN",
			Value: "value_at_/config/test-service/999999",
		},
		{ // 8. service: delete /config/test-service/LISTEN
			Action:         int(Action_Delete),
			CreateRevision: 1, ModRevision: 2, Version: 2,
			Key: "/config/test-service/LISTEN",
		},
		{ // 9. root: delete /config/LISTEN
			Action:         int(Action_Delete),
			CreateRevision: 1, ModRevision: 2, Version: 2,
			Key: "/config/LISTEN",
		},
		{ // 10. instance: delete /config/test-service/999999/LISTEN
			Action:         int(Action_Delete),
			CreateRevision: 1, ModRevision: 2, Version: 2,
			Key: "/config/test-service/999999/LISTEN",
		},
		{ // 11. service: create /config/test-service/LISTEN -> value_at_/config/test-service
			Action:         int(Action_Create),
			CreateRevision: 1, ModRevision: 1, Version: 1,
			Key:   "/config/test-service/LISTEN",
			Value: "value_at_/config/test-service",
		},
	}
	go func() {
		// 每隔一秒触发一个
		for _, tr := range triggers {
			time.Sleep(time.Second)
			chs[tr.Key] <- tr.MakeResponse(tr.Key)
		}
	}()
	return m, nil
}

func Test_Close(t *testing.T) {
	mockController = gomock.NewController(t)
	NewEtcdClient = NewMockEtcdClient
	config, err := New()
	assert.Equal(t, err, nil)
	defer config.Close()
}

func Test_Get(t *testing.T) {
	mockController = gomock.NewController(t)
	NewEtcdClient = NewMockEtcdClient

	for _, tc := range getCases {
		// don't test prefixed case, these cases should be test in test by env
		if len(tc.Prefix) > 0 {
			continue
		}
		config, err := NewWithVariable(tc.InstanceName, tc.InstanceId, "")
		assert.Equal(t, err, nil)
		val, err := config.Get(tc.Key)
		assert.Equal(t, err, nil)
		assert.Equal(t, val, tc.Expect)
		config.Close()
	}
}

func Test_Get_Env(t *testing.T) {
	mockController = gomock.NewController(t)
	NewEtcdClient = NewMockEtcdClient

	for i, tc := range getCases {
		t.Log("test case:", i)
		os.Unsetenv(INSTANCE_NAME)
		os.Unsetenv(INSTANCE_ID)
		os.Unsetenv(ETCD_CLUSTER)

		os.Setenv(tc.Prefix+INSTANCE_NAME, tc.InstanceName)
		os.Setenv(tc.Prefix+INSTANCE_ID, tc.InstanceId)
		os.Setenv(tc.Prefix+ETCD_CLUSTER, "")

		config, err := NewWithPrefix(tc.Prefix)
		assert.Equal(t, err, nil)
		val, err := config.Get(tc.Key)
		assert.Equal(t, err, nil)
		assert.Equal(t, val, tc.Expect)
		config.Close()
	}
}

func Test_Get_Prefix(t *testing.T) {
	mockController = gomock.NewController(t)
	NewEtcdClient = NewMockEtcdClient

	for _, tc := range getCases {
		os.Unsetenv(INSTANCE_NAME)
		os.Unsetenv(INSTANCE_ID)
		os.Unsetenv(ETCD_CLUSTER)

		os.Setenv(INSTANCE_NAME, tc.InstanceName)
		os.Setenv(INSTANCE_ID, tc.InstanceId)
		os.Setenv(ETCD_CLUSTER, "")

		config, err := NewWithPrefix(tc.Prefix)
		assert.Equal(t, err, nil)
		val, err := config.Get(tc.Key)
		assert.Equal(t, err, nil)
		assert.Equal(t, val, tc.Expect)
		config.Close()
	}
}

func TrueWatch(t *testing.T) {
	cli, err := originNewEtcdClient(clientv3.Config{
		Endpoints:   []string{"http://127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
		DialOptions: []grpc.DialOption{grpc.WithInsecure()},
	})
	assert.Equal(t, err, nil)
	defer cli.Close()

	ch := cli.Watch(context.TODO(), "/test")
	assert.NotEqual(t, ch, nil)

	for watchRes := range ch {
		if watchRes.Canceled {
			t.Log("canceld")
		}

		if watchRes.Created {
			t.Log("created")
		}

		if nil != watchRes.Err() {
			t.Log("err in watch", watchRes.Err())
		}
		t.Log(watchRes)
		for _, ev := range watchRes.Events {
			if ev.Type == mvccpb.DELETE {
				t.Log("key deleted", ev)
				continue
			}
			t.Log(ev, ev.PrevKv, ev.Kv)
			assert.Equal(t, err, nil)
			if ev.IsCreate() {
				t.Logf("create [%s: %s]\n", ev.Kv.Key, ev.Kv.Value)
			} else if ev.IsModify() {
				t.Logf("modify current kv: %v [%s, %s]\n",
					ev.Kv, ev.Kv.Key, ev.Kv.Value)
			}
		}
	}
}

func MockWatch(t *testing.T) {
	mockController = gomock.NewController(t)
	NewEtcdClient = NewMockEtcdClient

	cli, err := NewEtcdClient(clientv3.Config{
		Endpoints:   []string{"http://127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
		DialOptions: []grpc.DialOption{grpc.WithInsecure()},
	})
	assert.Equal(t, err, nil)
	defer cli.Close()

	ch := cli.Watch(context.TODO(), "/config/test-service/LISTEN")
	assert.NotEqual(t, ch, nil)
	t.Log("mockWatch->ch:", ch)

	count := 0
	for watchRes := range ch {
		count++
		if watchRes.Canceled {
			t.Log("canceld")
		}

		if watchRes.Created {
			t.Log("created")
		}

		if nil != watchRes.Err() {
			t.Log("err in watch", watchRes.Err())
		}
		t.Log(watchRes)
		for _, ev := range watchRes.Events {
			if ev.Type == mvccpb.DELETE {
				t.Log("key deleted", ev)
				continue
			}
			t.Log(ev, ev.PrevKv, ev.Kv)
			assert.Equal(t, err, nil)
			if ev.IsCreate() {
				t.Logf("create [%s: %s]\n", ev.Kv.Key, ev.Kv.Value)
			} else if ev.IsModify() {
				t.Logf("modify current kv: %v [%s, %s]\n",
					ev.Kv, ev.Kv.Key, ev.Kv.Value)
			} else {
				panic("unknown")
			}
		}
		if count == 3 {
			// test finished
			break
		}
	}
}

func Test_Watch10(t *testing.T) {
	ch := make(chan int)
	for i := 0; i < 30; i++ {
		go func() {
			Test_Watch(t)
			ch <- 1
		}()
	}
	for i := 0; i < 30; i++ {
		<-ch
	}
}

func Test_Watch(t *testing.T) {
	mockController = gomock.NewController(t)
	NewEtcdClient = NewMockEtcdClient_forWatch

	config, err := NewWithVariable("test-service", "999999", "")
	assert.Equal(t, err, nil)
	ch, err := config.Watch("LISTEN")
	assert.Equal(t, err, nil)
	assert.NotEqual(t, ch, nil)

	expects := []string{
		"newvalue_at_/config/test-service",        // 3.
		"newvalue_at_/config/test-service/999999", // 4.
		"newvalue_at_/config/test-service_new",    // 6.
		"value_at_/config/test-service/999999",    // 7.
		"",                                        // 10.
		"value_at_/config/test-service",           // 11.
	}
	offset := 0
	for str := range ch {
		t.Logf("[offset: %v] atch got %s\n", offset, str)
		assert.Equal(t, expects[offset], str)
		offset++
		if offset >= len(expects) {
			break
		}
	}
	config.Close()
}

func ProtoTest(t *testing.T) {
	cli, err := NewEtcdClient(clientv3.Config{
		Endpoints:   []string{"http://127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	})
	assert.Equal(t, err, nil)
	getResp, err := cli.Get(context.TODO(), "/slkdfjdl")
	assert.Equal(t, err, nil)
	t.Log(getResp)
}

func Test_Temp(t *testing.T) {
	// TrueWatch(t)
}
