package rabbitmq

import (
	"fmt"
	"seckill_go/utils"
	"strings"
	"time"

	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

// ConsumeMessages 消费队列消息
func (r *RabbitMQ) ConsumeMessages(queueName string, handler func([]byte) error) error {
	channel, err := r.GetChannel()
	if err != nil {
		return err
	}

	// 声明死信交换器
	dlxExchange := queueName + "_dlx"
	err = r.channel.ExchangeDeclare(
		dlxExchange, // 交换器名称
		"direct",    // 交换器类型
		true,        // 持久化
		false,       // 自动删除
		false,       // 内部交换器
		false,       // 非阻塞
		nil,         // 额外参数
	)
	if err != nil {
		return fmt.Errorf("声明死信交换器失败: %v", err)
	}

	// 声明死信队列
	dlxQueue := queueName + "_dlx_queue"
	_, err = r.channel.QueueDeclare(
		dlxQueue, // 队列名称
		true,     // 持久化
		false,    // 自动删除
		false,    // 排他性
		false,    // 非阻塞
		nil,      // 额外参数
	)
	if err != nil {
		return fmt.Errorf("声明死信队列失败: %v", err)
	}

	// 绑定死信队列到死信交换器
	err = r.channel.QueueBind(
		dlxQueue,    // 队列名称
		queueName,   // 路由键
		dlxExchange, // 交换器名称
		false,       // 非阻塞
		nil,         // 额外参数
	)
	if err != nil {
		return fmt.Errorf("绑定死信队列失败: %v", err)
	}

	// 声明主队列，设置死信交换器
	args := map[string]interface{}{
		"x-dead-letter-exchange":    dlxExchange,
		"x-dead-letter-routing-key": queueName,
		"x-message-ttl":             60000, // 消息过期时间 60秒
		"x-max-length":              10000, // 队列最大长度
	}

	// 声明主队列
	_, err = r.channel.QueueDeclare(
		queueName, // 队列名称
		true,      // 持久化
		false,     // 自动删除
		false,     // 排他性
		false,     // 非阻塞
		args,      // 额外参数
	)
	if err != nil {
		return fmt.Errorf("声明队列失败: %v", err)
	}

	// 消费消息
	msgs, err := channel.Consume(
		queueName, // 队列名称
		"",        // 消费者标签
		false,     // 自动确认
		false,     // 排他性
		false,     // 非阻塞
		false,     // 额外参数
		nil,       // amqp.Table 类型的参数
	)
	if err != nil {
		return fmt.Errorf("消费消息失败: %v", err)
	}

	// 启动goroutine处理消息
	go func() {
		for d := range msgs {
			utils.Logger.Info("收到消息", zap.String("message", string(d.Body)))

			// 获取重试次数
			retryCount := 0
			if retryVal, ok := d.Headers["x-retry-count"]; ok {
				if rc, ok := retryVal.(int64); ok {
					retryCount = int(rc)
				}
			}

			// 调用处理函数
			err := handler(d.Body)
			if err != nil {
				utils.Logger.Error("处理消息失败", zap.Error(err), zap.Int("retryCount", retryCount))

				// 检查重试次数
				if retryCount < 3 {
					// 增加重试次数
					retryCount++
					d.Headers["x-retry-count"] = int64(retryCount)
					// 重新发布消息，设置延迟
					err = r.channel.Publish(
						dlxExchange, // 死信交换器
						queueName,   // 路由键
						false,       // 强制
						false,       // 立即
						amqp.Publishing{
							ContentType:  "text/plain",
							Body:         d.Body,
							Headers:      d.Headers,
							DeliveryMode: amqp.Persistent,
						})
					if err != nil {
						utils.Logger.Error("重新发布消息失败", zap.Error(err))
					}
				} else {
					// 重试次数超过限制，直接拒绝，消息会进入死信队列
					utils.Logger.Warn("消息处理失败，已超过最大重试次数", zap.String("message", string(d.Body)))
				}
				// 拒绝消息
				d.Nack(false, false)
			} else {
				// 确认消息
				d.Ack(false)
			}
		}
	}()

	// 启动死信队列消费者
	r.consumeDeadLetterQueue(dlxQueue)

	// 启动定时对账任务
	go r.startReconciliationTask()

	utils.Logger.Info("消费者已启动", zap.String("queue", queueName))
	return nil
}

// consumeDeadLetterQueue 消费死信队列消息
func (r *RabbitMQ) consumeDeadLetterQueue(queueName string) error {
	channel, err := r.GetChannel()
	if err != nil {
		return err
	}

	// 消费死信队列消息
	msgs, err := channel.Consume(
		queueName, // 队列名称
		"",        // 消费者标签
		false,     // 自动确认
		false,     // 排他性
		false,     // 非阻塞
		false,     // 额外参数
		nil,       // amqp.Table 类型的参数
	)
	if err != nil {
		return fmt.Errorf("消费死信队列消息失败: %v", err)
	}

	// 启动goroutine处理死信消息
	go func() {
		for d := range msgs {
			utils.Logger.Warn("处理死信消息", zap.String("message", string(d.Body)))
			// 这里可以添加死信消息处理逻辑，例如：
			// 1. 记录到错误日志
			// 2. 发送告警通知
			// 3. 手动处理

			// 确认消息
			d.Ack(false)
		}
	}()

	utils.Logger.Info("死信队列消费者已启动", zap.String("queue", queueName))
	return nil
}

// startReconciliationTask 启动定时对账任务
func (r *RabbitMQ) startReconciliationTask() {
	ticker := time.NewTicker(5 * time.Minute) // 每5分钟执行一次对账
	defer ticker.Stop()

	for range ticker.C {
		utils.Logger.Info("开始执行对账任务")
		// 这里可以添加对账逻辑，例如：
		// 1. 检查Redis中的库存与数据库中的库存是否一致
		// 2. 检查消息队列中的待处理订单
		// 3. 处理异常订单
		utils.Logger.Info("对账任务执行完成")
	}
}

// HandleOrderMessage 处理订单消息
func HandleOrderMessage(message []byte) error {
	// 解析消息
	msgStr := string(message)
	parts := strings.Split(msgStr, ":")
	if len(parts) != 2 {
		return fmt.Errorf("无效的消息格式: %s", msgStr)
	}

	userID := parts[0]
	productID := parts[1]

	utils.Logger.Info("处理订单消息", zap.String("userID", userID), zap.String("productID", productID))

	// 模拟订单创建
	utils.Logger.Info("订单创建成功", zap.String("userID", userID), zap.String("productID", productID))

	return nil
}
