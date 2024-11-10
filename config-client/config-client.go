package client

import (
	"context"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/niming-dev/ddd-demo/go-common/uferror"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
)

const (
	CONFIG_PREFIX = "/config"

	INSTANCE_NAME = "INSTANCE_NAME"
	INSTANCE_ID   = "INSTANCE_ID"
	// 使用逗号分隔的多个endpoint，例如http://127.0.0.1:2222,http://127.0.0.1:3333
	ETCD_CLUSTER = "ETCD_CLUSTER"

	ERR_CODE_TIMEOUT   = -100
	ERR_GRPC           = -101
	ERR_NOT_FOUND      = -102
	ERR_TYPE_NOT_MATCH = -200
)

var (
	ErrTimeout      = uferror.New(ERR_CODE_TIMEOUT, "timeout")
	ErrGrpc         = uferror.New(ERR_GRPC, "")
	ErrNotFound     = uferror.New(ERR_NOT_FOUND, "not found")
	ErrTypeNotMatch = uferror.New(ERR_TYPE_NOT_MATCH, "type not match")
)

type ConfigClient struct {
	InstanceName string
	InstanceId   string
	EtcdCluster  string

	m                sync.Mutex
	Client           EtcdClient
	GrpcContext      context.Context
	GrpcContextClear context.CancelFunc
}

func (client *ConfigClient) GetEnvString() (string, string, string) {
	return client.InstanceName, client.InstanceId, client.EtcdCluster
}

func New() (*ConfigClient, uferror.UFError) {
	return NewWithPrefix("")
}

// 可以为环境变量指定前缀，便于测试时启动多个不同的实例
// 设置了DEV_为PREFIX后，会读取环境变量DEV_INSTANCE_NAME, DEV_INSTANCE_ID, DEV_ETCD_CLUSTER
func NewWithPrefix(prefix string) (*ConfigClient, uferror.UFError) {
	log.Printf(
		`INSTANCE_NAME: %s\n
		INSTANCE_ID: %s\n,
		ETCD_CLUSTER: %s\n`,
		os.Getenv(prefix+INSTANCE_NAME),
		os.Getenv(prefix+INSTANCE_ID),
		os.Getenv(prefix+ETCD_CLUSTER),
	)
	return NewWithVariable(
		os.Getenv(prefix+INSTANCE_NAME),
		os.Getenv(prefix+INSTANCE_ID),
		os.Getenv(prefix+ETCD_CLUSTER),
	)
}

func NewWithVariable(instanceName, instanceId, etcdCluster string) (*ConfigClient, uferror.UFError) {
	client := &ConfigClient{
		InstanceName: instanceName,
		InstanceId:   instanceId,
		EtcdCluster:  etcdCluster,
	}

	endpoints := strings.Split(client.EtcdCluster, ",")

	cli, err := NewEtcdClient(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
		DialOptions: []grpc.DialOption{grpc.WithInsecure()},
	})

	// etcd clientv3 >= v3.2.10, grpc/grpc-go >= v1.7.3
	if err == context.DeadlineExceeded {
		return nil, ErrTimeout.New()
	}
	// etcd clientv3 <= v3.2.9, grpc/grpc-go <= v1.2.1
	if err == context.DeadlineExceeded {
		return nil, ErrTimeout.New()
	}

	if nil != err {
		return nil, ErrGrpc.New(err)
	}
	client.Client = cli

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	client.GrpcContext = ctx
	client.GrpcContextClear = cancel
	return client, nil
}

func (client *ConfigClient) Close() {
	client.m.Lock()
	defer client.m.Unlock()

	if nil != client.GrpcContextClear {
		client.GrpcContextClear()
		client.GrpcContextClear = nil
	}

	if nil != client.Client {
		client.Client.Close()
	}
}
func (client *ConfigClient) GetTreeKeys(key string) []string {
	return []string{
		CONFIG_PREFIX + "/" + key,
		CONFIG_PREFIX + "/" + client.InstanceName + "/" + key,
		CONFIG_PREFIX + "/" + client.InstanceName + "/" + client.InstanceId + "/" + key,
	}
}
func (client *ConfigClient) innerGet(key string) ([]string, uferror.UFError) {
	// ret := ""
	keys := client.GetTreeKeys(key)

	ret := []string{}
	for _, k := range keys {
		resp, err := client.Client.Get(client.GrpcContext, k)
		if nil != err {
			return nil, ErrGrpc.New(err)
		}

		str := ""
		if len(resp.Kvs) > 0 {
			str = string(resp.Kvs[0].Value)
		}
		ret = append(ret, str)
	}
	return ret, nil
}

func (client *ConfigClient) innerPicker(arr []string) (int, string, uferror.UFError) {
	for i := len(arr) - 1; i >= 0; i-- {
		if len(arr[i]) > 0 {
			return i, arr[i], nil
		}
	}
	return -1, "", ErrNotFound.New()
}

func (client *ConfigClient) Get(key string) (string, uferror.UFError) {
	arr, e := client.innerGet(key)
	if e != nil {
		return "", nil
	}
	_, s, e := client.innerPicker(arr)
	return s, e
}

func (client *ConfigClient) Watch(key string) (<-chan string, uferror.UFError) {
	vals, err := client.innerGet(key)
	if nil != err {
		return nil, ErrGrpc.New(err)
	}
	oldLevel, oldValue, e := client.innerPicker(vals)
	if nil != e {
		return nil, e
	}

	ch := make(chan string)
	keys := client.GetTreeKeys(key)
	chResps := make([]reflect.SelectCase, len(keys))

	// watch 多级目录
	for i, k := range keys {
		respCh := client.Client.Watch(context.TODO(), k)
		chResps[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(respCh)}
	}

	go func() {
		for {
			chosen, value, ok := reflect.Select(chResps)
			if !ok {
				break
			}

			resp := value.Interface().(clientv3.WatchResponse)
			if nil != resp.Err() {
				log.Println("watch error:", resp.Err())
				continue
			}
			// 如果是删除，需要更新vals，否则更新缓存值
			for _, ev := range resp.Events {
				if ev.Type == clientv3.EventTypeDelete {
					vals[chosen] = ""
				} else {
					vals[chosen] = string(ev.Kv.Value)
				}
			}

			// 忽略上级目录的修改
			if chosen < oldLevel {
				continue
			}

			for _, ev := range resp.Events {
				// 如果删除当前节点，则需要向上追溯
				if ev.Type == clientv3.EventTypeDelete {
					newValue := ""
					newLevel := -1
					for i := oldLevel - 1; i >= 0; i-- {
						if len(vals[i]) > 0 {
							newValue = vals[i]
							newLevel = i
							break
						}
					}
					// 如果是-1,说明上级都没有设置，返回空字符串后，之后任何级别的设置都会触发配置项更新
					oldLevel = newLevel
					if newValue != oldValue {
						oldValue = newValue
						ch <- newValue
					}
				} else {
					newValue := string(ev.Kv.Value)
					oldLevel = chosen
					if newValue != oldValue {
						// 更新为新的变量
						oldValue = newValue
						ch <- newValue
					}
				}
			}
		}
	}()

	return ch, nil
}

func (client *ConfigClient) GetInt(key string) (int64, uferror.UFError) {
	str, uferr := client.Get(key)
	if nil != uferr {
		return 0, uferr
	}

	i, err := strconv.ParseInt(str, 10, 64)
	if nil != err {
		return 0, ErrTypeNotMatch.New(err)
	}
	return i, nil
}

func (client *ConfigClient) GetWithDefault(key, defaultValue string) string {
	ret, uferr := client.Get(key)
	if nil != uferr {
		return defaultValue
	} else {
		return ret
	}
}

func (client *ConfigClient) GetIntWithDefault(key string, defaultValue int64) int64 {
	ret, uferr := client.GetInt(key)
	if nil != uferr {
		return defaultValue
	} else {
		return ret
	}
}
