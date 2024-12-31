package settings

import (
	"flag"
	"fmt"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var conf = new(Settings)
var once sync.Once

type MysqlConfig struct {
	Host               string `mapstructure:"host"`
	Port               int    `mapstructure:"port"`
	User               string `mapstructure:"user"`
	Password           string `mapstructure:"password"`
	DB                 string `mapstructure:"db"`
	*MySQLLoggerConfig `mapstructure:"logger"`
}

type MySQLLoggerConfig struct {
	SlowThreshold             int  `mapstructure:"slowThreshold"`
	LogLevel                  int  `mapstructure:"logLevel"`
	IgnoreRecordNotFoundError bool `mapstructure:"ignoreRecordNotFoundError"`
	Colorful                  bool `mapstructure:"colorful"`
	ParameterizedQueries      bool `mapstructure:"parameterizedQueries"`
}

type RedisConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	DB   int    `mapstructure:"db"`
}

type LoggerConfig struct {
	Level            int      `mapstructure:"level"`
	Format           string   `mapstructure:"format"`
	OutputPaths      []string `mapstructure:"outputPaths"`
	ErrorOutputPaths []string `mapstructure:"errorOutputPaths"`
	MaxSize          int      `mapstructure:"maxSize"`
	MaxBackups       int      `mapstructure:"maxBackups"`
	MaxAge           int      `mapstructure:"maxAge"`
	Compress         bool     `mapstructure:"compress"`
}

type Etcd struct {
	Endpoints []string `mapstructure:"endpoints"`
	Timeout   int      `mapstructure:"timeout"`
}

type ServiceRegistry struct {
	ID              string `mapstructure:"id"`
	Name            string `mapstructure:"name"`
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	LeaseTime       int64  `mapstructure:"leaseTime"`
	Interval        int64  `mapstructure:"interval"`
	Timeout         int64  `mapstructure:"timeout"`
	DeregisterAfter int64  `mapstructure:"deregisterAfter"`
}

type KafkaConfig struct {
	Brokers []string `mapstructure:"brokers"`
}

type Settings struct {
	Host             string `mapstructure:"host"`
	Port             int    `mapstructure:"port"`
	Timeout          int    `mapstructure:"timeout"`
	PasswordSecret   string `mapstructure:"password_secret"`
	Mode             string `mapstructure:"mode"`
	*MysqlConfig     `mapstructure:"mysql"`
	*RedisConfig     `mapstructure:"redis"`
	*LoggerConfig    `mapstructure:"logger"`
	*Etcd            `mapstructure:"etcd"`
	*ServiceRegistry `mapstructure:"service_registry"`
	*KafkaConfig     `mapstructure:"kafka"`
}

// mustInitConfig 用于初始化配置文件
// 从命令行中读取配置文件路径
// 在初始化失败时，会触发 panic
func mustInitConfig() {

	// 从命令行中读取配置文件路径
	configFilePath := flag.String("config", "./conf/config.yaml", "配置文件路径")
	viper.SetConfigFile(*configFilePath)

	// 设置mysql和redis的默认端口和host
	viper.SetDefault("mysql.host", "localhost")
	viper.SetDefault("mysql.port", 3306)

	viper.SetDefault("mysql.logger.slowThreshold", 100)
	viper.SetDefault("mysql.logger.logLevel", 1)
	viper.SetDefault("mysql.logger.ignoreRecordNotFoundError", true)
	viper.SetDefault("mysql.logger.colorful", true)

	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", 6379)
	viper.SetDefault("redis.db", 0)

	viper.SetDefault("port", 8080)
	viper.SetDefault("host", "localhost")
	viper.SetDefault("logger.level", "debug")
	viper.SetDefault("timeout", 10)
	viper.SetDefault("mode", "release")

	// 用于判断配置文件是否被修改
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("ReadInConfig failed, err: %v", err))
	}

	// 将配置文件解析到结构体中
	if err := viper.Unmarshal(&conf); err != nil {
		panic(fmt.Errorf("unmarshal to conf failed, err:%v", err))
	}
}

// GetConfig 用于获取配置文件
// 使用单例模式，确保配置文件只被初始化一次
func GetConfig() *Settings {
	once.Do(mustInitConfig)
	return conf
}
