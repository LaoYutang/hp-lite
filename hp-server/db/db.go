package db

import (
	"hp-server-lib/entity"
	"log"
	"os"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func init() {
	err := os.MkdirAll("./data", os.ModePerm)
	if err != nil {
		log.Fatal("创建目录失败", err)
	}

	DB, err = gorm.Open(sqlite.Open("./data/hp-lite.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // 设置日志级别为 Info
	})
	if err != nil {
		log.Fatal("数据库连接失败", err)
	}

	// 获取底层的 sql.DB 实例
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("获取DB实例失败", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(10)         // 设置最大空闲连接数
	sqlDB.SetMaxOpenConns(50)         // 设置最大打开连接数
	sqlDB.SetConnMaxLifetime(30 * 60) // 设置连接的最大生命周期（单位：秒）

	//自动创建表
	DB.AutoMigrate(&entity.UserCustomEntity{})
	DB.AutoMigrate(&entity.UserDeviceEntity{})
	DB.AutoMigrate(&entity.UserConfigEntity{})
	DB.AutoMigrate(&entity.UserStatisticsEntity{})
	DB.AutoMigrate(&entity.UserDomainEntity{})
}
