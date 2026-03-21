package service

import (
	"context"
	"errors"
	"seckill_go/rabbitmq"
	"seckill_go/utils"
	"strconv"

	"go.uber.org/zap"
)

// ErrProductNotFound 自定义错误
var (
	ErrProductNotFound = errors.New("商品不存在")
)

// 模拟商品数据
var mockProducts = []map[string]interface{}{
	{
		"id":    1,
		"name":  "iPhone 15 Pro",
		"price": 7999.99,
		"stock": 100,
	},
	{
		"id":    2,
		"name":  "MacBook Pro 14",
		"price": 12999.99,
		"stock": 50,
	},
	{
		"id":    3,
		"name":  "AirPods Pro 2",
		"price": 1899.99,
		"stock": 200,
	},
	{
		"id":    4,
		"name":  "iPad Pro 12.9",
		"price": 8999.99,
		"stock": 50,
	},
}

// ReduceProductStock 扣减商品库存（使用新的库存服务）
func ReduceProductStock(ctx context.Context, productID uint, userID string) error {
	stockService := GetStockService()

	err := stockService.DeductStock(ctx, productID)
	if err != nil {
		if errors.Is(err, ErrStockNotEnough) {
			utils.Logger.Error("库存不足", zap.Uint("productID", productID))
			return errors.New("库存不足")
		}
		if errors.Is(err, ErrStockNotFound) {
			utils.Logger.Error("商品不存在", zap.Uint("productID", productID))
			return ErrProductNotFound
		}
		utils.Logger.Error("库存扣减失败", zap.Uint("productID", productID), zap.Error(err))
		return err
	}

	err = PublishOrderMessageReliably(ctx, userID, strconv.Itoa(int(productID)))
	if err != nil {
		utils.Logger.Error("订单消息发送失败，库存已扣减",
			zap.Uint("productID", productID),
			zap.String("userID", userID),
			zap.Error(err))
	}

	return nil
}

// PublishOrderCreationMsg 发送订单创建消息
func PublishOrderCreationMsg(id string, productId string) error {
	utils.Logger.Info("发送订单创建消息", zap.String("id", id), zap.String("productId", productId))

	// 检查RabbitMQ是否连接
	rmq, err := rabbitmq.GetInstance()
	if err != nil {
		utils.Logger.Error("获取RabbitMQ实例失败", zap.Error(err))
		// 失败也不影响流程
		return nil
	}

	// 声明队列
	err = rmq.DeclareQueue("seckill_order_queue")
	if err != nil {
		utils.Logger.Error("声明队列失败", zap.Error(err))
		// 失败也不影响流程
		return nil
	}

	// 发送消息
	message := id + ":" + productId
	err = rmq.SendToQueue("seckill_order_queue", []byte(message))
	if err != nil {
		utils.Logger.Error("发送消息失败", zap.Error(err))
		// 失败也不影响流程
		return nil
	}

	utils.Logger.Info("消息发送成功", zap.String("message", message))
	return nil
}

// GetProductList 获取所有商品列表
func GetProductList(ctx context.Context) ([]map[string]interface{}, error) {
	utils.Logger.Info("获取商品列表", zap.Int("count", len(mockProducts)))

	stockService := GetStockService()
	if err := stockService.PreheatAllStocks(ctx); err != nil {
		utils.Logger.Warn("库存预热部分失败", zap.Error(err))
	}

	return mockProducts, nil
}

// GetProductInfo 先查缓存，缓存未命中则查库并同步到缓存
// ids: 需要查询的商品ID列表
func GetProductInfo(ctx context.Context, ids []uint64) ([]map[string]interface{}, error) {
	utils.Logger.Info("处理获取指定商品信息请求", zap.Int("id_count", len(ids)))

	// 检查输入ID是否有效
	if len(ids) == 0 {
		utils.Logger.Warn("未提供有效的商品ID列表")
		return nil, errors.New("无效的商品ID列表")
	}

	// 模拟根据ID获取商品
	var result []map[string]interface{}
	for _, id := range ids {
		for _, product := range mockProducts {
			if product["id"] == id {
				result = append(result, product)
				break
			}
		}
	}

	// 检查是否有商品数据返回
	if len(result) == 0 {
		utils.Logger.Warn("未找到任何商品", zap.Any("ids", ids))
		return nil, ErrProductNotFound // 返回自定义的"商品不存在"错误
	}

	// 记录成功获取的商品数量
	utils.Logger.Info("成功获取商品信息",
		zap.Int("requested_count", len(ids)),
		zap.Int("returned_count", len(result)))

	return result, nil
}

// CreateProduct 创建商品
func CreateProduct(ctx context.Context, product map[string]interface{}) error {
	// 生成新的商品ID
	newID := uint64(len(mockProducts) + 1)
	product["id"] = newID

	// 添加到模拟数据
	mockProducts = append(mockProducts, product)

	utils.Logger.Info("商品创建成功", zap.Uint64("id", newID), zap.String("name", product["name"].(string)))
	return nil
}

// UpdateProduct 更新商品
func UpdateProduct(ctx context.Context, id uint64, product map[string]interface{}) error {
	// 查找商品
	for i, p := range mockProducts {
		if p["id"] == id {
			// 更新商品信息
			if name, ok := product["name"]; ok {
				mockProducts[i]["name"] = name
			}
			if price, ok := product["price"]; ok {
				mockProducts[i]["price"] = price
			}
			if stock, ok := product["stock"]; ok {
				mockProducts[i]["stock"] = stock
			}

			utils.Logger.Info("商品更新成功", zap.Uint64("id", id))
			return nil
		}
	}

	utils.Logger.Warn("商品不存在", zap.Uint64("id", id))
	return ErrProductNotFound
}

// DeleteProduct 删除商品
func DeleteProduct(ctx context.Context, id uint64) error {
	// 查找并删除商品
	for i, p := range mockProducts {
		if p["id"] == id {
			// 从切片中删除元素
			mockProducts = append(mockProducts[:i], mockProducts[i+1:]...)
			utils.Logger.Info("商品删除成功", zap.Uint64("id", id))
			return nil
		}
	}

	utils.Logger.Warn("商品不存在", zap.Uint64("id", id))
	return ErrProductNotFound
}

// GetProductByID 根据ID获取商品
func GetProductByID(ctx context.Context, id uint64) (map[string]interface{}, error) {
	// 查找商品
	for _, product := range mockProducts {
		if product["id"] == id {
			return product, nil
		}
	}

	utils.Logger.Warn("商品不存在", zap.Uint64("id", id))
	return nil, ErrProductNotFound
}

// ResetProductStock 重置商品库存（用于性能测试）
func ResetProductStock(ctx context.Context, productID uint, stock int) error {
	stockService := GetStockService()
	return stockService.ResetStock(ctx, productID, stock)
}
