package etcd

import (
	"GinTalk/settings"
	"sync"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

var (
	client *clientv3.Client
	once   sync.Once
)

// newClient 用于创建 etcd client
func newClient(endpoints []string, timeout time.Duration) {
	var err error
	// 创建 etcd client
	client, err = clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: timeout,
	})
	if err != nil {
		zap.L().Fatal("创建 etcd client 失败", zap.Error(err))
		return
	}
	zap.L().Info("创建 etcd client 成功")
	return
}

// GetClient 用于获取 etcd client
// 通过 sync.Once 来保证只初始化一次
//
// 返回值:
//   - *clientv3.Client: etcd client 指针
func GetClient() *clientv3.Client {
	once.Do(func() {
		cfg := settings.GetConfig().Etcd
		newClient(cfg.Endpoints, time.Duration(cfg.Timeout)*time.Second)
	})
	return client
}
