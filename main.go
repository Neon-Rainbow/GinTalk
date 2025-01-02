package main

import (
	"GinTalk/dao/MySQL"
	"GinTalk/dao/Redis"
	"GinTalk/etcd"
	"GinTalk/kafka"
	"GinTalk/logger"
	"GinTalk/metrics"
	"GinTalk/pkg/snowflake"
	"GinTalk/router"
	"GinTalk/settings"
	"fmt"
	"go.uber.org/zap"
)

func main() {

	// 设置机器号
	snowflake.SetMachineID(1)

	if err := logger.SetupGlobalLogger(settings.GetConfig().LoggerConfig); err != nil {
		panic(fmt.Sprintf("初始化日志失败: %v\n", err))
	}

	// 初始化 Prometheus
	metrics.NewMetrics().AutoUpdateMetrics()

	// 初始化配置
	kafka.InitKafkaManager()

	etcd.NewService()
	if err := etcd.GetService().Register(); err != nil {
		zap.L().Fatal("注册服务失败", zap.Error(err))
	}

	defer kafka.GetKafkaManager().Close()
	defer MySQL.Close()
	defer Redis.Close()

	// 初始化路由
	r := router.SetupRouter()
	err := r.Run(fmt.Sprintf("%s:%d", settings.GetConfig().Host, settings.GetConfig().Port))
	if err != nil {
		zap.L().Fatal("启动服务失败", zap.Error(err))
	}
}
