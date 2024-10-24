package main

import (
	"GinTalk/container"
	"GinTalk/dao/MySQL"
	"GinTalk/dao/Redis"
	"GinTalk/kafka"
	"GinTalk/logger"
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
		fmt.Printf("初始化日志库失败,错误原因: %v\n", err)
	}

	// 清空 Kafka 中的 Topic
	err := kafka.ResetKafkaTopic()
	if err != nil {
		zap.L().Error("kafka.ResetKafkaTopic() 错误", zap.Error(err))
	}

	// 初始化IOC容器
	c := container.BuildContainer()

	defer MySQL.Close()
	defer Redis.Close()

	// 初始化路由
	r := router.SetupRouter(c)
	err = r.Run(fmt.Sprintf("%s:%d", settings.GetConfig().Host, settings.GetConfig().Port))
	if err != nil {
		fmt.Printf("启动失败,错误原因: %v\n", err)
	}
}
