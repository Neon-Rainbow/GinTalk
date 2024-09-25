package main

import (
	"fmt"
	"forum-gin/dao/MySQL"
	"forum-gin/dao/Redis"
	"forum-gin/logger"
	"forum-gin/router"
	"forum-gin/settings"
)

func main() {
	if err := settings.Init(); err != nil {
		fmt.Printf("初始化配置文件失败,错误原因: %v\n", err)
	}

	if err := MySQL.Init(settings.Conf.MysqlConfig); err != nil {
		fmt.Printf("初始化数据库连接失败,错误原因: %v\n", err)
	}
	defer MySQL.Close()

	if err := Redis.Init(settings.Conf.RedisConfig); err != nil {
		fmt.Printf("初始化Redis连接失败,错误原因: %v\n", err)
	}
	defer Redis.Close()

	if err := logger.SetupGlobalLogger(settings.Conf.LoggerConfig); err != nil {
		fmt.Printf("初始化日志库失败,错误原因: %v\n", err)
	}

	// 初始化路由
	r := router.SetupRouter()
	err := r.Run(fmt.Sprintf(":%d", settings.Conf.Port))
	if err != nil {
		fmt.Printf("启动失败,错误原因: %v\n", err)
	}
}
