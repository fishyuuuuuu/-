package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RateLimitMiddleware IP速率限制中间件
func RateLimitMiddleware() gin.HandlerFunc {
	ipRateLimiter := GetIPRateLimiter()

	return func(c *gin.Context) {
		// 获取客户端IP
		clientIP := c.ClientIP()

		// 检查是否允许请求
		if !ipRateLimiter.Allow(clientIP) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "请求过于频繁，请稍后再试",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
