# GinTalk
使用Gin框架搭建的论坛项目

## 项目部署:
项目的MySQL以及Redis使用了Docker部署
项目使用docker部署,docker具体可以修改[docker-compose.yml](docker-compose.yml)

使用下列指令启动所有在 docker-compose.yml 文件中定义的服务
```shell
docker-compose up
```

项目的具体配置信息在[配置文件](./conf/config.yaml)中,具体配置需求可以修改该文件

## 运行项目
```shell
go run main.go
```

## 项目日志
项目的日志内容由配置文件决定,日志库使用了[zap日志库](https://github.com/uber-go/zap)