package etcd

import (
	"GinTalk/settings"
	"context"
	"encoding/json"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"sync"
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

type Options interface {
	apply(*Service)
}

type idOption string

func (id idOption) apply(o *Service) {
	o.ID = string(id)
}

func WithID(id string) Options {
	return idOption(id)
}

type nameOption string

func (name nameOption) apply(o *Service) {
	o.Name = string(name)
}

func WithName(name string) Options {
	return nameOption(name)
}

type hostOption string

func (host hostOption) apply(o *Service) {
	o.Host = string(host)
}

func WithHost(host string) Options {
	return hostOption(host)
}

type portOption int

func (port portOption) apply(o *Service) {
	o.Port = int(port)
}

func WithPort(port int) Options {
	return portOption(port)
}

type leaseTimeOption int64

func (leaseTime leaseTimeOption) apply(o *Service) {
	o.LeaseTime = int64(leaseTime)
}

func WithLeaseTime(leaseTime int64) Options {
	return leaseTimeOption(leaseTime)
}

type intervalOption int64

func (interval intervalOption) apply(o *Service) {
	o.Interval = int64(interval)
}

func WithInterval(interval int64) Options {
	return intervalOption(interval)
}

type timeoutOption int64

func (timeout timeoutOption) apply(o *Service) {
	o.Timeout = int64(timeout)
}

func WithTimeout(timeout int64) Options {
	return timeoutOption(timeout)
}

type deregisterAfterOption int64

func (deregisterAfter deregisterAfterOption) apply(o *Service) {
	o.DeregisterAfter = int64(deregisterAfter)
}

func WithDeregisterAfter(deregisterAfter int64) Options {
	return deregisterAfterOption(deregisterAfter)
}

type configOption struct {
	settings.ServiceRegistry
}

func (c configOption) apply(o *Service) {
	o.ID = c.ID
	o.Name = c.Name
	o.Host = c.Host
	o.Port = c.Port
	o.LeaseTime = c.LeaseTime
	o.Interval = c.Interval
	o.Timeout = c.Timeout
	o.DeregisterAfter = c.DeregisterAfter
}

func WithConfig(c settings.ServiceRegistry) Options {
	return configOption{c}
}

var (
	service            *Service
	newEtcdServiceOnce sync.Once
)

// newService 用于注册服务
// 参数:
//   - options: 服务配置
//
// 返回值:
//   - Service: 服务
//
// 示例:
//
//	service := etcd.Register(etcd.WithID("test"), etcd.WithName("test"), etcd.WithHost(" localhost"), etcd.WithPort(8080))
//	service := etcd.Register(etcd.WithConfig(settings.GetConfig().ServiceRegistry))
func newService(options ...Options) *Service {
	s := &Service{}
	for _, option := range options {
		option.apply(s)
	}
	return s
}

func GetService() *Service {
	newEtcdServiceOnce.Do(func() {
		service = newService(
			WithConfig(*settings.GetConfig().ServiceRegistry),
		)
	})
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
