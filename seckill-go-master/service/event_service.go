package service

import (
	"context"
	"errors"
	"fmt"
	"seckill_go/utils"
	"time"

	"go.uber.org/zap"
)

// ErrEventNotFound 自定义错误
var (
	ErrEventNotFound   = errors.New("活动不存在")
	ErrEventNotBound   = errors.New("商品未绑定秒杀活动")
	ErrEventNotStarted = errors.New("秒杀未开始")
	ErrEventEnded      = errors.New("活动已结束")
	ErrEventInvalid    = errors.New("活动数据无效")
)

const (
	SeckillStatusNoEvent    = "no_event"
	SeckillStatusNotStarted = "not_started"
	SeckillStatusOngoing    = "ongoing"
	SeckillStatusEnded      = "ended"
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

func toUintFromAny(v interface{}) (uint, bool) {
	switch val := v.(type) {
	case uint:
		return val, true
	case uint64:
		return uint(val), true
	case int:
		if val < 0 {
			return 0, false
		}
		return uint(val), true
	case int64:
		if val < 0 {
			return 0, false
		}
		return uint(val), true
	case float64:
		if val < 0 {
			return 0, false
		}
		return uint(val), true
	default:
		return 0, false
	}
}

func parseEventData(event map[string]interface{}) (productID uint, startTime time.Time, endTime time.Time, err error) {
	pid, ok := toUintFromAny(event["product_id"])
	if !ok || pid == 0 {
		return 0, time.Time{}, time.Time{}, fmt.Errorf("%w: product_id", ErrEventInvalid)
	}
	start, ok := event["start_time"].(time.Time)
	if !ok {
		return 0, time.Time{}, time.Time{}, fmt.Errorf("%w: start_time", ErrEventInvalid)
	}
	end, ok := event["end_time"].(time.Time)
	if !ok {
		return 0, time.Time{}, time.Time{}, fmt.Errorf("%w: end_time", ErrEventInvalid)
	}
	if !start.Before(end) {
		return 0, time.Time{}, time.Time{}, fmt.Errorf("%w: start/end", ErrEventInvalid)
	}
	return pid, start, end, nil
}

func productExists(productID uint) bool {
	for _, p := range mockProducts {
		id, ok := toUintFromAny(p["id"])
		if ok && id == productID {
			return true
		}
	}
	return false
}

func findEventByProductID(productID uint) (int, map[string]interface{}) {
	for i, e := range mockEvents {
		pid, ok := toUintFromAny(e["product_id"])
		if ok && pid == productID {
			return i, e
		}
	}
	return -1, nil
}

func eventIDEquals(event map[string]interface{}, id uint) bool {
	eventID, ok := toUintFromAny(event["id"])
	return ok && eventID == id
}

func nextEventID() uint {
	var maxID uint
	for _, e := range mockEvents {
		if id, ok := toUintFromAny(e["id"]); ok && id > maxID {
			maxID = id
		}
	}
	return maxID + 1
}

func buildSeckillInfoByEvent(event map[string]interface{}, now time.Time) map[string]interface{} {
	productID, startTime, endTime, err := parseEventData(event)
	if err != nil {
		return map[string]interface{}{
			"has_event":    false,
			"start_time":   "",
			"end_time":     "",
			"can_seckill":  false,
			"status":       SeckillStatusNoEvent,
			"status_text":  "商品未绑定秒杀活动",
			"product_id":   productID,
			"current_time": now.Format(time.RFC3339),
		}
	}

	status := SeckillStatusOngoing
	statusText := "秒杀进行中"
	canSeckill := true
	switch {
	case now.Before(startTime):
		status = SeckillStatusNotStarted
		statusText = "秒杀未开始"
		canSeckill = false
	case now.After(endTime):
		status = SeckillStatusEnded
		statusText = "活动已结束"
		canSeckill = false
	}

	return map[string]interface{}{
		"has_event":    true,
		"start_time":   startTime.Format(time.RFC3339),
		"end_time":     endTime.Format(time.RFC3339),
		"can_seckill":  canSeckill,
		"status":       status,
		"status_text":  statusText,
		"product_id":   productID,
		"current_time": now.Format(time.RFC3339),
	}
}

func BuildSeckillInfoByProductID(ctx context.Context, productID uint) map[string]interface{} {
	_ = ctx
	now := time.Now()
	_, event := findEventByProductID(productID)
	if event == nil {
		return map[string]interface{}{
			"has_event":    false,
			"start_time":   "",
			"end_time":     "",
			"can_seckill":  false,
			"status":       SeckillStatusNoEvent,
			"status_text":  "商品未绑定秒杀活动",
			"product_id":   productID,
			"current_time": now.Format(time.RFC3339),
		}
	}
	return buildSeckillInfoByEvent(event, now)
}

func ValidateSeckillWindow(ctx context.Context, productID uint) error {
	_ = ctx
	info := BuildSeckillInfoByProductID(ctx, productID)
	canSeckill, _ := info["can_seckill"].(bool)
	if canSeckill {
		return nil
	}
	status, _ := info["status"].(string)
	switch status {
	case SeckillStatusNoEvent:
		return ErrEventNotBound
	case SeckillStatusNotStarted:
		return ErrEventNotStarted
	case SeckillStatusEnded:
		return ErrEventEnded
	default:
		return ErrEventInvalid
	}
}

// CreateEvent 创建秒杀活动
func CreateEvent(ctx context.Context, event map[string]interface{}) error {
	productID, startTime, endTime, err := parseEventData(event)
	if err != nil {
		return err
	}
	if !productExists(productID) {
		return ErrProductNotFound
	}
	if idx, _ := findEventByProductID(productID); idx >= 0 {
		mockEvents[idx]["start_time"] = startTime
		mockEvents[idx]["end_time"] = endTime
		utils.Logger.Info("商品已存在活动，执行覆盖更新", zap.Uint("product_id", productID))
		return nil
	}

	// 生成新的活动ID
	newID := nextEventID()
	event["id"] = newID
	event["product_id"] = productID
	event["start_time"] = startTime
	event["end_time"] = endTime

	// 添加到模拟数据
	mockEvents = append(mockEvents, event)

	utils.Logger.Info("活动创建成功", zap.Uint("id", newID), zap.Uint("product_id", productID))
	return nil
}

// UpdateEvent 更新秒杀活动
func UpdateEvent(ctx context.Context, id uint, event map[string]interface{}) error {
	// 查找活动
	for i, e := range mockEvents {
		if eventIDEquals(e, id) {
			// 更新活动信息
			if productID, ok := event["product_id"]; ok {
				if pid, castOK := toUintFromAny(productID); castOK && pid > 0 {
					if !productExists(pid) {
						return ErrProductNotFound
					}
					mockEvents[i]["product_id"] = pid
				} else {
					return ErrEventInvalid
				}
			}
			if startTime, ok := event["start_time"]; ok {
				start, castOK := startTime.(time.Time)
				if !castOK {
					return ErrEventInvalid
				}
				mockEvents[i]["start_time"] = start
			}
			if endTime, ok := event["end_time"]; ok {
				end, castOK := endTime.(time.Time)
				if !castOK {
					return ErrEventInvalid
				}
				mockEvents[i]["end_time"] = end
			}

			_, start, end, parseErr := parseEventData(mockEvents[i])
			if parseErr != nil || !start.Before(end) {
				return ErrEventInvalid
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
		if eventIDEquals(e, id) {
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
		if eventIDEquals(event, id) {
			return event, nil
		}
	}

	utils.Logger.Warn("活动不存在", zap.Uint("id", id))
	return nil, ErrEventNotFound
}

func GetEventByProductID(ctx context.Context, productID uint) (map[string]interface{}, error) {
	_ = ctx
	_, event := findEventByProductID(productID)
	if event == nil {
		return nil, ErrEventNotFound
	}
	return event, nil
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
		if eventIDEquals(e, id) {
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
		if eventIDEquals(e, id) {
			// 设置活动结束时间为当前时间
			mockEvents[i]["end_time"] = time.Now()
			utils.Logger.Info("活动停止成功", zap.Uint("id", id))
			return nil
		}
	}

	utils.Logger.Warn("活动不存在", zap.Uint("id", id))
	return ErrEventNotFound
}
