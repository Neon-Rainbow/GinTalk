package MySQL

import (
	"GinTalk/settings"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"sync"
	"time"
)

var (
	db   *gorm.DB
	once sync.Once // 确保单例模式
)

// initDB 初始化数据库连接
func initDB(config *settings.MysqlConfig) (err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.DB,
	)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // 日志输出到标准输出
		logger.Config{
			SlowThreshold:             time.Duration(config.SlowThreshold) * time.Millisecond, // 慢 SQL 阈值
			LogLevel:                  logger.LogLevel(config.MySQLLoggerConfig.LogLevel),     // 级别
			IgnoreRecordNotFoundError: true,                                                   // 忽略 ErrRecordNotFound 错误
			Colorful:                  config.MySQLLoggerConfig.Colorful,                      // 彩色打印
			ParameterizedQueries:      config.MySQLLoggerConfig.ParameterizedQueries,          // 参数化查询
		},
	)

	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger:                 newLogger,
		SkipDefaultTransaction: false,
		PrepareStmt:            false,
	})
	if err != nil {
		return err
	}

	return nil
}

// Close 关闭数据库连接
func Close() {
	sqlDB, _ := db.DB()
	err := sqlDB.Close()
	if err != nil {
		return
	}
}

// GetDB 获取数据库连接
// 使用单例模式，确保数据库连接只初始化一次
func GetDB() *gorm.DB {
	// 使用 sync.Once 确保数据库只初始化一次
	once.Do(func() {
		cfg := settings.GetConfig().MysqlConfig
		if err := initDB(cfg); err != nil {
			zap.L().Fatal("initDB failed", zap.Error(err))
		}
	})

	return db
}
