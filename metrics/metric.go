package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	ProjectNameSpace = "GinTalk"
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
		Namespace: ProjectNameSpace,
		Subsystem: "cpu",
		Name:      "usage",
		Help:      "CPU 使用率",
	}, []string{"instance"})

	// memoryUsageGauge 内存使用率
	memoryUsageGauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: ProjectNameSpace,
		Subsystem: "memory",
		Name:      "usage",
		Help:      "内存使用率",
	}, []string{"instance"})

	// goroutineGauge goroutine数量
	goroutineGauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: ProjectNameSpace,
		Subsystem: "goroutine",
		Name:      "num",
		Help:      "goroutine数量",
	}, []string{"instance"})

	// processNumGauge 进程数量
	processNumGauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: ProjectNameSpace,
		Subsystem: "process",
		Name:      "num",
		Help:      "进程数量",
	}, []string{"instance"})

	m := &metrics{
		cpuUsageGauge:    cpuUsageGauge,
		memoryUsageGauge: memoryUsageGauge,
		goroutineGauge:   goroutineGauge,
		processNumGauge:  processNumGauge,
	}

	// 注册指标
	prometheus.MustRegister(cpuUsageGauge, memoryUsageGauge, goroutineGauge, processNumGauge)
	return m
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
