package main

import (
	"GinTalk/dao/MySQL"
	"GinTalk/dao/Redis"
	"GinTalk/etcd"
	"GinTalk/kafka"
	"GinTalk/logger"
	"GinTalk/pkg/snowflake"
	"GinTalk/router"
	"GinTalk/settings"
	"fmt"
)

func main() {

	// 设置机器号
	snowflake.SetMachineID(1)

	if err := logger.SetupGlobalLogger(settings.GetConfig().LoggerConfig); err != nil {
		fmt.Printf("初始化日志库失败,错误原因: %v\n", err)
	}

	// 初始化配置
	kafka.InitKafkaManager()
	defer kafka.GetKafkaManager().Close()

	defer MySQL.Close()
	defer Redis.Close()

	if err := etcd.GetService().Register(); err != nil {
		fmt.Printf("注册服务失败,错误原因: %v\n", err)
	}

	// 初始化路由
	r := router.SetupRouter()
	err := r.Run(fmt.Sprintf("%s:%d", settings.GetConfig().Host, settings.GetConfig().Port))
	if err != nil {
		fmt.Printf("启动失败,错误原因: %v\n", err)
	}
}
