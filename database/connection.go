package database

import (
	"app_api/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Init 修改后的Init函数
func Init() error {
	// 从全局配置获取参数
	dbConfig := config.Get().Database

	db, err := gorm.Open(mysql.Open(dbConfig.DSN), &gorm.Config{})
	if err != nil {
		return err
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(dbConfig.MaxIdleConns)
	sqlDB.SetMaxOpenConns(dbConfig.MaxOpenConns)

	DB = db
	return nil
}

func GetDB() *gorm.DB {
	return DB
}
