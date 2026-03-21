package utils

import (
	"fmt"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Metrics 系统性能指标
type Metrics struct {
	RequestCount      int64           // 总请求数
	SuccessCount      int64           // 成功请求数
	ErrorCount        int64           // 错误请求数
	TotalResponseTime time.Duration   // 总响应时间
	AvgResponseTime   time.Duration   // 平均响应时间
	P99ResponseTime   time.Duration   // P99响应时间
	ResponseTimes     []time.Duration // 响应时间列表
	mu                sync.Mutex      // 互斥锁
}

// NewMetrics 创建一个新的性能指标实例
func NewMetrics() *Metrics {
	return &Metrics{
		ResponseTimes: make([]time.Duration, 0, 1000),
	}
}

// RecordRequest 记录请求信息
func (m *Metrics) RecordRequest(duration time.Duration, success bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.RequestCount++
	m.TotalResponseTime += duration
	m.ResponseTimes = append(m.ResponseTimes, duration)

	// 限制响应时间列表长度
	if len(m.ResponseTimes) > 1000 {
		m.ResponseTimes = m.ResponseTimes[len(m.ResponseTimes)-1000:]
	}

	if success {
		m.SuccessCount++
	} else {
		m.ErrorCount++
	}

	// 更新平均响应时间
	m.AvgResponseTime = m.TotalResponseTime / time.Duration(m.RequestCount)

	// 更新 P99 响应时间
	if len(m.ResponseTimes) > 0 {
		m.P99ResponseTime = m.calculateP99()
	}
}

// calculateP99 计算 P99 响应时间
func (m *Metrics) calculateP99() time.Duration {
	// 复制响应时间列表
	times := make([]time.Duration, len(m.ResponseTimes))
	copy(times, m.ResponseTimes)

	// 排序
	for i := 0; i < len(times); i++ {
		for j := i + 1; j < len(times); j++ {
			if times[i] > times[j] {
				times[i], times[j] = times[j], times[i]
			}
		}
	}

	// 计算 P99 索引
	index := int(float64(len(times)) * 0.99)
	if index >= len(times) {
		index = len(times) - 1
	}

	return times[index]
}

// GetMetrics 获取性能指标
func (m *Metrics) GetMetrics() map[string]interface{} {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 获取系统信息
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	return map[string]interface{}{
		"request_count":     m.RequestCount,
		"success_count":     m.SuccessCount,
		"error_count":       m.ErrorCount,
		"avg_response_time": m.AvgResponseTime.Milliseconds(),
		"p99_response_time": m.P99ResponseTime.Milliseconds(),
		"memory_alloc":      memStats.Alloc / 1024 / 1024,      // MB
		"memory_total":      memStats.TotalAlloc / 1024 / 1024, // MB
		"goroutines":        runtime.NumGoroutine(),
	}
}

// 全局性能指标实例
var globalMetrics = NewMetrics()

// MonitoringMiddleware 性能监控中间件
func MonitoringMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录开始时间
		start := time.Now()

		// 处理请求
		c.Next()

		// 计算响应时间
		duration := time.Since(start)

		// 记录请求信息
		success := c.Writer.Status() < http.StatusBadRequest
		globalMetrics.RecordRequest(duration, success)

		// 记录请求日志
		Logger.Info("请求处理完成",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("duration", duration),
			zap.String("ip", c.ClientIP()),
		)
	}
}

// GetMetrics 获取系统性能指标
func GetMetrics() map[string]interface{} {
	return globalMetrics.GetMetrics()
}

// StartMetricsServer 启动性能指标服务器
func StartMetricsServer(port int) {
	// 创建一个新的路由
	r := gin.Default()

	// 添加监控中间件
	r.Use(MonitoringMiddleware())

	// 健康检查接口
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// 性能指标接口
	r.GET("/metrics", func(c *gin.Context) {
		metrics := GetMetrics()
		c.JSON(http.StatusOK, metrics)
	})

	// 启动服务器
	serverAddr := fmt.Sprintf(":%d", port)
	Logger.Info("性能指标服务器启动", zap.String("addr", serverAddr))

	go func() {
		if err := r.Run(serverAddr); err != nil {
			Logger.Error("性能指标服务器启动失败", zap.Error(err))
		}
	}()
}
