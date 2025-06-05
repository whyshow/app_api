package database

import (
	"app_api/config"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Init 修改后的Init函数
func Init() error {
	// 从全局配置获取参数
	dbConfig := config.Get().Database

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.DbName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
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
