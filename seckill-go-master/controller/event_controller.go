package controller

import (
	"errors"
	"net/http"
	"seckill_go/service"
	"seckill_go/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func parseEventTime(value interface{}) (time.Time, bool) {
	raw, ok := value.(string)
	if !ok {
		return time.Time{}, false
	}
	layouts := []string{
		time.RFC3339,
		"2006-01-02T15:04:05",
		"2006-01-02T15:04",
	}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, raw); err == nil {
			return t, true
		}
	}
	return time.Time{}, false
}

// CreateEventHandler 创建秒杀活动请求
func CreateEventHandler(c *gin.Context) {
	var event map[string]interface{}
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	// 验证必要字段
	productValue, ok := event["product_id"]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "商品ID不能为空"})
		return
	}
	switch v := productValue.(type) {
	case float64:
		event["product_id"] = uint(v)
	case int:
		event["product_id"] = uint(v)
	case uint:
		event["product_id"] = v
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "商品ID格式错误"})
		return
	}
	if _, ok := event["start_time"]; !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "开始时间不能为空"})
		return
	}
	if _, ok := event["end_time"]; !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "结束时间不能为空"})
		return
	}

	// 解析时间
	startTime, ok := parseEventTime(event["start_time"])
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "开始时间格式错误"})
		return
	}
	event["start_time"] = startTime

	endTime, ok := parseEventTime(event["end_time"])
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "结束时间格式错误"})
		return
	}
	event["end_time"] = endTime

	// 验证时间顺序
	if endTime.Before(startTime) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "结束时间不能早于开始时间"})
		return
	}

	if err := service.CreateEvent(c, event); err != nil {
		if errors.Is(err, service.ErrProductNotFound) || errors.Is(err, service.ErrEventInvalid) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建活动失败"})
		return
	}

	seckillInfo := service.BuildSeckillInfoByProductID(c, event["product_id"].(uint))
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "活动创建成功",
		"data": map[string]interface{}{
			"product_id": event["product_id"],
			"start_time": startTime.Format(time.RFC3339),
			"end_time":   endTime.Format(time.RFC3339),
			"seckill":    seckillInfo,
		},
	})
}

// UpdateEventHandler 更新秒杀活动请求
func UpdateEventHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的活动ID"})
		return
	}

	var event map[string]interface{}
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	// 解析时间
	if productValue, ok := event["product_id"]; ok {
		switch v := productValue.(type) {
		case float64:
			event["product_id"] = uint(v)
		case int:
			event["product_id"] = uint(v)
		case uint:
			event["product_id"] = v
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "商品ID格式错误"})
			return
		}
	}

	// 解析时间
	if startTimeStr, ok := event["start_time"].(string); ok {
		startTime, parsed := parseEventTime(startTimeStr)
		if !parsed {
			c.JSON(http.StatusBadRequest, gin.H{"error": "开始时间格式错误"})
			return
		}
		event["start_time"] = startTime
	}

	if endTimeStr, ok := event["end_time"].(string); ok {
		endTime, parsed := parseEventTime(endTimeStr)
		if !parsed {
			c.JSON(http.StatusBadRequest, gin.H{"error": "结束时间格式错误"})
			return
		}
		event["end_time"] = endTime
	}

	// 验证时间顺序
	if startTime, ok := event["start_time"].(time.Time); ok {
		if endTime, ok := event["end_time"].(time.Time); ok {
			if endTime.Before(startTime) {
				c.JSON(http.StatusBadRequest, gin.H{"error": "结束时间不能早于开始时间"})
				return
			}
		}
	}

	err = service.UpdateEvent(c, uint(id), event)
	if err != nil {
		if errors.Is(err, service.ErrEventNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "活动不存在"})
			return
		}
		if errors.Is(err, service.ErrProductNotFound) || errors.Is(err, service.ErrEventInvalid) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新活动失败"})
		return
	}

	updatedEvent, getErr := service.GetEventByID(c, uint(id))
	if getErr != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "活动更新成功",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "活动更新成功",
		"data":    updatedEvent,
	})
}

// DeleteEventHandler 删除秒杀活动请求
func DeleteEventHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的活动ID"})
		return
	}

	err = service.DeleteEvent(c, uint(id))
	if err != nil {
		if err == service.ErrEventNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "活动不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除活动失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "活动删除成功"})
}

// GetEventByIDHandler 根据ID获取活动请求
func GetEventByIDHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的活动ID"})
		return
	}

	event, err := service.GetEventByID(c, uint(id))
	if err != nil {
		if err == service.ErrEventNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "活动不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取活动失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    event,
	})
}

// GetEventListHandler 获取活动列表请求
func GetEventListHandler(c *gin.Context) {
	events, err := service.GetEventList(c)
	if err != nil {
		utils.Logger.Error("获取活动列表失败", zap.Error(err))
		c.JSON(500, gin.H{
			"code":    500,
			"message": "无法获取活动列表",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    events,
	})
}

// StartEventHandler 启动秒杀活动请求
func StartEventHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的活动ID"})
		return
	}

	err = service.StartEvent(c, uint(id))
	if err != nil {
		if err == service.ErrEventNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "活动不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "启动活动失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "活动启动成功"})
}

// StopEventHandler 停止秒杀活动请求
func StopEventHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的活动ID"})
		return
	}

	err = service.StopEvent(c, uint(id))
	if err != nil {
		if err == service.ErrEventNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "活动不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "停止活动失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "活动停止成功"})
}
