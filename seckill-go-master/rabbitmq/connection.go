package rabbitmq

import (
	"fmt"
	"github.com/streadway/amqp"
	"seckill_go/config" // 替换为你的实际项目路径
	"sync"
	"time"
)

var (
	instance *RabbitMQ
	mu       sync.Mutex
)

// RabbitMQ 连接管理器
type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	config  config.RabbitMQConfig
	mu      sync.Mutex // 用于并发安全
}

// Init 初始化RabbitMQ连接，只能在main中调用一次
func Init() error {
	mu.Lock()
	defer mu.Unlock()
	// 从全局配置中获取RabbitMQ配置
	cfg := config.GetConfig().RabbitMQ
	instance = &RabbitMQ{
		config: cfg,
	}
	// 初始化连接
	if err := instance.connect(); err != nil {
		return err
	}
	return nil
}

// GetInstance 获取已初始化的RabbitMQ实例
func GetInstance() (*RabbitMQ, error) {
	mu.Lock()
	defer mu.Unlock()
	// 检查连接是否有效
	if instance.conn == nil || instance.conn.IsClosed() {
		if err := instance.Reconnect(); err != nil {
			return nil, err
		}
	}
	return instance, nil
}

// 建立与RabbitMQ的连接
func (r *RabbitMQ) connect() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// 构建连接字符串
	addr := fmt.Sprintf("amqp://%s:%s@%s:%d/",
		r.config.Username,
		r.config.Password,
		r.config.Host,
		r.config.Port)

	// 连接到RabbitMQ服务器
	conn, err := amqp.Dial(addr)
	if err != nil {
		return fmt.Errorf("无法连接到RabbitMQ: %v", err)
	}

	r.conn = conn

	// 创建一个通道
	channel, err := conn.Channel()
	if err != nil {
		r.conn.Close()
		return fmt.Errorf("无法创建通道: %v", err)
	}

	r.channel = channel
	return nil
}

// GetChannel 获取RabbitMQ通道
func (r *RabbitMQ) GetChannel() (*amqp.Channel, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// 检查连接是否有效
	if r.conn == nil || r.conn.IsClosed() {
		if err := r.Reconnect(); err != nil {
			return nil, err
		}
	}

	// 检查通道是否有效
	if r.channel == nil {
		var err error
		r.channel, err = r.conn.Channel()
		if err != nil {
			return nil, fmt.Errorf("无法创建通道: %v", err)
		}
	}

	return r.channel, nil
}

// Reconnect 重新连接RabbitMQ
func (r *RabbitMQ) Reconnect() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	// 关闭现有连接
	if r.conn != nil {
		r.conn.Close()
	}
	// 重试连接机制
	for i := 0; i < 5; i++ {
		if err := r.connect(); err != nil {
			fmt.Printf("重连失败，重试中... (%d/5)\n", i+1)
			time.Sleep(2 * time.Second)
			continue
		}
		fmt.Println("重新连接RabbitMQ成功")
		return nil
	}
	return fmt.Errorf("多次尝试后仍无法连接到RabbitMQ")
}

// Close 关闭连接和通道
func (r *RabbitMQ) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.channel != nil {
		if err := r.channel.Close(); err != nil {
			return err
		}
	}

	if r.conn != nil {
		return r.conn.Close()
	}

	return nil
}
