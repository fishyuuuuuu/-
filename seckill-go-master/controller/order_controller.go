package controller

import (
	"net/http"
	"seckill_go/service"
	"seckill_go/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GetOrderListHandler 获取订单列表请求
func GetOrderListHandler(c *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户未登录"})
		return
	}

	orders, err := service.GetOrderList(c, userID.(uint))
	if err != nil {
		utils.Logger.Error("获取订单列表失败", zap.Error(err))
		c.JSON(500, gin.H{
			"code":    500,
			"message": "无法获取订单列表",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    orders,
	})
}

// GetOrderByIDHandler 根据ID获取订单请求
func GetOrderByIDHandler(c *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户未登录"})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的订单ID"})
		return
	}

	order, err := service.GetOrderByID(c, uint(id), userID.(uint))
	if err != nil {
		if err == service.ErrOrderNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "订单不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取订单失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    order,
	})
}

// CancelOrderHandler 取消订单请求
func CancelOrderHandler(c *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户未登录"})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的订单ID"})
		return
	}

	err = service.CancelOrder(c, uint(id), userID.(uint))
	if err != nil {
		if err == service.ErrOrderNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "订单不存在"})
			return
		}
		if err == service.ErrOrderNotCancelable {
			c.JSON(http.StatusBadRequest, gin.H{"error": "订单状态不允许取消"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "取消订单失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "订单取消成功"})
}
