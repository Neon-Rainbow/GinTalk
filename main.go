package main

import (
	"GinTalk/container"
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

	// 初始化IOC容器
	container.InitContainer()

	// 初始化路由
	r := router.SetupRouter()
	err := r.Run(fmt.Sprintf("%s:%d", settings.GetConfig().Host, settings.GetConfig().Port))
	if err != nil {
		fmt.Printf("启动失败,错误原因: %v\n", err)
	}
}
