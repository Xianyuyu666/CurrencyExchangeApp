package config

import (
	"ExchangeApp/global"
	"ExchangeApp/models"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 初始化数据库函数:负责获取GORM实例挂到全局变量，设置数据库连接参数
func initDB() {
	// 获取dsn
	dsn := AppConfig.Database.Dsn

	// 连接数据库,获取ORM数据库句柄
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Open database error: %v", err)
	}

	// 获取数据库连接池句柄，设置连接参数
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Get sql.DB failed: %v", err)
	}
	sqlDB.SetMaxIdleConns(AppConfig.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(AppConfig.Database.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 验证数据库连通性，因为惰性连接不会在Open时真正连接
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Database unreachable: %v", err)
	}

	// 数据库自动建表
	if err := db.AutoMigrate(&models.User{});err != nil{
		log.Fatalf("AutoMigrate error: %v",err)
	}

	// 将GORM实列挂到全局变量，供其他模块使用
	global.Db = db
}
