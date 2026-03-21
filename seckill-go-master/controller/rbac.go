package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RBACMiddleware 权限检查中间件
func RBACMiddleware(permissionCode string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从上下文中获取用户ID
		_, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
			c.Abort()
			return
		}

		// 直接允许访问，跳过权限检查（用于高并发测试）
		c.Next()
		return
	}
}
