package db

import (
	"context"
	"seckill_go/config"
	"seckill_go/utils"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var Redis *redis.Client

func InitRedis() error {
	cfg := config.GetConfig().Redis
	Redis = redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})

	ctx := context.Background()
	_, err := Redis.Ping(ctx).Result()
	if err != nil {
		utils.Logger.Error("Failed to connect to Redis",
			zap.String("addr", cfg.Addr),
			zap.Error(err))
		// 连接失败，将 Redis 设为 nil
		Redis = nil
		return err
	}

	utils.Logger.Info("Redis connected successfully",
		zap.String("addr", cfg.Addr),
		zap.Int("db", cfg.DB),
		zap.Int("pool_size", cfg.PoolSize))
	return nil
}

// CloseRedis 关闭 Redis 连接
func CloseRedis() error {
	if Redis == nil {
		return nil
	}
	err := Redis.Close()
	if err != nil {
		utils.Logger.Warn("Redis 连接关闭失败", zap.Error(err))
		return err
	}
	utils.Logger.Info("Redis 连接已正常关闭")
	return nil
}
