package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

// metrics 指标
// 其中包含了CPU使用率, 内存使用率, goroutine数量, 进程数量, 自定义计数器
type metrics struct {
	// cpuUsageGauge CPU使用率
	cpuUsageGauge *prometheus.GaugeVec

	// memoryUsageGauge 内存使用率
	memoryUsageGauge *prometheus.GaugeVec

	// goroutineGauge goroutine数量
	goroutineGauge *prometheus.GaugeVec

	// processNumGauge 进程数量
	processNumGauge *prometheus.GaugeVec
}

func NewMetrics() *metrics {

	// cpuUsageGauge CPU使用率
	cpuUsageGauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "monitor",
		Subsystem: "cpu",
		Name:      "usage",
		Help:      "The percentage of CPU usage",
	}, []string{"instance"})

	// memoryUsageGauge 内存使用率
	memoryUsageGauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "monitor",
		Subsystem: "memory",
		Name:      "usage",
		Help:      "The percentage of memory usage",
	}, []string{"instance"})

	// goroutineGauge goroutine数量
	goroutineGauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "monitor",
		Subsystem: "goroutine",
		Name:      "num",
		Help:      "The number of goroutines that currently exist",
	}, []string{"instance"})

	// processNumGauge 进程数量
	processNumGauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "monitor",
		Subsystem: "process",
		Name:      "num",
		Help:      "The number of processes that currently exist",
	}, []string{"instance"})

	customCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "monitor",
		Subsystem: "custom",
		Name:      "num",
		Help:      "The number of custom counter",
	}, []string{"instance"})

	m := &metrics{
		cpuUsageGauge:    cpuUsageGauge,
		memoryUsageGauge: memoryUsageGauge,
		goroutineGauge:   goroutineGauge,
		processNumGauge:  processNumGauge,
	}

	// 注册指标
	prometheus.MustRegister(cpuUsageGauge, memoryUsageGauge, goroutineGauge, processNumGauge, customCounter)
	return m
}

// updateCpuUsage 更新CPU使用率
//
// 参数:
//   - instance: 实例名称
//   - value: CPU使用率
func (m *metrics) updateCpuUsage(instance string, value float64) {
	m.cpuUsageGauge.WithLabelValues(instance).Set(value)
}

// updateMemoryUsage 更新内存使用率
//
// 参数:
//   - instance: 实例名称
//   - value: 内存使用率
func (m *metrics) updateMemoryUsage(instance string, value float64) {
	m.memoryUsageGauge.WithLabelValues(instance).Set(value)
}

// updateGoroutineNum 更新goroutine数量
//
// 参数:
//   - instance: 实例名称
//   - value: goroutine数量
func (m *metrics) updateGoroutineNum(instance string, value float64) {
	m.goroutineGauge.WithLabelValues(instance).Set(value)
}

// updateProcessNum 更新进程数量
//
// 参数:
//   - instance: 实例名称
//   - value: 进程数量
func (m *metrics) updateProcessNum(instance string, value float64) {
	m.processNumGauge.WithLabelValues(instance).Set(value)
}

// HttpRequest http请求指标
var HttpRequest = NewHttpRequestMetrics()

type HttpRequestMetrics struct {
	httpRequestCounter *prometheus.CounterVec
}

// NewHttpRequestMetrics 创建http请求指标
func NewHttpRequestMetrics() *HttpRequestMetrics {
	httpRequestCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "monitor",
		Subsystem: "http",
		Name:      "request",
		Help:      "The number of http request",
	}, []string{"method", "path", "status"})

	prometheus.MustRegister(httpRequestCounter)
	return &HttpRequestMetrics{
		httpRequestCounter: httpRequestCounter,
	}
}

// AddCounter 添加http请求计数器
//
// 参数:
//   - method: 请求方法, 如 GET, POST
//   - path: 请求路径
//   - status: 请求状态码, 如 200, 404
//
// 使用示例:
//
//	metrics.HttpRequest.AddCounter("GET", "/ping", "200")
func (m *HttpRequestMetrics) AddCounter(method, path, status string) {
	m.httpRequestCounter.WithLabelValues(method, path, status).Add(1)
}
