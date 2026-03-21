package service

import (
	"context"
	"errors"
	"seckill_go/utils"
	"time"

	"go.uber.org/zap"
)

// ErrEventNotFound 自定义错误
var (
	ErrEventNotFound = errors.New("活动不存在")
)

// 模拟活动数据
var mockEvents = []map[string]interface{}{
	{
		"id":         1,
		"product_id": 1,
		"start_time": time.Now().Add(-1 * time.Hour),
		"end_time":   time.Now().Add(23 * time.Hour),
	},
	{
		"id":         2,
		"product_id": 2,
		"start_time": time.Now().Add(-1 * time.Hour),
		"end_time":   time.Now().Add(23 * time.Hour),
	},
}

// CreateEvent 创建秒杀活动
func CreateEvent(ctx context.Context, event map[string]interface{}) error {
	// 生成新的活动ID
	newID := uint(len(mockEvents) + 1)
	event["id"] = newID

	// 添加到模拟数据
	mockEvents = append(mockEvents, event)

	utils.Logger.Info("活动创建成功", zap.Uint("id", newID), zap.Uint("product_id", event["product_id"].(uint)))
	return nil
}

// UpdateEvent 更新秒杀活动
func UpdateEvent(ctx context.Context, id uint, event map[string]interface{}) error {
	// 查找活动
	for i, e := range mockEvents {
		if e["id"] == id {
			// 更新活动信息
			if productID, ok := event["product_id"]; ok {
				mockEvents[i]["product_id"] = productID
			}
			if startTime, ok := event["start_time"]; ok {
				mockEvents[i]["start_time"] = startTime
			}
			if endTime, ok := event["end_time"]; ok {
				mockEvents[i]["end_time"] = endTime
			}

			utils.Logger.Info("活动更新成功", zap.Uint("id", id))
			return nil
		}
	}

	utils.Logger.Warn("活动不存在", zap.Uint("id", id))
	return ErrEventNotFound
}

// DeleteEvent 删除秒杀活动
func DeleteEvent(ctx context.Context, id uint) error {
	// 查找并删除活动
	for i, e := range mockEvents {
		if e["id"] == id {
			// 从切片中删除元素
			mockEvents = append(mockEvents[:i], mockEvents[i+1:]...)
			utils.Logger.Info("活动删除成功", zap.Uint("id", id))
			return nil
		}
	}

	utils.Logger.Warn("活动不存在", zap.Uint("id", id))
	return ErrEventNotFound
}

// GetEventByID 根据ID获取活动
func GetEventByID(ctx context.Context, id uint) (map[string]interface{}, error) {
	// 查找活动
	for _, event := range mockEvents {
		if event["id"] == id {
			return event, nil
		}
	}

	utils.Logger.Warn("活动不存在", zap.Uint("id", id))
	return nil, ErrEventNotFound
}

// GetEventList 获取活动列表
func GetEventList(ctx context.Context) ([]map[string]interface{}, error) {
	utils.Logger.Info("获取活动列表", zap.Int("count", len(mockEvents)))
	return mockEvents, nil
}

// StartEvent 启动秒杀活动
func StartEvent(ctx context.Context, id uint) error {
	// 查找活动
	for i, e := range mockEvents {
		if e["id"] == id {
			// 设置活动开始时间为当前时间
			mockEvents[i]["start_time"] = time.Now()
			utils.Logger.Info("活动启动成功", zap.Uint("id", id))
			return nil
		}
	}

	utils.Logger.Warn("活动不存在", zap.Uint("id", id))
	return ErrEventNotFound
}

// StopEvent 停止秒杀活动
func StopEvent(ctx context.Context, id uint) error {
	// 查找活动
	for i, e := range mockEvents {
		if e["id"] == id {
			// 设置活动结束时间为当前时间
			mockEvents[i]["end_time"] = time.Now()
			utils.Logger.Info("活动停止成功", zap.Uint("id", id))
			return nil
		}
	}

	utils.Logger.Warn("活动不存在", zap.Uint("id", id))
	return ErrEventNotFound
}
