package MySQL

import (
	"GinTalk/model"
	"GinTalk/settings"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var db *gorm.DB

func Init(config *settings.MysqlConfig) (err error) {
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
			SlowThreshold:             time.Second, // 慢 SQL 阈值
			LogLevel:                  logger.Info, // 级别
			IgnoreRecordNotFoundError: true,        // 忽略 ErrRecordNotFound 错误
			Colorful:                  false,       // 禁用彩色打印
		},
	)

	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(
		&model.User{},
		&model.Community{},
		&model.Comment{},
		&model.Post{},
	)
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

func GetDB() *gorm.DB {
	if db == nil {
		return nil
	}
	return db
}
