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
	StockKeyPrefix    = "seckill:stock:"
	StockLockPrefix   = "seckill:stock:lock:"
	StockReconcileKey = "seckill:stock:reconcile"
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

func (s *StockService) PreheatAllStocks(ctx context.Context) error {
	utils.Logger.Info("开始预热所有商品库存")

	if s.redis == nil {
		utils.Logger.Warn("Redis未连接，跳过库存预热")
		return nil
	}

	successCount := 0
	failCount := 0

	for _, product := range mockProducts {
		productID := uint(product["id"].(int))
		stock := product["stock"].(int)

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
	if s.redis == nil {
		return s.deductStockFromMemory(productID)
	}

	stockKey := StockKeyPrefix + strconv.Itoa(int(productID))

	luaScript := `
		local stockKey = KEYS[1]
		local stock = tonumber(redis.call('get', stockKey))
		if stock and stock > 0 then
			local newStock = redis.call('decr', stockKey)
			return {1, newStock}
		else
			return {0, 0}
		end
	`

	result, err := s.redis.Eval(ctx, luaScript, []string{stockKey}).Result()
	if err != nil {
		utils.Logger.Error("Redis执行Lua脚本失败", zap.Uint("productID", productID), zap.Error(err))
		return s.deductStockFromMemory(productID)
	}

	resultSlice := result.([]interface{})
	success := resultSlice[0].(int64)
	newStock := resultSlice[1].(int64)

	if success == 0 {
		utils.Logger.Warn("库存不足", zap.Uint("productID", productID))
		return ErrStockNotEnough
	}

	utils.Logger.Info("库存扣减成功",
		zap.Uint("productID", productID),
		zap.Int64("newStock", newStock))

	return nil
}

func (s *StockService) deductStockFromMemory(productID uint) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range mockProducts {
		if mockProducts[i]["id"] == int(productID) {
			stock := mockProducts[i]["stock"].(int)
			if stock > 0 {
				mockProducts[i]["stock"] = stock - 1
				utils.Logger.Info("内存库存扣减成功",
					zap.Uint("productID", productID),
					zap.Int("newStock", stock-1))
				return nil
			}
			return ErrStockNotEnough
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

func (s *StockService) getStockFromMemory(productID uint) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, product := range mockProducts {
		if product["id"] == int(productID) {
			return product["stock"].(int), nil
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
