package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

// AppConfig 配置结构的
type AppConfig struct {
	Database DatabaseConfig `yaml:"database"` // Database 数据库配置
	Redis    RedisConfig    `yaml:"redis"`    // Redis 配置
	RabbitMQ RabbitMQConfig `yaml:"rabbitmq"` // RabbitMQ 配置
}

// DatabaseConfig 数据库子配置
type DatabaseConfig struct {
	DSN          string `yaml:"dsn"`            // 数据库连接字符串
	MaxOpenConns int    `yaml:"max_open_conns"` // 最大打开连接数
	MaxIdleConns int    `yaml:"max_idle_conns"` // 最大空闲连接数
	LogLevel     string `yaml:"log_level"`      // 日志级别：silent/error/warn/info
}

// RedisConfig Redis 子配置
type RedisConfig struct {
	Addr     string `yaml:"addr"`      // Redis 地址
	Password string `yaml:"password"`  // Redis 密码
	DB       int    `yaml:"db"`        // Redis 数据库编号
	PoolSize int    `yaml:"pool_size"` // 连接池大小
}

// RabbitMQConfig RabbitMQ 子配置
type RabbitMQConfig struct {
	Host     string `yaml:"host"`     // RabbitMQ 地址
	Port     int    `yaml:"port"`     // RabbitMQ 端口
	Username string `yaml:"username"` // RabbitMQ 用户名
	Password string `yaml:"password"` // RabbitMQ 密码
}

var appConfig AppConfig

// Init 加载配置文件到 appConfig 变量
func Init() error {
	configPath := "config.yaml"
	// 读取配置文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		return err // 返回错误，让调用者处理
	}
	// 2. 解析 YAML 到 appConfig 变量
	if err := yaml.Unmarshal(data, &appConfig); err != nil {
		return err
	}
	return nil
}

// GetConfig 提供对外的“只读访问方法”，避免直接暴露全局变量
func GetConfig() AppConfig {
	return appConfig
}
