package etcd

import "GinTalk/settings"

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
}

func WithConfig(c settings.ServiceRegistry) Options {
	return configOption{c}
}

type defaultOption struct{}

func (d defaultOption) apply(o *Service) {
	o.ID = "default_id"
	o.Name = "default_name"
	o.Host = "localhost"
	o.Port = 8080
	o.LeaseTime = 60
	o.Interval = 10
	o.Timeout = 10
	o.DeregisterAfter = 60
}

func WithDefault() Options {
	return defaultOption{}
}
