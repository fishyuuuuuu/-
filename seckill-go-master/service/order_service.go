package service

import (
	"context"
	"errors"
	"seckill_go/db"
	"seckill_go/model"
	"seckill_go/utils"
	"sync"
	"time"

	"go.uber.org/zap"
)

// ErrOrderNotFound 自定义错误
var (
	ErrOrderNotFound      = errors.New("订单不存在")
	ErrOrderNotCancelable = errors.New("订单状态不允许取消")
)

var (
	mockOrders   []map[string]interface{}
	mockOrdersMu sync.RWMutex
	nextOrderID  uint = 5
)

func init() {
	mockOrders = []map[string]interface{}{
		{
			"id":           1,
			"user_id":      2,
			"user_name":    "user1",
			"product_id":   1,
			"product_name": "iPhone 15 Pro",
			"order_amount": 7999.99,
			"status":       1,
			"created_at":   time.Now().Add(-24 * time.Hour),
		},
		{
			"id":           2,
			"user_id":      2,
			"user_name":    "user1",
			"product_id":   2,
			"product_name": "MacBook Air M3",
			"order_amount": 11999.00,
			"status":       0,
			"created_at":   time.Now().Add(-12 * time.Hour),
		},
		{
			"id":           3,
			"user_id":      3,
			"user_name":    "user2",
			"product_id":   3,
			"product_name": "AirPods Pro 2",
			"order_amount": 1899.00,
			"status":       3,
			"created_at":   time.Now().Add(-6 * time.Hour),
		},
		{
			"id":           4,
			"user_id":      3,
			"user_name":    "user2",
			"product_id":   4,
			"product_name": "iPad Pro 12.9",
			"order_amount": 9299.00,
			"status":       2,
			"created_at":   time.Now().Add(-1 * time.Hour),
		},
	}
}

// CreateOrder 创建订单（已废弃，使用消息队列异步创建）
func CreateOrder(ctx context.Context, userID string, productID string) error {
	// 此函数已废弃，订单通过消息队列异步创建
	utils.Logger.Warn("CreateOrder函数已废弃，请使用消息队列")
	return nil
}

// CreateMockOrder 创建模拟订单
func CreateMockOrder(userID uint, productID uint) (map[string]interface{}, error) {
	mockOrdersMu.Lock()
	defer mockOrdersMu.Unlock()

	var productName string
	var price float64

	for _, p := range mockProducts {
		if p["id"].(int) == int(productID) {
			productName = p["name"].(string)
			price = p["price"].(float64)
			break
		}
	}

	if productName == "" {
		productName = "未知商品"
		price = 0.0
	}

	var userName string
	if userID == 2 {
		userName = "user1"
	} else if userID == 3 {
		userName = "user2"
	} else {
		userName = "user" + string(rune(userID))
	}

	newOrder := map[string]interface{}{
		"id":           nextOrderID,
		"user_id":      userID,
		"user_name":    userName,
		"product_id":   productID,
		"product_name": productName,
		"order_amount": price,
		"status":       0,
		"created_at":   time.Now(),
	}

	nextOrderID++
	mockOrders = append([]map[string]interface{}{newOrder}, mockOrders...)

	return newOrder, nil
}

// GetOrderList 获取用户订单列表
func GetOrderList(ctx context.Context, userID uint) ([]map[string]interface{}, error) {
	var result []map[string]interface{}

	if db.DB != nil {
		var orders []model.Order
		if err := db.DB.Where("user_id = ?", userID).Order("created_at DESC").Find(&orders).Error; err == nil && len(orders) > 0 {
			for _, order := range orders {
				var product model.Product
				db.DB.First(&product, order.ProductID)

				result = append(result, map[string]interface{}{
					"id":           order.ID,
					"user_id":      order.UserID,
					"product_id":   order.ProductID,
					"product_name": product.Name,
					"order_amount": order.OrderAmount,
					"status":       order.Status,
					"created_at":   order.CreatedAt,
				})
			}
			return result, nil
		}
	}

	mockOrdersMu.RLock()
	defer mockOrdersMu.RUnlock()

	for _, order := range mockOrders {
		if order["user_id"].(uint) == userID {
			result = append(result, map[string]interface{}{
				"id":           order["id"],
				"user_id":      order["user_id"],
				"product_id":   order["product_id"],
				"product_name": order["product_name"],
				"order_amount": order["order_amount"],
				"status":       order["status"],
				"created_at":   order["created_at"],
			})
		}
	}

	utils.Logger.Info("获取订单列表", zap.Uint("userID", userID), zap.Int("count", len(result)))
	return result, nil
}

// GetAllOrders 获取所有订单（管理员使用）
func GetAllOrders(ctx context.Context) ([]map[string]interface{}, error) {
	var result []map[string]interface{}

	if db.DB != nil {
		var orders []model.Order
		if err := db.DB.Order("created_at DESC").Find(&orders).Error; err == nil && len(orders) > 0 {
			for _, order := range orders {
				var product model.Product
				db.DB.First(&product, order.ProductID)

				var user model.User
				db.DB.First(&user, order.UserID)

				result = append(result, map[string]interface{}{
					"id":           order.ID,
					"user_id":      order.UserID,
					"user_name":    user.Username,
					"product_id":   order.ProductID,
					"product_name": product.Name,
					"order_amount": order.OrderAmount,
					"status":       order.Status,
					"created_at":   order.CreatedAt,
				})
			}
			utils.Logger.Info("获取所有订单", zap.Int("count", len(result)))
			return result, nil
		}
	}

	mockOrdersMu.RLock()
	defer mockOrdersMu.RUnlock()

	for _, order := range mockOrders {
		result = append(result, map[string]interface{}{
			"id":           order["id"],
			"user_id":      order["user_id"],
			"user_name":    order["user_name"],
			"product_id":   order["product_id"],
			"product_name": order["product_name"],
			"order_amount": order["order_amount"],
			"status":       order["status"],
			"created_at":   order["created_at"],
		})
	}

	utils.Logger.Info("获取所有订单(模拟数据)", zap.Int("count", len(result)))
	return result, nil
}

// GetOrderByID 根据ID获取订单
func GetOrderByID(ctx context.Context, id uint, userID uint) (map[string]interface{}, error) {
	var order model.Order
	if db.DB != nil {
		if err := db.DB.First(&order, id).Error; err == nil {
			if order.UserID == userID {
				var product model.Product
				db.DB.First(&product, order.ProductID)

				return map[string]interface{}{
					"id":           order.ID,
					"user_id":      order.UserID,
					"product_id":   order.ProductID,
					"product_name": product.Name,
					"order_amount": order.OrderAmount,
					"status":       order.Status,
					"created_at":   order.CreatedAt,
				}, nil
			}
			utils.Logger.Warn("订单不属于当前用户", zap.Uint("orderID", id), zap.Uint("userID", userID))
			return nil, ErrOrderNotFound
		}
	}

	mockOrdersMu.RLock()
	defer mockOrdersMu.RUnlock()

	for _, order := range mockOrders {
		if order["id"].(uint) == id {
			if order["user_id"].(uint) == userID {
				return map[string]interface{}{
					"id":           order["id"],
					"user_id":      order["user_id"],
					"product_id":   order["product_id"],
					"product_name": order["product_name"],
					"order_amount": order["order_amount"],
					"status":       order["status"],
					"created_at":   order["created_at"],
				}, nil
			}
			utils.Logger.Warn("订单不属于当前用户", zap.Uint("orderID", id), zap.Uint("userID", userID))
			return nil, ErrOrderNotFound
		}
	}

	utils.Logger.Warn("订单不存在", zap.Uint("id", id))
	return nil, ErrOrderNotFound
}

// GetOrderByIDAdmin 根据ID获取订单（管理员使用，不检查用户权限）
func GetOrderByIDAdmin(ctx context.Context, id uint) (map[string]interface{}, error) {
	var order model.Order
	if db.DB != nil {
		if err := db.DB.First(&order, id).Error; err == nil {
			var product model.Product
			db.DB.First(&product, order.ProductID)

			var user model.User
			db.DB.First(&user, order.UserID)

			return map[string]interface{}{
				"id":           order.ID,
				"user_id":      order.UserID,
				"user_name":    user.Username,
				"product_id":   order.ProductID,
				"product_name": product.Name,
				"order_amount": order.OrderAmount,
				"status":       order.Status,
				"created_at":   order.CreatedAt,
			}, nil
		}
	}

	mockOrdersMu.RLock()
	defer mockOrdersMu.RUnlock()

	for _, order := range mockOrders {
		if order["id"].(uint) == id {
			return map[string]interface{}{
				"id":           order["id"],
				"user_id":      order["user_id"],
				"user_name":    order["user_name"],
				"product_id":   order["product_id"],
				"product_name": order["product_name"],
				"order_amount": order["order_amount"],
				"status":       order["status"],
				"created_at":   order["created_at"],
			}, nil
		}
	}

	utils.Logger.Warn("订单不存在", zap.Uint("id", id))
	return nil, ErrOrderNotFound
}

// CancelOrder 取消订单
func CancelOrder(ctx context.Context, id uint, userID uint) error {
	var order model.Order
	if db.DB != nil {
		if err := db.DB.First(&order, id).Error; err == nil {
			if order.UserID != userID {
				utils.Logger.Warn("订单不属于当前用户", zap.Uint("orderID", id), zap.Uint("userID", userID))
				return ErrOrderNotFound
			}
			if order.Status != model.OrderStatusPending {
				utils.Logger.Warn("订单状态不允许取消", zap.Uint("orderID", id), zap.Int("status", int(order.Status)))
				return ErrOrderNotCancelable
			}
			order.Status = model.OrderStatusCanceled
			if err := db.DB.Save(&order).Error; err != nil {
				utils.Logger.Error("取消订单失败", zap.Uint("id", id), zap.Error(err))
				return err
			}
			utils.Logger.Info("订单取消成功", zap.Uint("id", id))
			return nil
		}
	}

	mockOrdersMu.Lock()
	defer mockOrdersMu.Unlock()

	for i, order := range mockOrders {
		if order["id"].(uint) == id {
			if order["user_id"].(uint) != userID {
				utils.Logger.Warn("订单不属于当前用户", zap.Uint("orderID", id), zap.Uint("userID", userID))
				return ErrOrderNotFound
			}
			if order["status"].(int) != 0 {
				utils.Logger.Warn("订单状态不允许取消", zap.Uint("orderID", id), zap.Int("status", order["status"].(int)))
				return ErrOrderNotCancelable
			}
			mockOrders[i]["status"] = 2
			utils.Logger.Info("订单取消成功", zap.Uint("id", id))
			return nil
		}
	}

	utils.Logger.Warn("订单不存在", zap.Uint("id", id))
	return ErrOrderNotFound
}

// ClearMockOrders 清空模拟订单数据（用于测试环境重置）
func ClearMockOrders() {
	mockOrdersMu.Lock()
	defer mockOrdersMu.Unlock()

	utils.Logger.Info("清空订单数据", zap.Int("count", len(mockOrders)))
	mockOrders = []map[string]interface{}{}
	nextOrderID = 1
}

// UpdateOrderStatus 更新订单状态（管理员使用）
func UpdateOrderStatus(ctx context.Context, id uint, status model.OrderStatus) error {
	var order model.Order
	if db.DB != nil {
		if err := db.DB.First(&order, id).Error; err == nil {
			order.Status = status
			if err := db.DB.Save(&order).Error; err != nil {
				utils.Logger.Error("更新订单状态失败", zap.Uint("id", id), zap.Error(err))
				return err
			}
			utils.Logger.Info("订单状态更新成功", zap.Uint("id", id), zap.Int("status", int(status)))
			return nil
		}
	}

	mockOrdersMu.Lock()
	defer mockOrdersMu.Unlock()

	for i, order := range mockOrders {
		if order["id"].(uint) == id {
			mockOrders[i]["status"] = int(status)
			utils.Logger.Info("订单状态更新成功", zap.Uint("id", id), zap.Int("status", int(status)))
			return nil
		}
	}

	utils.Logger.Warn("订单不存在", zap.Uint("id", id))
	return ErrOrderNotFound
}
