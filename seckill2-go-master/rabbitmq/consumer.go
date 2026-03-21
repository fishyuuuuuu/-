package rabbitmq

import (
	"fmt"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

type Consumer struct {
	rmq    *RabbitMQ
	logger *zap.Logger
}

// NewConsumer 创建新的 RabbitMQ 消费者
func NewConsumer(logger *zap.Logger) (*Consumer, error) {
	// 获取 RabbitMQ 实例
	rmq, err := GetInstance()
	if err != nil {
		logger.Error("无法获取 RabbitMQ 实例", zap.Error(err))
		return nil, fmt.Errorf("无法获取 RabbitMQ 实例: %v", err)
	}

	return &Consumer{
		rmq:    rmq,
		logger: logger,
	}, nil
}

// ConsumeFromQueue 从指定队列消费消息
func (c *Consumer) ConsumeFromQueue(queueName, consumerTag string) (<-chan amqp.Delivery, error) {
	channel, err := c.rmq.GetChannel()
	if err != nil {
		c.logger.Error("无法获取 RabbitMQ 通道", zap.Error(err))
		return nil, fmt.Errorf("无法获取通道: %v", err)
	}

	// 声明队列
	if err := c.rmq.DeclareQueue(queueName); err != nil {
		c.logger.Error("无法声明队列", zap.Error(err))
		return nil, fmt.Errorf("无法声明队列 %s: %v", queueName, err)
	}

	// 设置 QoS
	if err := channel.Qos(1, 0, false); err != nil {
		c.logger.Error("无法设置 QoS", zap.Error(err))
		return nil, fmt.Errorf("无法设置 QoS: %v", err)
	}

	// 注册消费者
	msgs, err := channel.Consume(
		queueName,   // 队列名称
		consumerTag, // 消费者标签
		false,       // 手动确认
		false,       // 非独占
		false,       // 不等待
		false,       // 非局部
		nil,         // 额外参数
	)
	if err != nil {
		c.logger.Error("无法注册消费者", zap.Error(err))
		return nil, fmt.Errorf("无法注册消费者: %v", err)
	}

	c.logger.Info("成功注册消费者", zap.String("queue", queueName), zap.String("consumer_tag", consumerTag))
	return msgs, nil
}

// ConsumeFromExchange 从交换机消费消息
func (c *Consumer) ConsumeFromExchange(queueName, routingKey, exchangeName, consumerTag string) (<-chan amqp.Delivery, error) {
	channel, err := c.rmq.GetChannel()
	if err != nil {
		c.logger.Error("无法获取 RabbitMQ 通道", zap.Error(err))
		return nil, fmt.Errorf("无法获取通道: %v", err)
	}

	// 声明队列
	if err := c.rmq.DeclareQueue(queueName); err != nil {
		c.logger.Error("无法声明队列", zap.Error(err))
		return nil, fmt.Errorf("无法声明队列 %s: %v", queueName, err)
	}

	// 绑定队列到交换机
	if err := c.rmq.BindQueue(queueName, routingKey, exchangeName); err != nil {
		c.logger.Error("无法绑定队列到交换机", zap.Error(err))
		return nil, fmt.Errorf("无法绑定队列 %s 到交换机 %s: %v", queueName, exchangeName, err)
	}

	// 设置 QoS
	if err := channel.Qos(1, 0, false); err != nil {
		c.logger.Error("无法设置 QoS", zap.Error(err))
		return nil, fmt.Errorf("无法设置 QoS: %v", err)
	}

	// 注册消费者
	msgs, err := channel.Consume(
		queueName,   // 队列名称
		consumerTag, // 消费者标签
		false,       // 手动确认
		false,       // 非独占
		false,       // 不等待
		false,       // 非局部
		nil,         // 额外参数
	)
	if err != nil {
		c.logger.Error("无法注册消费者", zap.Error(err))
		return nil, fmt.Errorf("无法注册消费者: %v", err)
	}

	c.logger.Info("成功注册消费者",
		zap.String("queue", queueName),
		zap.String("exchange", exchangeName),
		zap.String("routing_key", routingKey),
		zap.String("consumer_tag", consumerTag))
	return msgs, nil
}

// StartConsuming 开始消费消息并处理
func (c *Consumer) StartConsuming(queueName, consumerTag string, useExchange bool, exchangeName, routingKey string) error {
	var msgs <-chan amqp.Delivery
	var err error

	if useExchange {
		// 从交换机消费
		msgs, err = c.ConsumeFromExchange(queueName, routingKey, exchangeName, consumerTag)
	} else {
		// 直接从队列消费
		msgs, err = c.ConsumeFromQueue(queueName, consumerTag)
	}
	if err != nil {
		return err
	}

	// 在 Goroutine 中处理消息
	go func() {
		// 处理消息
		for d := range msgs {
			c.logger.Info("收到消息", zap.String("body", string(d.Body)))
			// 模拟消息处理
			// 手动确认消息
			if err := d.Ack(false); err != nil {
				c.logger.Error("消息确认失败", zap.Error(err))
			} else {
				c.logger.Info("消息处理完成", zap.String("body", string(d.Body)))
			}
		}

		c.logger.Warn("消息通道已关闭", zap.String("queue", queueName))
		// 通道关闭后尝试重连
		if err := c.rmq.Reconnect(); err != nil {
			c.logger.Error("重连失败", zap.Error(err))
		} else {
			c.logger.Info("重连成功，重新启动消费者")
			// 递归调用以重新启动消费
			if err := c.StartConsuming(queueName, consumerTag, useExchange, exchangeName, routingKey); err != nil {
				c.logger.Error("重新启动消费者失败", zap.Error(err))
			}
		}
	}()

	c.logger.Info("消费者已启动，等待消息...", zap.String("queue", queueName))
	return nil
}

// Reconnect 重新连接 RabbitMQ
func (c *Consumer) Reconnect() error {
	c.logger.Info("尝试重新连接 RabbitMQ")
	if err := c.rmq.Reconnect(); err != nil {
		c.logger.Error("RabbitMQ 重连失败", zap.Error(err))
		return fmt.Errorf("RabbitMQ 重连失败: %v", err)
	}
	c.logger.Info("RabbitMQ 重连成功")
	return nil
}

// Close 关闭消费者
func (c *Consumer) Close() error {
	c.logger.Info("正在关闭 RabbitMQ 消费者")
	if err := c.rmq.Close(); err != nil {
		c.logger.Error("关闭 RabbitMQ 连接失败", zap.Error(err))
		return fmt.Errorf("关闭 RabbitMQ 连接失败: %v", err)
	}
	c.logger.Info("RabbitMQ 消费者已关闭")
	return nil
}
