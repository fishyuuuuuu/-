package service

import (
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"seckill2_go/model"
	"seckill2_go/rabbitmq"
	"strconv"
	"strings"
)

// OrderMessage 定义订单消息结构
type OrderMessage struct {
	UserID    uint64 // 订单ID
	ProductID uint64 // 商品ID，匹配 model.Product.ID
}

// OrderService 订单服务
type OrderService struct {
	consumer *rabbitmq.Consumer
	db       *gorm.DB
	logger   *zap.Logger
}

// NewOrderService 创建订单服务
func NewOrderService(logger *zap.Logger, db *gorm.DB) (*OrderService, error) {
	consumer, err := rabbitmq.NewConsumer(logger)
	if err != nil {
		logger.Error("无法创建 RabbitMQ 消费者", zap.Error(err))
		return nil, fmt.Errorf("无法创建 RabbitMQ 消费者: %v", err)
	}

	return &OrderService{
		consumer: consumer,
		db:       db,
		logger:   logger,
	}, nil
}

// StartConsumingOrders 开始消费订单消息
func (s *OrderService) StartConsumingOrders(queueName, consumerTag string) error {
	// 从队列消费消息
	msgs, err := s.consumer.ConsumeFromQueue(queueName, consumerTag)
	if err != nil {
		s.logger.Error("无法启动订单消息消费", zap.String("queue", queueName), zap.Error(err))
		return fmt.Errorf("无法启动订单消息消费: %v", err)
	}

	// 在 Goroutine 中处理消息
	go func() {
		for d := range msgs {
			// 解析消息：格式为 user_id:product_id
			body := string(d.Body)
			parts := strings.Split(body, ":")
			if len(parts) != 2 {
				s.logger.Error("消息格式无效，期望 order_id:product_id",
					zap.String("body", body))
				// 可选择重新入队：d.Nack(true)
				continue
			}

			// 转换为 OrderMessage
			msg := OrderMessage{
				UserID:    parseString2Uint64(parts[0], s.logger),
				ProductID: parseString2Uint64(parts[1], s.logger),
			}

			s.logger.Info("收到订单消息",
				zap.Uint64("user_id", msg.UserID),
				zap.Uint64("product_id", msg.ProductID))

			// 更新商品库存
			if err := s.updateProductStock(msg.ProductID, msg.UserID); err != nil {
				s.logger.Error("更新商品库存失败",
					zap.Uint64("user_id", msg.UserID),
					zap.Uint64("product_id", msg.ProductID),
					zap.Error(err))
				// 可选择重新入队：d.Nack(true)
				continue
			}

			// 手动确认消息
			if err := d.Ack(false); err != nil {
				s.logger.Error("消息确认失败",
					zap.Uint64("user_id", msg.UserID),
					zap.Uint64("product_id", msg.ProductID),
					zap.Error(err))
			} else {
				s.logger.Info("订单消息处理完成",
					zap.Uint64("order_id", msg.UserID),
					zap.Uint64("product_id", msg.ProductID))
			}
		}

		s.logger.Warn("消息通道已关闭", zap.String("queue", queueName))
		// 尝试重连
		if err := s.consumer.Reconnect(); err != nil {
			s.logger.Error("重连失败", zap.Error(err))
		} else {
			s.logger.Info("重连成功，重新启动订单消费者")
			if err := s.StartConsumingOrders(queueName, consumerTag); err != nil {
				s.logger.Error("重新启动订单消费者失败", zap.Error(err))
			}
		}
	}()

	s.logger.Info("订单消费者已启动，等待消息...", zap.String("queue", queueName))
	return nil
}

// parseString2Uint64 将字符串转换为 uint64，处理转换错误
func parseString2Uint64(str string, logger *zap.Logger) uint64 {
	ui64, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		logger.Error("解析字符串失败，设置为默认值 1",
			zap.String("str", str),
			zap.Error(err))
		return 1 // 默认值
	}
	return ui64
}

// updateProductStock 更新商品库存
func (s *OrderService) updateProductStock(productID uint64, userID uint64) error {
	var product model.Product
	// 查找商品
	if err := s.db.Where("id = ? AND deleted_at IS NULL", productID).First(&product).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("商品不存在: product_id=%d", productID)
		}
		return fmt.Errorf("查询商品失败: %v", err)
	}

	// 检查库存
	if product.Stock <= 0 {
		return fmt.Errorf("库存不足: product_id=%d, stock=%d", productID, product.Stock)
	}

	// 直接更新库存
	result := s.db.Model(&model.Product{}).
		Where("id = ? AND deleted_at IS NULL", productID).
		Update("stock", product.Stock-1)

	if result.Error != nil {
		return fmt.Errorf("更新库存失败: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("库存更新失败，可能商品已被删除: product_id=%d", productID)
	}

	var orderAmount float64
	if err := s.db.Model(&model.Product{}).Select("price").Where("id = ?", productID).First(&orderAmount).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("商品不存在: id=%d", productID)
		}
		return fmt.Errorf("查询订单金额失败: %v", err)
	}
	s.logger.Info("订单金额", zap.Float64("order_amount", orderAmount))

	// 插入订单
	order := model.Order{
		UserID:      userID,
		ProductID:   productID,
		Status:      model.OrderStatusPending,
		OrderAmount: orderAmount,
	}

	if err := s.db.Create(&order).Error; err != nil {
		return fmt.Errorf("创建订单失败: %v", err)
	}

	s.logger.Info("订单创建成功",
		zap.Uint64("user_id", order.UserID),
		zap.Uint64("product_id", order.ProductID),
		zap.Float64("order_amount", order.OrderAmount))

	s.logger.Info("商品库存更新成功",
		zap.Uint64("product_id", productID),
		zap.Int("new_stock", product.Stock-1))
	return nil
}

// Close 关闭订单服务
func (s *OrderService) Close() error {
	s.logger.Info("正在关闭订单服务")
	if err := s.consumer.Close(); err != nil {
		s.logger.Error("关闭 RabbitMQ 消费者失败", zap.Error(err))
		return fmt.Errorf("关闭 RabbitMQ 消费者失败: %v", err)
	}
	s.logger.Info("订单服务已关闭")
	return nil
}
