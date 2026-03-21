package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"seckill2_go/config"
	"seckill2_go/model"
	"time"
)

// DB 全局数据库连接实例
var DB *gorm.DB

// Init 初始化数据库连接
func Init() error {
	// 1. 从配置中获取数据库连接信息
	dbConfig := config.GetConfig().Database
	// 2. 配置 GORM 日志级别
	var logLevel logger.LogLevel
	switch dbConfig.LogLevel {
	case "silent":
		logLevel = logger.Silent
	case "error":
		logLevel = logger.Error
	case "warn":
		logLevel = logger.Warn
	case "info":
		logLevel = logger.Info
	default:
		logLevel = logger.Warn // 默认警告级别
	}
	// 3. 初始化 GORM DB 实例
	var err error
	DB, err = gorm.Open(mysql.Open(dbConfig.DSN), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return err
	}
	// 配置数据库连接池
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	// 设置最大打开连接数
	sqlDB.SetMaxOpenConns(dbConfig.MaxOpenConns)
	// 设置最大空闲连接数
	sqlDB.SetMaxIdleConns(dbConfig.MaxIdleConns)
	// 设置连接的最大存活时间
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 自动迁移数据表结构（根据模型定义创建或更新表）
	if err := autoMigrate(); err != nil {
		return err
	}

	log.Println("数据库初始化成功")
	return nil
}

// 自动迁移数据表
func autoMigrate() error {
	// 在这里添加需要自动迁移的模型
	return DB.AutoMigrate(
		&model.User{}, // 用户表
		&model.Product{},
		// 可以添加更多模型...
	)
}

// Close 关闭数据库连接
func Close() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
