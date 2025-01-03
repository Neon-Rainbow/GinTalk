# GinTalk
使用Gin框架搭建的论坛项目

## 项目部署:
项目的MySQL,Redis,etcd,kafka,prometheus使用了Docker部署
项目使用docker部署,docker具体可以修改[docker-compose文件](docker-compose)

使用下列指令启动所有在 docker-compose.yml 文件中定义的服务
```shell
docker-compose up
```

这些服务也都可以自行部署

项目的具体配置信息在[配置文件](./conf/config.yaml)中,具体配置需求可以修改该文件

项目需要在MySQL中建表,建表语句在[该文件](./model/create_table.sql)中

## 运行项目
```shell
go run main.go
```

## 项目日志
项目的日志内容由配置文件决定,日志库使用了[zap日志库](https://github.com/uber-go/zap)

## 项目介绍
项目使用了 Gin 框架搭建了一个类似 [Reddit](https://www.reddit.com/) 的论坛项目,项目实现了完整的评论和点赞系统,这两个系统的设计如下:

评论系统:
评论系统参考了抖音,知乎的设计,使用二级模式设计,评论热度的计算公式:
$$
hot = 0.6 * 评论数 + 0.4 * 点赞数
$$

对于每一个帖子,统计出热度前100的评论,这些评论使用 ZSET 存储在 Redis 中,100条之后的评论存储在 MySQL 中,当用户请求评论时,先从 Redis 中获取前100条评论,如果 Redis 中没有,则从 MySQL 中获取

# 项目运行时数据
项目使用了prometheus来监控数据,监控的数据包括 CPU 使用率,内存使用率,goroutine数量,进程数量,接口访问次数
