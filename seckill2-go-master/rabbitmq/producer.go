package rabbitmq

import (
	"fmt"
	"github.com/streadway/amqp"
)

// DeclareExchange 声明交换机
func (r *RabbitMQ) DeclareExchange(exchangeName, exchangeType string) error {
	channel, err := r.GetChannel()
	if err != nil {
		return err
	}

	// 声明交换机，持久化
	return channel.ExchangeDeclare(
		exchangeName, // 交换机名称
		exchangeType, // 交换机类型：direct, topic, fanout, headers
		true,         // 持久化
		false,        // 自动删除
		false,        // 内部的
		false,        // 非阻塞
		nil,          // 额外参数
	)
}

// DeclareQueue 声明队列
func (r *RabbitMQ) DeclareQueue(queueName string) error {
	channel, err := r.GetChannel()
	if err != nil {
		return err
	}

	// 声明队列，持久化
	_, err = channel.QueueDeclare(
		queueName, // 队列名称
		true,      // 持久化
		false,     // 自动删除
		false,     // 排他性
		false,     // 非阻塞
		nil,       // 额外参数
	)
	return err
}

// BindQueue 绑定队列到交换机
func (r *RabbitMQ) BindQueue(queueName, routingKey, exchangeName string) error {
	channel, err := r.GetChannel()
	if err != nil {
		return err
	}

	return channel.QueueBind(
		queueName,    // 队列名称
		routingKey,   // 路由键
		exchangeName, // 交换机名称
		false,        // 非阻塞
		nil,          // 额外参数
	)
}

// SendToQueue 发送消息到指定队列
func (r *RabbitMQ) SendToQueue(queueName string, message []byte) error {
	channel, err := r.GetChannel()
	if err != nil {
		return err
	}

	// 发送消息，持久化消息
	err = channel.Publish(
		"",        // 交换机名称，空表示使用默认交换机
		queueName, // 路由键，使用队列名称
		false,     // 强制的
		false,     // 立即的
		amqp.Publishing{
			ContentType:  "text/plain",
			Body:         message,
			DeliveryMode: amqp.Persistent, // 持久化消息
		},
	)

	if err != nil {
		return fmt.Errorf("发送消息失败: %v", err)
	}
	return nil
}

// SendToExchange 发送消息到交换机
func (r *RabbitMQ) SendToExchange(exchangeName, routingKey string, message []byte) error {
	channel, err := r.GetChannel()
	if err != nil {
		return err
	}

	// 发送消息到交换机
	err = channel.Publish(
		exchangeName, // 交换机名称
		routingKey,   // 路由键
		false,        // 强制的
		false,        // 立即的
		amqp.Publishing{
			ContentType:  "text/plain",
			Body:         message,
			DeliveryMode: amqp.Persistent, // 持久化消息
		},
	)

	if err != nil {
		return fmt.Errorf("发送消息到交换机失败: %v", err)
	}
	return nil
}
