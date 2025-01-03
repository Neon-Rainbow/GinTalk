package controller

import (
	"GinTalk/metrics"
	"github.com/gin-gonic/gin"
	"strconv"
)

// PrometheusMiddleware Prometheus 中间件
// 用于统计接口访问次数, 并将数据上报到 Prometheus
func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.FullPath()

		method := c.Request.Method

		c.Next()

		statusCode := c.Writer.Status()

		// 统计接口访问次数
		metrics.HttpCountRequest.AddCounter(method, path, strconv.Itoa(statusCode))

		// 统计接口访问耗时
		metrics.HttpDuration.AddHistogram(method, path, strconv.Itoa(statusCode), c.GetFloat64("cost"))
	}
}
