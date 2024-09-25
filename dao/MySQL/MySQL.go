package MySQL

import (
	"fmt"
	"forum-gin/model"
	"forum-gin/settings"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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

	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
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
