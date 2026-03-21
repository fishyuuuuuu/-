package utils

import (
	"go.uber.org/zap"
)

var Logger *zap.Logger

// InitLogger 初始化 Zap 日志
func InitLogger() error {
	var err error
	// 使用生产环境配置，输出 JSON 格式日志
	Logger, err = zap.NewProduction()
	if err != nil {
		return err
	}
	Logger.Info("zap日志库初始化完成")
	return nil
}
