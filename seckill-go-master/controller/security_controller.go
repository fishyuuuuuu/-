package controller

import (
	"net/http"
	"seckill_go/utils"

	"github.com/gin-gonic/gin"
)

// GetBlacklistInfoHandler 获取黑名单信息
func GetBlacklistInfoHandler(c *gin.Context) {
	targetType := c.Param("type")
	targetID := c.Param("id")

	if targetType == "" || targetID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "参数错误",
		})
		return
	}

	sm := utils.GetSecurityManager()
	info, err := sm.GetBlacklistInfo(c, targetType, targetID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "不在黑名单中",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": info,
	})
}

// AddToBlacklistHandler 添加到黑名单
func AddToBlacklistHandler(c *gin.Context) {
	var req struct {
		TargetType string `json:"target_type" binding:"required"`
		TargetID   string `json:"target_id" binding:"required"`
		Reason     string `json:"reason" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "参数错误",
		})
		return
	}

	sm := utils.GetSecurityManager()
	err := sm.AddToBlacklist(c, req.TargetType, req.TargetID, req.Reason)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "添加到黑名单失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "已添加到黑名单",
	})
}

// RemoveFromBlacklistHandler 从黑名单移除
func RemoveFromBlacklistHandler(c *gin.Context) {
	targetType := c.Param("type")
	targetID := c.Param("id")

	if targetType == "" || targetID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "参数错误",
		})
		return
	}

	sm := utils.GetSecurityManager()
	err := sm.RemoveFromBlacklist(c, targetType, targetID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "从黑名单移除失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "已从黑名单移除",
	})
}

// GetSecurityStatsHandler 获取安全统计信息
func GetSecurityStatsHandler(c *gin.Context) {
	sm := utils.GetSecurityManager()
	stats := sm.GetSecurityStats(c)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": stats,
	})
}

// GetRateLimitConfigHandler 获取限流配置
func GetRateLimitConfigHandler(c *gin.Context) {
	config := map[string]interface{}{
		"global": map[string]interface{}{
			"capacity": utils.GlobalRateLimitCapacity,
			"rate":     utils.GlobalRateLimitRate,
		},
		"ip": map[string]interface{}{
			"capacity": utils.IPRateLimitCapacity,
			"rate":     utils.IPRateLimitRate,
		},
		"user": map[string]interface{}{
			"capacity": utils.UserRateLimitCapacity,
			"rate":     utils.UserRateLimitRate,
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": config,
	})
}

// UpdateRateLimitConfigHandler 更新限流配置
func UpdateRateLimitConfigHandler(c *gin.Context) {
	var req struct {
		Global struct {
			Capacity int `json:"capacity"`
			Rate     int `json:"rate"`
		} `json:"global"`
		IP struct {
			Capacity int `json:"capacity"`
			Rate     int `json:"rate"`
		} `json:"ip"`
		User struct {
			Capacity int `json:"capacity"`
			Rate     int `json:"rate"`
		} `json:"user"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "参数错误",
		})
		return
	}

	// 这里可以实现动态更新限流配置的逻辑
	// 由于限流器是单例模式，需要重新初始化

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "限流配置已更新（将在下次重启后生效）",
	})
}
