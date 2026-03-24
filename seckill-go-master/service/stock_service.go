package service

import (
	"context"
	"errors"
	"seckill_go/db"
	"seckill_go/utils"
	"strconv"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var (
	ErrStockNotFound    = errors.New("库存不存在")
	ErrStockNotEnough   = errors.New("库存不足")
	ErrStockPreheatFail = errors.New("库存预热失败")
)

const (
	StockKeyPrefix     = "seckill:stock:"
	StockLockPrefix    = "seckill:stock:lock:"
	StockReconcileKey  = "seckill:stock:reconcile"
	SeckillActivityKey = "seckill:activity:"
)

type StockService struct {
	redis *redis.Client
	mu    sync.RWMutex
}

var stockService *StockService
var stockServiceOnce sync.Once

func GetStockService() *StockService {
	stockServiceOnce.Do(func() {
		stockService = &StockService{
			redis: db.Redis,
		}
	})
	return stockService
}

func (s *StockService) PreheatStock(ctx context.Context, productID uint, stock int) error {
	if s.redis == nil {
		utils.Logger.Warn("Redis未连接，跳过库存预热", zap.Uint("productID", productID))
		return nil
	}

	stockKey := StockKeyPrefix + strconv.Itoa(int(productID))

	exists, err := s.redis.Exists(ctx, stockKey).Result()
	if err != nil {
		utils.Logger.Error("检查库存键失败", zap.Uint("productID", productID), zap.Error(err))
		return ErrStockPreheatFail
	}

	if exists > 0 {
		utils.Logger.Info("库存已存在，跳过预热", zap.Uint("productID", productID))
		return nil
	}

	_, err = s.redis.Set(ctx, stockKey, stock, 0).Result()
	if err != nil {
		utils.Logger.Error("库存预热失败", zap.Uint("productID", productID), zap.Error(err))
		return ErrStockPreheatFail
	}

	utils.Logger.Info("库存预热成功", zap.Uint("productID", productID), zap.Int("stock", stock))
	return nil
}

// SetSeckillActivity 设置秒杀活动
func (s *StockService) SetSeckillActivity(ctx context.Context, productID uint, stock int, startTime, endTime time.Time) error {
	if s.redis == nil {
		utils.Logger.Warn("Redis未连接，无法设置秒杀活动", zap.Uint("productID", productID))
		return errors.New("redis not connected")
	}

	activityKey := SeckillActivityKey + strconv.Itoa(int(productID))

	// 设置活动信息
	activityData := map[string]interface{}{
		"start_time": startTime.Unix(),
		"end_time":   endTime.Unix(),
		"stock":      stock,
		"status":     "pending",
	}

	// 存储活动信息
	for k, v := range activityData {
		key := activityKey + ":" + k
		_, err := s.redis.Set(ctx, key, v, 24*time.Hour).Result()
		if err != nil {
			utils.Logger.Error("设置秒杀活动失败", zap.Uint("productID", productID), zap.Error(err))
			return err
		}
	}

	// 预热库存
	err := s.PreheatStock(ctx, productID, stock)
	if err != nil {
		utils.Logger.Error("秒杀活动库存预热失败", zap.Uint("productID", productID), zap.Error(err))
		return err
	}

	utils.Logger.Info("秒杀活动设置成功",
		zap.Uint("productID", productID),
		zap.Int("stock", stock),
		zap.Time("startTime", startTime),
		zap.Time("endTime", endTime))
	return nil
}

// CheckSeckillActivity 检查秒杀活动状态
func (s *StockService) CheckSeckillActivity(ctx context.Context, productID uint) (bool, error) {
	if s.redis == nil {
		return false, errors.New("redis not connected")
	}

	activityKey := SeckillActivityKey + strconv.Itoa(int(productID))
	startTimeKey := activityKey + ":start_time"
	endTimeKey := activityKey + ":end_time"

	// 获取活动时间
	startTimeStr, err := s.redis.Get(ctx, startTimeKey).Result()
	if err != nil {
		if err == redis.Nil {
			return false, nil // 活动不存在
		}
		return false, err
	}

	endTimeStr, err := s.redis.Get(ctx, endTimeKey).Result()
	if err != nil {
		if err == redis.Nil {
			return false, nil // 活动不存在
		}
		return false, err
	}

	// 解析时间
	startTime, _ := strconv.ParseInt(startTimeStr, 10, 64)
	endTime, _ := strconv.ParseInt(endTimeStr, 10, 64)
	now := time.Now().Unix()

	// 检查活动是否在进行中
	isActive := now >= startTime && now <= endTime
	return isActive, nil
}

func (s *StockService) PreheatAllStocks(ctx context.Context) error {
	utils.Logger.Info("开始预热所有商品库存")

	if s.redis == nil {
		utils.Logger.Warn("Redis未连接，跳过库存预热")
		return nil
	}

	successCount := 0
	failCount := 0

	for _, product := range mockProducts {
		// 安全地转换ID
		var productID uint
		switch v := product["id"].(type) {
		case int:
			productID = uint(v)
		case float64:
			productID = uint(v)
		default:
			utils.Logger.Warn("无法转换商品ID，跳过", zap.Any("id", product["id"]))
			failCount++
			continue
		}

		// 安全地转换库存
		var stock int
		switch v := product["stock"].(type) {
		case int:
			stock = v
		case float64:
			stock = int(v)
		default:
			utils.Logger.Warn("无法转换商品库存，跳过", zap.Any("stock", product["stock"]))
			failCount++
			continue
		}

		if err := s.PreheatStock(ctx, productID, stock); err != nil {
			failCount++
		} else {
			successCount++
		}
	}

	utils.Logger.Info("库存预热完成",
		zap.Int("success", successCount),
		zap.Int("fail", failCount))

	if failCount > 0 {
		return ErrStockPreheatFail
	}
	return nil
}

func (s *StockService) DeductStock(ctx context.Context, productID uint) error {
	// 检查秒杀活动状态
	isActive, err := s.CheckSeckillActivity(ctx, productID)
	if err == nil && !isActive {
		// 检查是否是普通商品（非秒杀活动）
		activityKey := SeckillActivityKey + strconv.Itoa(int(productID))
		exists, _ := s.redis.Exists(ctx, activityKey+":start_time").Result()
		if exists > 0 {
			// 是秒杀商品但活动未开始或已结束
			return errors.New("秒杀活动未开始或已结束")
		}
	}

	// 如果Redis可用，优先使用Redis进行原子扣减
	if s.redis != nil {
		stockKey := StockKeyPrefix + strconv.Itoa(int(productID))

		// 先检查库存是否存在，不存在则尝试从内存同步
		exists, err := s.redis.Exists(ctx, stockKey).Result()
		if err != nil {
			utils.Logger.Error("检查Redis库存键失败", zap.Uint("productID", productID), zap.Error(err))
			return s.deductStockFromMemory(productID)
		}

		// 如果Redis中没有库存，尝试从内存同步
		if exists == 0 {
			memStock, err := s.getStockFromMemory(productID)
			if err != nil {
				return err
			}
			// 将内存库存同步到Redis
			_, err = s.redis.Set(ctx, stockKey, memStock, 0).Result()
			if err != nil {
				utils.Logger.Error("同步库存到Redis失败", zap.Uint("productID", productID), zap.Error(err))
				return s.deductStockFromMemory(productID)
			}
			utils.Logger.Info("库存已从内存同步到Redis", zap.Uint("productID", productID), zap.Int("stock", memStock))
		}

		luaScript := `
		local stockKey = KEYS[1]
		local lockKey = KEYS[2]
		local stock = tonumber(redis.call('get', stockKey))
		if stock == nil then
			return {-1, 0}
		elseif stock > 0 then
			-- 使用原子操作扣减库存
			local newStock = redis.call('decr', stockKey)
			-- 设置库存锁定，防止超卖
			redis.call('set', lockKey, 1, 'EX', 30)
			return {1, newStock}
		else
			return {0, 0}
		end
	`

		// 设置上下文超时，避免Redis操作阻塞
		redisCtx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
		defer cancel()

		lockKey := StockLockPrefix + strconv.Itoa(int(productID))
		result, err := s.redis.Eval(redisCtx, luaScript, []string{stockKey, lockKey}).Result()
		if err != nil {
			utils.Logger.Error("Redis执行Lua脚本失败", zap.Uint("productID", productID), zap.Error(err))
			// Redis失败时回退到内存库存扣减
			return s.deductStockFromMemory(productID)
		}

		// 类型断言安全检查
		resultSlice, ok := result.([]interface{})
		if !ok || len(resultSlice) < 2 {
			utils.Logger.Error("Redis返回结果格式错误", zap.Uint("productID", productID))
			return s.deductStockFromMemory(productID)
		}

		status, ok := resultSlice[0].(int64)
		if !ok {
			utils.Logger.Error("Redis返回结果类型错误", zap.Uint("productID", productID))
			return s.deductStockFromMemory(productID)
		}

		if status == -1 {
			utils.Logger.Error("Redis库存键不存在", zap.Uint("productID", productID))
			return s.deductStockFromMemory(productID)
		}

		if status == 0 {
			utils.Logger.Warn("库存不足", zap.Uint("productID", productID))
			return ErrStockNotEnough
		}

		newStock, _ := resultSlice[1].(int64)
		utils.Logger.Info("Redis库存扣减成功",
			zap.Uint("productID", productID),
			zap.Int64("newStock", newStock))

		return nil
	}

	// Redis不可用，使用内存扣减
	return s.deductStockFromMemory(productID)
}

func (s *StockService) deductStockFromMemory(productID uint) error {
	// 使用严格的互斥锁确保并发安全
	s.mu.Lock()
	defer s.mu.Unlock()

	// 查找商品
	for i := range mockProducts {
		if mockProducts[i]["id"] == int(productID) {
			// 安全地获取库存值
			var stock int
			switch v := mockProducts[i]["stock"].(type) {
			case int:
				stock = v
			case float64:
				stock = int(v)
			default:
				utils.Logger.Error("库存类型错误", zap.Uint("productID", productID), zap.Any("type", v))
				return errors.New("invalid stock type")
			}

			// 检查库存
			if stock <= 0 {
				return ErrStockNotEnough
			}

			// 扣减库存（在锁保护下完成）
			mockProducts[i]["stock"] = stock - 1

			utils.Logger.Info("内存库存扣减成功",
				zap.Uint("productID", productID),
				zap.Int("newStock", stock-1))
			return nil
		}
	}

	return ErrStockNotFound
}

func (s *StockService) GetStock(ctx context.Context, productID uint) (int, error) {
	if s.redis == nil {
		return s.getStockFromMemory(productID)
	}

	stockKey := StockKeyPrefix + strconv.Itoa(int(productID))
	stock, err := s.redis.Get(ctx, stockKey).Int()
	if err != nil {
		if err == redis.Nil {
			return s.getStockFromMemory(productID)
		}
		utils.Logger.Error("获取Redis库存失败", zap.Uint("productID", productID), zap.Error(err))
		return 0, err
	}

	return stock, nil
}

// ReleaseStock 释放库存（用于订单超时或失败）
func (s *StockService) ReleaseStock(ctx context.Context, productID uint) error {
	if s.redis == nil {
		// 内存库存释放
		s.mu.Lock()
		defer s.mu.Unlock()

		for i := range mockProducts {
			if mockProducts[i]["id"] == int(productID) {
				// 安全地获取库存值
				var stock int
				switch v := mockProducts[i]["stock"].(type) {
				case int:
					stock = v
				case float64:
					stock = int(v)
				default:
					utils.Logger.Error("库存类型错误", zap.Uint("productID", productID), zap.Any("type", v))
					return errors.New("invalid stock type")
				}

				// 释放库存
				mockProducts[i]["stock"] = stock + 1

				utils.Logger.Info("内存库存释放成功",
					zap.Uint("productID", productID),
					zap.Int("newStock", stock+1))
				return nil
			}
		}

		return ErrStockNotFound
	}

	stockKey := StockKeyPrefix + strconv.Itoa(int(productID))
	lockKey := StockLockPrefix + strconv.Itoa(int(productID))

	luaScript := `
		local stockKey = KEYS[1]
		local lockKey = KEYS[2]
		local stock = tonumber(redis.call('get', stockKey))
		if stock then
			local newStock = redis.call('incr', stockKey)
			-- 释放库存锁定
			redis.call('del', lockKey)
			return {1, newStock}
		else
			return {0, 0}
		end
	`

	result, err := s.redis.Eval(ctx, luaScript, []string{stockKey, lockKey}).Result()
	if err != nil {
		utils.Logger.Error("Redis释放库存失败", zap.Uint("productID", productID), zap.Error(err))
		return err
	}

	resultSlice := result.([]interface{})
	success := resultSlice[0].(int64)
	newStock := resultSlice[1].(int64)

	if success == 0 {
		utils.Logger.Warn("库存键不存在，无法释放", zap.Uint("productID", productID))
		return ErrStockNotFound
	}

	utils.Logger.Info("Redis库存释放成功",
		zap.Uint("productID", productID),
		zap.Int64("newStock", newStock))

	return nil
}

func (s *StockService) getStockFromMemory(productID uint) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, product := range mockProducts {
		if product["id"] == int(productID) {
			// 安全类型断言
			switch v := product["stock"].(type) {
			case int:
				return v, nil
			case float64:
				return int(v), nil
			default:
				return 0, errors.New("invalid stock type")
			}
		}
	}
	return 0, ErrStockNotFound
}

func (s *StockService) RestoreStock(ctx context.Context, productID uint, quantity int) error {
	if s.redis == nil {
		return s.restoreStockToMemory(productID, quantity)
	}

	stockKey := StockKeyPrefix + strconv.Itoa(int(productID))

	luaScript := `
		local stockKey = KEYS[1]
		local quantity = tonumber(ARGV[1])
		local stock = tonumber(redis.call('get', stockKey))
		if stock then
			local newStock = redis.call('incrby', stockKey, quantity)
			return {1, newStock}
		else
			return {0, 0}
		end
	`

	result, err := s.redis.Eval(ctx, luaScript, []string{stockKey}, quantity).Result()
	if err != nil {
		utils.Logger.Error("Redis恢复库存失败", zap.Uint("productID", productID), zap.Error(err))
		return s.restoreStockToMemory(productID, quantity)
	}

	resultSlice := result.([]interface{})
	success := resultSlice[0].(int64)
	newStock := resultSlice[1].(int64)

	if success == 0 {
		utils.Logger.Warn("库存键不存在，无法恢复", zap.Uint("productID", productID))
		return ErrStockNotFound
	}

	utils.Logger.Info("库存恢复成功",
		zap.Uint("productID", productID),
		zap.Int("quantity", quantity),
		zap.Int64("newStock", newStock))

	return nil
}

func (s *StockService) restoreStockToMemory(productID uint, quantity int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range mockProducts {
		if mockProducts[i]["id"] == int(productID) {
			stock := mockProducts[i]["stock"].(int)
			mockProducts[i]["stock"] = stock + quantity
			utils.Logger.Info("内存库存恢复成功",
				zap.Uint("productID", productID),
				zap.Int("quantity", quantity),
				zap.Int("newStock", stock+quantity))
			return nil
		}
	}
	return ErrStockNotFound
}

func (s *StockService) ReconcileStock(ctx context.Context) ([]map[string]interface{}, error) {
	utils.Logger.Info("开始库存对账")

	if s.redis == nil {
		utils.Logger.Warn("Redis未连接，跳过库存对账")
		return nil, nil
	}

	var inconsistencies []map[string]interface{}

	for _, product := range mockProducts {
		productID := uint(product["id"].(int))
		dbStock := product["stock"].(int)

		redisStock, err := s.GetStock(ctx, productID)
		if err != nil {
			utils.Logger.Error("获取Redis库存失败",
				zap.Uint("productID", productID),
				zap.Error(err))
			continue
		}

		if redisStock != dbStock {
			inconsistency := map[string]interface{}{
				"product_id":  productID,
				"db_stock":    dbStock,
				"redis_stock": redisStock,
				"difference":  dbStock - redisStock,
				"timestamp":   time.Now().Format(time.RFC3339),
			}
			inconsistencies = append(inconsistencies, inconsistency)

			utils.Logger.Warn("库存不一致",
				zap.Uint("productID", productID),
				zap.Int("dbStock", dbStock),
				zap.Int("redisStock", redisStock),
				zap.Int("difference", dbStock-redisStock))
		}
	}

	utils.Logger.Info("库存对账完成",
		zap.Int("totalProducts", len(mockProducts)),
		zap.Int("inconsistencies", len(inconsistencies)))

	return inconsistencies, nil
}

func (s *StockService) AutoFixStock(ctx context.Context, productID uint) error {
	utils.Logger.Info("开始自动修复库存", zap.Uint("productID", productID))

	if s.redis == nil {
		utils.Logger.Warn("Redis未连接，无法自动修复")
		return nil
	}

	var dbStock int
	found := false
	for _, product := range mockProducts {
		if product["id"] == int(productID) {
			dbStock = product["stock"].(int)
			found = true
			break
		}
	}

	if !found {
		utils.Logger.Warn("商品不存在", zap.Uint("productID", productID))
		return ErrStockNotFound
	}

	stockKey := StockKeyPrefix + strconv.Itoa(int(productID))
	_, err := s.redis.Set(ctx, stockKey, dbStock, 0).Result()
	if err != nil {
		utils.Logger.Error("自动修复库存失败",
			zap.Uint("productID", productID),
			zap.Error(err))
		return err
	}

	utils.Logger.Info("库存自动修复成功",
		zap.Uint("productID", productID),
		zap.Int("stock", dbStock))

	return nil
}

func (s *StockService) StartReconciliationTask(interval time.Duration) {
	if s.redis == nil {
		utils.Logger.Warn("Redis未连接，不启动库存对账任务")
		return
	}

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
			ctx := context.Background()
			inconsistencies, err := s.ReconcileStock(ctx)
			if err != nil {
				utils.Logger.Error("库存对账任务失败", zap.Error(err))
				continue
			}

			for _, inc := range inconsistencies {
				productID := inc["product_id"].(uint)
				if err := s.AutoFixStock(ctx, productID); err != nil {
					utils.Logger.Error("自动修复库存失败",
						zap.Uint("productID", productID),
						zap.Error(err))
				}
			}
		}
	}()

	utils.Logger.Info("库存对账任务已启动", zap.Duration("interval", interval))
}

func (s *StockService) ResetStock(ctx context.Context, productID uint, stock int) error {
	utils.Logger.Info("重置商品库存", zap.Uint("productID", productID), zap.Int("stock", stock))

	s.mu.Lock()
	defer s.mu.Unlock()

	found := false
	for i := range mockProducts {
		if mockProducts[i]["id"] == int(productID) {
			mockProducts[i]["stock"] = stock
			found = true
			break
		}
	}

	if !found {
		return ErrStockNotFound
	}

	if s.redis != nil {
		stockKey := StockKeyPrefix + strconv.Itoa(int(productID))
		_, err := s.redis.Set(ctx, stockKey, stock, 0).Result()
		if err != nil {
			utils.Logger.Error("Redis设置库存失败",
				zap.Uint("productID", productID),
				zap.Error(err))
			return err
		}
		utils.Logger.Info("库存已同步到Redis",
			zap.Uint("productID", productID),
			zap.Int("stock", stock))
	}

	return nil
}
