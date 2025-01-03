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

// HttpCountRequest http请求指标
var HttpCountRequest = NewHttpRequestMetrics()

type HttpRequestMetrics struct {
	httpRequestCounter *prometheus.CounterVec
}

// NewHttpRequestMetrics 创建http请求指标
func NewHttpRequestMetrics() *HttpRequestMetrics {
	httpRequestCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: ProjectNameSpace,
		Subsystem: "http",
		Name:      "request",
		Help:      "HTTP 请求次数",
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
//	metrics.HttpCountRequest.AddCounter("GET", "/ping", "200")
func (m *HttpRequestMetrics) AddCounter(method, path, status string) {
	m.httpRequestCounter.WithLabelValues(method, path, status).Add(1)
}

var HttpDuration = NewHttpDurationMetrics()

type HttpDurationMetrics struct {
	httpDurationHistogram *prometheus.HistogramVec
}

// NewHttpDurationMetrics 创建http请求耗时指标
func NewHttpDurationMetrics() *HttpDurationMetrics {
	httpDurationHistogram := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: ProjectNameSpace,
		Subsystem: "http",
		Name:      "duration",
		Help:      "HTTP 请求耗时",
		Buckets:   prometheus.DefBuckets,
	}, []string{"method", "path", "status"})

	prometheus.MustRegister(httpDurationHistogram)
	return &HttpDurationMetrics{
		httpDurationHistogram: httpDurationHistogram,
	}
}

func (m *HttpDurationMetrics) AddHistogram(method, path, status string, duration float64) {
	m.httpDurationHistogram.WithLabelValues(method, path, status).Observe(duration)
}
