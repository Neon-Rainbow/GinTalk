package etcd

import (
	"GinTalk/settings"
	"context"
	"encoding/json"
	"fmt"
	"sync"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

// Service 服务
// 用于注册服务到 etcd
type Service struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Host            string `json:"host"`
	Port            int    `json:"port"`
	LeaseTime       int64  `json:"lease_time"`
	Interval        int64  `json:"interval"`
	Timeout         int64  `json:"timeout"`
	DeregisterAfter int64  `json:"deregister_after"`
}

var (
	service            *Service
	newEtcdServiceOnce sync.Once
)

// NewService 用于注册服务
// 参数:
//   - options: 服务配置
//
// 返回值:
//   - Service: 服务
//
// 示例:
//
// 	service := etcd.Register(etcd.WithID("test"), etcd.WithName("test"), etcd.WithHost(" localhost"), etcd.WithPort(8080))
// 	service := etcd.Register(etcd.WithConfig(settings.GetConfig().ServiceRegistry))
//	service := etcd.Register()
func NewService(options ...Options) *Service {
	cfg := settings.GetConfig().ServiceRegistry
	s := &Service{
		ID:        cfg.ID,
		Name:      cfg.Name,
		Host:      cfg.Host,
		Port:      cfg.Port,
		LeaseTime: cfg.LeaseTime,
	}

	if len(options) == 0 {
		options = append(options, WithDefault())
	}
	for _, option := range options {
		option.apply(s)
	}
	service = s
	return s
}

func GetService() *Service {
	if service == nil {
		zap.L().Fatal("etcd服务未初始化")
	}
	return service
}

// Register 用于注册服务
func (s *Service) Register() error {
	// 注册服务
	// 创建租约
	leaseResp, err := GetClient().Grant(context.TODO(), s.LeaseTime)
	if err != nil {
		zap.L().Error("创建租约失败", zap.Error(err))
		return err
	}

	// 服务信息序列化为 JSON
	serviceKey := fmt.Sprintf("/service/%v/%v", s.Name, s.ID)
	serviceValue, err := json.Marshal(s)
	if err != nil {
		zap.L().Error("序列化服务信息失败", zap.Error(err))
		return err
	}

	// 注册服务信息
	_, err = GetClient().Put(context.TODO(), serviceKey, string(serviceValue), clientv3.WithLease(leaseResp.ID))
	if err != nil {
		zap.L().Error("注册服务失败", zap.Error(err))
		return err
	}

	// 开始续租
	go func(*clientv3.LeaseGrantResponse) {
		s.keepAlive(leaseResp.ID)
	}(leaseResp)

	zap.L().Info("服务注册成功", zap.String("key", serviceKey))
	return nil
}

// keepAlive 用于保持租约
func (s *Service) keepAlive(leaseID clientv3.LeaseID) {
	keepAliveChan, err := GetClient().KeepAlive(context.TODO(), leaseID)
	if err != nil {
		zap.L().Error("续租失败", zap.Error(err))
		return
	}

	for ka := range keepAliveChan {
		if ka == nil {
			zap.L().Error("续租失败")
			return
		}
	}
}
