package main

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log"
	"seckill2_go/config"
	"seckill2_go/db"
	"seckill2_go/rabbitmq"
	"seckill2_go/service"
	"seckill2_go/utils"
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
	defer func() {
		if err := utils.Logger.Sync(); err != nil {
			utils.Logger.Warn("日志刷新失败", zap.Error(err))
		}
	}()

	// 3. 初始化数据库
	if err := db.Init(); err != nil {
		utils.Logger.Fatal("初始化数据库失败", zap.Error(err))
	}
	defer func() {
		if err := db.Close(); err != nil {
			utils.Logger.Warn("数据库连接关闭失败", zap.Error(err))
		} else {
			utils.Logger.Info("数据库连接已正常关闭")
		}
	}()

	// 4. 初始化 RabbitMQ
	if err := rabbitmq.Init(); err != nil {
		log.Fatalf("RabbitMQ初始化失败: %v", err)
	}
	defer func() {
		rmq, _ := rabbitmq.GetInstance()
		if rmq != nil {
			if err := rmq.Close(); err != nil {
				utils.Logger.Warn("RabbitMQ 连接关闭失败", zap.Error(err))
			} else {
				utils.Logger.Info("RabbitMQ 连接已正常关闭")
			}
		}
	}()

	// 5. 初始化 OrderService
	orderService, err := service.NewOrderService(utils.Logger, db.DB)
	if err != nil {
		utils.Logger.Fatal("创建订单服务失败", zap.Error(err))
	}
	defer func() {
		if err := orderService.Close(); err != nil {
			utils.Logger.Warn("订单服务关闭失败", zap.Error(err))
		} else {
			utils.Logger.Info("订单服务已正常关闭")
		}
	}()

	// 6. 启动订单消息消费
	queueName := "seckill_queue"
	consumerTag := "order_consumer"
	if err := orderService.StartConsumingOrders(queueName, consumerTag); err != nil {
		utils.Logger.Fatal("启动订单消费者失败", zap.Error(err))
	}

	// 7. 初始化 Gin 服务（仅用于保持服务运行）
	r := gin.Default()
	utils.Logger.Info("服务启动，监听端口", zap.String("port", ":8081"))
	if err := r.Run(":8081"); err != nil {
		utils.Logger.Fatal("服务启动失败", zap.Error(err))
	}
}
