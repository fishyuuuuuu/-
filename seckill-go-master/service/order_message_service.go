package service

import (
	"context"
	"encoding/json"
	"errors"
	"seckill_go/rabbitmq"
	"seckill_go/utils"
	"sync"
	"time"

	"go.uber.org/zap"
)

var (
	ErrOrderMessageSendFail = errors.New("订单消息发送失败")
	ErrOrderMessageNotFound = errors.New("订单消息不存在")
)

type OrderMessage struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	ProductID  string    `json:"product_id"`
	Status     string    `json:"status"`
	RetryCount int       `json:"retry_count"`
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
}

type OrderMessageService struct {
	messages sync.Map
	mu       sync.RWMutex
}

var orderMessageService *OrderMessageService
var orderMessageServiceOnce sync.Once

func GetOrderMessageService() *OrderMessageService {
	orderMessageServiceOnce.Do(func() {
		orderMessageService = &OrderMessageService{
			messages: sync.Map{},
		}
	})
	return orderMessageService
}

func (s *OrderMessageService) SaveOrderMessage(ctx context.Context, msg *OrderMessage) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	msg.CreateTime = time.Now()
	msg.UpdateTime = time.Now()
	msg.Status = "pending"
	msg.RetryCount = 0

	s.messages.Store(msg.ID, msg)

	utils.Logger.Info("订单消息已保存到本地消息表",
		zap.String("messageID", msg.ID),
		zap.String("userID", msg.UserID),
		zap.String("productID", msg.ProductID))

	return nil
}

func (s *OrderMessageService) UpdateOrderMessageStatus(ctx context.Context, messageID string, status string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	value, ok := s.messages.Load(messageID)
	if !ok {
		return ErrOrderMessageNotFound
	}

	msg := value.(*OrderMessage)
	msg.Status = status
	msg.UpdateTime = time.Now()

	s.messages.Store(messageID, msg)

	utils.Logger.Info("订单消息状态已更新",
		zap.String("messageID", messageID),
		zap.String("status", status))

	return nil
}

func (s *OrderMessageService) IncrementRetryCount(ctx context.Context, messageID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	value, ok := s.messages.Load(messageID)
	if !ok {
		return ErrOrderMessageNotFound
	}

	msg := value.(*OrderMessage)
	msg.RetryCount++
	msg.UpdateTime = time.Now()

	s.messages.Store(messageID, msg)

	utils.Logger.Info("订单消息重试次数已更新",
		zap.String("messageID", messageID),
		zap.Int("retryCount", msg.RetryCount))

	return nil
}

func (s *OrderMessageService) GetPendingMessages(ctx context.Context) []*OrderMessage {
	var pendingMessages []*OrderMessage

	s.messages.Range(func(key, value interface{}) bool {
		msg := value.(*OrderMessage)
		if msg.Status == "pending" || (msg.Status == "failed" && msg.RetryCount < 3) {
			pendingMessages = append(pendingMessages, msg)
		}
		return true
	})

	return pendingMessages
}

func (s *OrderMessageService) SendOrderMessageWithRetry(ctx context.Context, msg *OrderMessage) error {
	err := s.SaveOrderMessage(ctx, msg)
	if err != nil {
		return err
	}

	err = s.sendToRabbitMQ(ctx, msg)
	if err != nil {
		utils.Logger.Error("订单消息发送失败，将等待重试",
			zap.String("messageID", msg.ID),
			zap.Error(err))
		s.UpdateOrderMessageStatus(ctx, msg.ID, "failed")
		return ErrOrderMessageSendFail
	}

	s.UpdateOrderMessageStatus(ctx, msg.ID, "sent")
	return nil
}

func (s *OrderMessageService) sendToRabbitMQ(ctx context.Context, msg *OrderMessage) error {
	rmq, err := rabbitmq.GetInstance()
	if err != nil {
		utils.Logger.Error("获取RabbitMQ实例失败", zap.Error(err))
		return err
	}

	err = rmq.DeclareQueue("seckill_order_queue")
	if err != nil {
		utils.Logger.Error("声明队列失败", zap.Error(err))
		return err
	}

	messageBytes, err := json.Marshal(msg)
	if err != nil {
		utils.Logger.Error("序列化消息失败", zap.Error(err))
		return err
	}

	err = rmq.SendToQueue("seckill_order_queue", messageBytes)
	if err != nil {
		utils.Logger.Error("发送消息到队列失败", zap.Error(err))
		return err
	}

	utils.Logger.Info("订单消息已发送到RabbitMQ",
		zap.String("messageID", msg.ID),
		zap.String("queue", "seckill_order_queue"))

	return nil
}

func (s *OrderMessageService) RetryFailedMessages(ctx context.Context) {
	pendingMessages := s.GetPendingMessages(ctx)

	utils.Logger.Info("开始重试失败的订单消息",
		zap.Int("pendingCount", len(pendingMessages)))

	for _, msg := range pendingMessages {
		if msg.RetryCount >= 3 {
			utils.Logger.Warn("订单消息重试次数已达上限，跳过",
				zap.String("messageID", msg.ID),
				zap.Int("retryCount", msg.RetryCount))
			s.UpdateOrderMessageStatus(ctx, msg.ID, "abandoned")
			continue
		}

		err := s.sendToRabbitMQ(ctx, msg)
		if err != nil {
			utils.Logger.Error("订单消息重试发送失败",
				zap.String("messageID", msg.ID),
				zap.Int("retryCount", msg.RetryCount),
				zap.Error(err))
			s.IncrementRetryCount(ctx, msg.ID)
			s.UpdateOrderMessageStatus(ctx, msg.ID, "failed")
		} else {
			utils.Logger.Info("订单消息重试发送成功",
				zap.String("messageID", msg.ID))
			s.UpdateOrderMessageStatus(ctx, msg.ID, "sent")
		}

		time.Sleep(100 * time.Millisecond)
	}
}

func (s *OrderMessageService) StartRetryTask(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
			ctx := context.Background()
			s.RetryFailedMessages(ctx)
		}
	}()

	utils.Logger.Info("订单消息重试任务已启动", zap.Duration("interval", interval))
}

func (s *OrderMessageService) GetOrderMessageStats(ctx context.Context) map[string]interface{} {
	stats := map[string]int{
		"pending":   0,
		"sent":      0,
		"failed":    0,
		"abandoned": 0,
		"total":     0,
	}

	s.messages.Range(func(key, value interface{}) bool {
		msg := value.(*OrderMessage)
		stats[msg.Status]++
		stats["total"]++
		return true
	})

	return map[string]interface{}{
		"stats":     stats,
		"timestamp": time.Now().Format(time.RFC3339),
	}
}

func CreateOrderMessage(userID, productID string) *OrderMessage {
	return &OrderMessage{
		ID:        generateMessageID(userID, productID),
		UserID:    userID,
		ProductID: productID,
		Status:    "pending",
	}
}

func generateMessageID(userID, productID string) string {
	return userID + "_" + productID + "_" + time.Now().Format("20060102150405")
}

func PublishOrderMessageReliably(ctx context.Context, userID, productID string) error {
	msgService := GetOrderMessageService()
	msg := CreateOrderMessage(userID, productID)

	return msgService.SendOrderMessageWithRetry(ctx, msg)
}

func InitOrderMessageService() {
	msgService := GetOrderMessageService()
	msgService.StartRetryTask(30 * time.Second)

	utils.Logger.Info("订单消息服务已初始化")
}
