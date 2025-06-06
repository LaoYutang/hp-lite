package db

import (
	"hp-server/internal/entity"
	"hp-server/pkg/logger"
	"os"
	"time"

	"github.com/glebarez/sqlite"
	"go.uber.org/zap"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

var DB *gorm.DB

type ExpandLogger struct {
	inner *zap.SugaredLogger
}

func (e *ExpandLogger) Printf(template string, args ...interface{}) {
	e.inner.With("module", "gorm").Infof(template, args...)
}

func init() {
	err := os.MkdirAll("./data", os.ModePerm)
	if err != nil {
		logger.Fatalf("创建目录失败: %v", err)
	}

	DB, err = gorm.Open(sqlite.Open("./data/hp-lite.db"), &gorm.Config{
		Logger: gormLogger.New(
			&ExpandLogger{inner: logger.Logger},
			gormLogger.Config{
				SlowThreshold: 200 * time.Millisecond,
				LogLevel:      gormLogger.Warn,
			}),
	})
	if err != nil {
		logger.Fatalf("数据库连接失败: %v", err)
	}

	// 获取底层的 sql.DB 实例
	sqlDB, err := DB.DB()
	if err != nil {
		logger.Fatalf("获取DB实例失败: %v", err)
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
