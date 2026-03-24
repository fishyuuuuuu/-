package main

import (
	"context"
	"log"
	"os"
	"seckill_go/config"
	"seckill_go/db"
	"seckill_go/rabbitmq"
	"seckill_go/router"
	"seckill_go/service"
	"seckill_go/utils"
	"time"

	"go.uber.org/zap"
)

func main() {
	// 1. 读取配置文件
	if err := config.Init(); err != nil {
		log.Fatalf("初始化配置文件失败: %v", err)
	}
	// 2. 初始化日志
	if err := utils.InitLogger(); err != nil {
		log.Fatalf("初始化日志失败: %v", err)
	}
	// 日志刷新
	defer func() {
		if err := utils.Logger.Sync(); err != nil {
			utils.Logger.Warn("日志刷新失败", zap.Error(err))
		}
	}()
	// 3. 初始化数据库
	if err := db.Init(); err != nil {
		utils.Logger.Warn("初始化数据库失败，将继续运行", zap.Error(err))
	} else {
		// 初始化默认角色和权限
		if err := service.InitDefaultRoles(); err != nil {
			utils.Logger.Warn("初始化默认角色和权限失败", zap.Error(err))
		} else {
			utils.Logger.Info("默认角色和权限初始化成功")
		}
		// 初始化默认管理员用户
		if err := service.InitDefaultAdminUser(); err != nil {
			utils.Logger.Warn("初始化默认管理员用户失败", zap.Error(err))
		} else {
			utils.Logger.Info("默认管理员用户初始化成功")
		}
		// 数据库连接关闭
		defer func() {
			if err := db.Close(); err != nil {
				utils.Logger.Warn("数据库连接关闭失败", zap.Error(err))
			} else {
				utils.Logger.Info("数据库连接已正常关闭")
			}
		}()
	}
	// 4. 初始化 Redis
	if err := db.InitRedis(); err != nil {
		utils.Logger.Warn("初始化 Redis 失败，将继续运行", zap.Error(err))
	} else {
		utils.SetRedisClient(db.Redis)
		utils.Logger.Info("Redis客户端已注入限流器")

		// 设置安全模块的Redis客户端
		utils.SetSecurityRedis(db.Redis)
		utils.Logger.Info("安全模块Redis客户端已设置")

		// 初始化库存服务并预热库存
		stockService := service.GetStockService()
		if err := stockService.PreheatAllStocks(context.Background()); err != nil {
			utils.Logger.Warn("库存预热失败", zap.Error(err))
		} else {
			utils.Logger.Info("库存预热成功")
		}

		// 启动库存对账任务（每5分钟执行一次）
		stockService.StartReconciliationTask(5 * time.Minute)

		// Redis 连接关闭
		defer func() {
			if err := db.CloseRedis(); err != nil {
				utils.Logger.Warn("Redis 连接关闭失败", zap.Error(err))
			} else {
				utils.Logger.Info("Redis 连接已正常关闭")
			}
		}()
	}
	// 5. 初始化RabbitMQ
	if err := rabbitmq.Init(); err != nil {
		utils.Logger.Warn("RabbitMQ初始化失败，将继续运行", zap.Error(err))
	} else {
		// 启动订单消费者
		rmq, err := rabbitmq.GetInstance()
		if err == nil {
			err = rmq.ConsumeMessages("seckill_order_queue", rabbitmq.HandleOrderMessage)
			if err != nil {
				utils.Logger.Warn("启动订单消费者失败", zap.Error(err))
			} else {
				utils.Logger.Info("订单消费者已启动")
			}
		}

		// 初始化订单消息服务（启动消息重试任务）
		service.InitOrderMessageService()

		defer func() {
			// rabbitMQ 关闭连接
			rmq, _ := rabbitmq.GetInstance()
			if rmq != nil {
				rmq.Close()
			}
		}()
	}
	// 创建必要的目录
	uploadDirs := []string{"./uploads/images", "./temp"}
	for _, dir := range uploadDirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			utils.Logger.Warn("创建目录失败", zap.String("dir", dir), zap.Error(err))
		} else {
			utils.Logger.Info("目录创建成功", zap.String("dir", dir))
		}
	}

	// 6. 初始化路由
	r := router.Init()
	// 启动性能指标服务器
	utils.StartMetricsServer(8082)

	// 7. 启动服务
	utils.Logger.Info("服务启动，监听端口 :8081")
	if err := r.Run(":8081"); err != nil {
		utils.Logger.Fatal("服务启动失败", zap.Error(err))
	}
}
