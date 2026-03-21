package service

import (
	"context"
	"errors"
	"seckill_go/model"
	"seckill_go/utils"
	"time"

	"go.uber.org/zap"
)

// ErrOrderNotFound 自定义错误
var (
	ErrOrderNotFound      = errors.New("订单不存在")
	ErrOrderNotCancelable = errors.New("订单状态不允许取消")
)

// 模拟订单数据
var mockOrders = []map[string]interface{}{
	{
		"id":           1,
		"user_id":      1,
		"product_id":   1,
		"order_amount": 7999.99,
		"status":       0,
		"created_at":   time.Now().Add(-1 * time.Hour),
	},
	{
		"id":           2,
		"user_id":      1,
		"product_id":   2,
		"order_amount": 12999.99,
		"status":       1,
		"created_at":   time.Now().Add(-2 * time.Hour),
	},
}

// CreateOrder 创建订单
func CreateOrder(ctx context.Context, userID string, productID string) error {
	// 生成新的订单ID
	newID := uint(len(mockOrders) + 1)

	// 查找商品信息
	var product map[string]interface{}
	for _, p := range mockProducts {
		if p["id"].(int) == utils.StringToInt(productID) {
			product = p
			break
		}
	}

	if product == nil {
		utils.Logger.Error("商品不存在", zap.String("productID", productID))
		return ErrProductNotFound
	}

	// 创建订单
	order := map[string]interface{}{
		"id":           newID,
		"user_id":      utils.StringToUint(userID),
		"product_id":   utils.StringToUint(productID),
		"order_amount": product["price"],
		"status":       model.OrderStatusPending,
		"created_at":   time.Now(),
	}

	// 添加到模拟数据
	mockOrders = append(mockOrders, order)

	utils.Logger.Info("订单创建成功", zap.Uint("id", newID), zap.String("userID", userID), zap.String("productID", productID))
	return nil
}

// GetOrderList 获取用户订单列表
func GetOrderList(ctx context.Context, userID uint) ([]map[string]interface{}, error) {
	var userOrders []map[string]interface{}

	// 查找用户的订单
	for _, order := range mockOrders {
		if order["user_id"].(uint) == userID {
			userOrders = append(userOrders, order)
		}
	}

	utils.Logger.Info("获取订单列表", zap.Uint("userID", userID), zap.Int("count", len(userOrders)))
	return userOrders, nil
}

// GetOrderByID 根据ID获取订单
func GetOrderByID(ctx context.Context, id uint, userID uint) (map[string]interface{}, error) {
	// 查找订单
	for _, order := range mockOrders {
		if order["id"].(uint) == id {
			// 检查订单是否属于该用户
			if order["user_id"].(uint) != userID {
				utils.Logger.Warn("订单不属于当前用户", zap.Uint("orderID", id), zap.Uint("userID", userID))
				return nil, ErrOrderNotFound
			}
			return order, nil
		}
	}

	utils.Logger.Warn("订单不存在", zap.Uint("id", id))
	return nil, ErrOrderNotFound
}

// CancelOrder 取消订单
func CancelOrder(ctx context.Context, id uint, userID uint) error {
	// 查找订单
	for i, order := range mockOrders {
		if order["id"].(uint) == id {
			// 检查订单是否属于该用户
			if order["user_id"].(uint) != userID {
				utils.Logger.Warn("订单不属于当前用户", zap.Uint("orderID", id), zap.Uint("userID", userID))
				return ErrOrderNotFound
			}

			// 检查订单状态
			if order["status"].(model.OrderStatus) != model.OrderStatusPending {
				utils.Logger.Warn("订单状态不允许取消", zap.Uint("orderID", id), zap.Int("status", int(order["status"].(model.OrderStatus))))
				return ErrOrderNotCancelable
			}

			// 更新订单状态
			mockOrders[i]["status"] = model.OrderStatusCanceled
			utils.Logger.Info("订单取消成功", zap.Uint("id", id))
			return nil
		}
	}

	utils.Logger.Warn("订单不存在", zap.Uint("id", id))
	return ErrOrderNotFound
}
