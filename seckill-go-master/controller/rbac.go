package controller

import (
	"net/http"
	"seckill_go/service"

	"github.com/gin-gonic/gin"
)

// RBACMiddleware 权限检查中间件
func RBACMiddleware(permissionCode string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDValue, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
			c.Abort()
			return
		}

		userID, ok := userIDValue.(uint)
		if !ok || userID == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
			c.Abort()
			return
		}

		allowed, err := service.CheckPermission(userID, permissionCode)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "无权限访问"})
			c.Abort()
			return
		}
		if !allowed {
			c.JSON(http.StatusForbidden, gin.H{"error": "无权限访问"})
			c.Abort()
			return
		}

		c.Next()
	}
}
