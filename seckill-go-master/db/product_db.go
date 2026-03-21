package db

import (
	"context"
	"errors"
	"fmt"
	"seckill_go/model"
	"seckill_go/utils"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// 商品哈希表键前缀：seckill:product:{id}
const productHashKeyPrefix = "seckill:product:%d"

// 生成商品哈希表键
func getProductHashKey(productID uint64) string {
	return fmt.Sprintf(productHashKeyPrefix, productID)
}

// DecrStock 直接操作哈希表中的stock字段扣减库存（原子操作）
// 仅使用哈希表存储，不依赖独立库存键
func DecrStock(ctx context.Context, productID uint64, num int) (bool, error) {
	// 入参校验
	if num <= 0 {
		utils.Logger.Warn("无效的库存扣减数量", zap.Int("num", num), zap.Uint64("product_id", productID))
		return false, fmt.Errorf("扣减数量必须大于0")
	}
	if productID == 0 {
		utils.Logger.Warn("无效的商品ID", zap.Uint64("product_id", productID))
		return false, fmt.Errorf("商品ID不能为空")
	}
	// LUA脚本：直接操作哈希表的stock字段，实现原子扣减
	luaScript := `
    -- 入参：
    -- KEYS[1]：商品哈希表键（seckill:product:{id}）
    -- ARGV[1]：扣减数量
    
    local hashKey = KEYS[1]
    local decrNum = tonumber(ARGV[1])
    
    -- 1. 从哈希表获取当前库存（字段不存在默认0）
    local currentStock = tonumber(redis.call('HGET', hashKey, 'stock') or "0")
    
    -- 2. 校验库存是否充足
    if currentStock < decrNum then
        return 0  -- 库存不足
    end
    
    -- 3. 原子扣减库存（直接修改哈希表的stock字段）
    redis.call('HINCRBY', hashKey, 'stock', -decrNum)
    
    -- 4. 返回成功标识
    return 1
	`
	// 准备脚本参数：KEYS传入商品哈希表键，ARGV传入扣减数量
	hashKey := getProductHashKey(productID)
	keys := []string{hashKey}
	args := []interface{}{num}
	// 执行LUA脚本（保证原子性）
	result, err := Redis.Eval(ctx, luaScript, keys, args).Int()
	if err != nil {
		utils.Logger.Error("库存扣减失败",
			zap.Uint64("product_id", productID),
			zap.String("hash_key", hashKey),
			zap.Error(err))
		return false, fmt.Errorf("redis操作失败：%w", err)
	}
	// 解析结果
	if result == 1 {
		utils.Logger.Info("库存扣减成功",
			zap.Uint64("product_id", productID),
			zap.Int("扣减数量", num),
			zap.String("hash_key", hashKey))
		return true, nil
	}

	// 库存不足时，查询当前实际库存并记录日志（便于排查）
	currentStock, _ := Redis.HGet(ctx, hashKey, "stock").Int()
	utils.Logger.Warn("库存不足",
		zap.Uint64("product_id", productID),
		zap.Int("请求扣减数量", num),
		zap.Int("当前库存", currentStock),
		zap.String("hash_key", hashKey))
	return false, nil
}

// GetAllProductIDs 从数据库获取全量商品主键ID
func GetAllProductIDs() ([]uint64, error) {
	var ids []uint64
	// 仅查主键ID，不查其他字段
	result := DB.Model(&model.Product{}).Select("id").Find(&ids)
	if result.Error != nil {
		utils.Logger.Error("数据库查询全量商品ID失败", zap.Error(result.Error))
		return nil, result.Error
	}
	return ids, nil
}

// GetProductList 根据id获取商品列表
func GetProductList(ctx context.Context, ids []uint64) ([]map[string]interface{}, error) {
	// 结果集
	var result []map[string]interface{}
	// 缓存未命中的ID列表
	var missIDs []uint64
	// 1. 先从缓存获取每个商品的信息
	for _, id := range ids {
		productKey := "seckill:product:" + strconv.FormatUint(id, 10)
		// 从Redis哈希表获取商品信息
		productData, err := Redis.HGetAll(ctx, productKey).Result()
		if err != nil && !errors.Is(err, redis.Nil) {
			utils.Logger.Warn("从缓存获取单个商品信息失败",
				zap.Uint64("product_id", id),
				zap.Error(err))
			missIDs = append(missIDs, id)
			continue
		}

		// 检查是否命中缓存
		if len(productData) > 0 {
			// 转换为map[string]interface{}
			productMap, err := convertHashToMap(id, productData)
			if err != nil {
				utils.Logger.Warn("转换缓存商品信息失败",
					zap.Uint64("product_id", id),
					zap.Error(err))
				missIDs = append(missIDs, id)
				continue
			}
			result = append(result, productMap)
		} else {
			// 缓存未命中，记录需要从数据库查询的ID
			missIDs = append(missIDs, id)
		}
	}

	utils.Logger.Info("缓存查询结果",
		zap.Int("hit_count", len(ids)-len(missIDs)),
		zap.Int("miss_count", len(missIDs)))

	// 2. 处理缓存未命中的商品，从数据库查询
	if len(missIDs) > 0 {
		utils.Logger.Info("缓存未命中，从数据库获取商品信息", zap.Int("miss_count", len(missIDs)))

		var dbProducts []model.Product
		// 从数据库查询未命中的商品
		resultDB := DB.Where("id IN (?)", missIDs).Find(&dbProducts)
		if resultDB.Error != nil {
			utils.Logger.Error("从数据库获取商品信息失败", zap.Error(resultDB.Error))
			return nil, resultDB.Error
		}

		// 处理查询结果
		dbResultMap := make(map[uint64]model.Product)
		for _, p := range dbProducts {
			dbResultMap[p.ID] = p

			// 转换为map并添加到结果集
			productMap := map[string]interface{}{
				"id":    p.ID,
				"name":  p.Name,
				"price": p.Price,
				"stock": p.Stock,
				//"version": p.Version,
			}
			result = append(result, productMap)

			// 同步到Redis缓存
			productKey := "seckill:product:" + strconv.FormatUint(uint64(p.ID), 10)
			if err := Redis.HSet(ctx, productKey, productMap).Err(); err != nil {
				utils.Logger.Warn("同步商品信息到缓存失败",
					zap.String("key", productKey),
					zap.Error(err))
			} else {
				// 设置过期时间
				_ = Redis.Expire(ctx, productKey, 2*time.Hour).Err()
				utils.Logger.Debug("同步商品信息到缓存成功", zap.String("key", productKey))
			}
		}

		// 检查是否有ID在数据库中也不存在
		for _, id := range missIDs {
			if _, exists := dbResultMap[id]; !exists {
				utils.Logger.Warn("商品在数据库中不存在", zap.Uint64("product_id", id))
			}
		}
	}

	// 3. 按原始ID顺序排序结果
	sortedResult := sortResultByIDs(ids, result)

	return sortedResult, nil
}

// convertHashToMap 将Redis哈希表结果转换为map[string]interface{}
func convertHashToMap(id uint64, hashData map[string]string) (map[string]interface{}, error) {
	price, err := strconv.ParseFloat(hashData["price"], 64)
	if err != nil {
		return nil, err
	}

	stock, err := strconv.Atoi(hashData["stock"])
	if err != nil {
		return nil, err
	}

	//version, err := strconv.Atoi(hashData["version"])
	//if err != nil {
	//	return nil, err
	//}

	return map[string]interface{}{
		"id":    id,
		"name":  hashData["name"],
		"price": price,
		"stock": stock,
		//"version": version,
	}, nil
}

// sortResultByIDs 按照原始ID列表的顺序对结果进行排序
func sortResultByIDs(ids []uint64, result []map[string]interface{}) []map[string]interface{} {
	// 创建ID到结果的映射
	resultMap := make(map[uint64]map[string]interface{})
	for _, item := range result {
		id := item["id"].(uint64)
		resultMap[id] = item
	}

	// 按照原始ID顺序构建结果集
	sorted := make([]map[string]interface{}, 0, len(ids))
	for _, id := range ids {
		if item, exists := resultMap[id]; exists {
			sorted = append(sorted, item)
		}
	}

	return sorted
}
