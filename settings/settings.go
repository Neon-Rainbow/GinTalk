package settings

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var Conf = new(Settings)

type MysqlConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DB       string `mapstructure:"db"`
}

type RedisConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type LoggerConfig struct {
	Level            string   `mapstructure:"level"`
	Format           string   `mapstructure:"format"`
	OutputPaths      []string `mapstructure:"outputPaths"`
	ErrorOutputPaths []string `mapstructure:"errorOutputPaths"`
	MaxSize          int      `mapstructure:"maxSize"`
	MaxBackups       int      `mapstructure:"maxBackups"`
	MaxAge           int      `mapstructure:"maxAge"`
	Compress         bool     `mapstructure:"compress"`
}

type Settings struct {
	Host           string `mapstructure:"host"`
	Port           int    `mapstructure:"port"`
	Timeout        int    `mapstructure:"timeout"`
	PasswordSecret string `mapstructure:"password_secret"`
	*MysqlConfig   `mapstructure:"mysql"`
	*RedisConfig   `mapstructure:"redis"`
	*LoggerConfig  `mapstructure:"logger"`
}

// Init 用于初始化配置文件
func Init() error {
	viper.SetConfigFile("./conf/config.yaml")

	// 设置mysql和redis的默认端口和host
	viper.SetDefault("mysql.host", "localhost")
	viper.SetDefault("mysql.port", 3306)
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", 6379)
	viper.SetDefault("port", 8080)
	viper.SetDefault("host", "localhost")
	viper.SetDefault("logger.level", "debug")
	viper.SetDefault("timeout", 10)

	// 用于判断配置文件是否被修改
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})

	// 读取配置文件
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("ReadInConfig failed, err: %v", err))
	}
	if err := viper.Unmarshal(&Conf); err != nil {
		panic(fmt.Errorf("unmarshal to Conf failed, err:%v", err))
	}
	return err
}
