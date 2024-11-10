package client

//go:generate mockgen -source=etcd-client.go -destination=mocks/etcd-client-mock.go -package=mocks EtcdClientMock

import (
	"context"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type EtcdClient interface {
	Close() error
	Get(ctx context.Context, key string, opts ...clientv3.OpOption) (*clientv3.GetResponse, error)
	Watch(ctx context.Context, key string, opts ...clientv3.OpOption) clientv3.WatchChan
}

func originNewEtcdClient(cfg clientv3.Config) (EtcdClient, error) {
	return clientv3.New(cfg)
}

type NewClientFunction func(cfg clientv3.Config) (EtcdClient, error)

var NewEtcdClient NewClientFunction = originNewEtcdClient
