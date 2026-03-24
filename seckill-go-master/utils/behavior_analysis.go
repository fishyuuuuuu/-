package utils

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// 行为分析相关常量
const (
	// 行为类型
	BehaviorTypeLogin     = "login"
	BehaviorTypeSeckill   = "seckill"
	BehaviorTypePageView  = "page_view"
	BehaviorTypeAddToCart = "add_to_cart"
	BehaviorTypeCheckout  = "checkout"

	// 行为时间窗口
	BehaviorWindowShort  = 1 * time.Minute
	BehaviorWindowMedium = 5 * time.Minute
	BehaviorWindowLong   = 10 * time.Minute

	// 风险评分阈值
	RiskScoreLow    = 30
	RiskScoreMedium = 60
	RiskScoreHigh   = 80

	// 行为键前缀
	BehaviorKeyPrefix = "security:behavior:"
)

// 用户行为记录
type UserBehavior struct {
	UserID    string    `json:"user_id"`
	IP        string    `json:"ip"`
	Behavior  string    `json:"behavior"`
	Timestamp time.Time `json:"timestamp"`
	Details   string    `json:"details"`
}

// 行为分析器
type BehaviorAnalyzer struct {
	redis         *redis.Client
	mu            sync.RWMutex
	behaviorCache map[string][]UserBehavior
}

var behaviorAnalyzer *BehaviorAnalyzer
var behaviorAnalyzerOnce sync.Once

// GetBehaviorAnalyzer 获取行为分析器实例
func GetBehaviorAnalyzer() *BehaviorAnalyzer {
	behaviorAnalyzerOnce.Do(func() {
		behaviorAnalyzer = &BehaviorAnalyzer{
			redis:         securityRedis,
			behaviorCache: make(map[string][]UserBehavior),
		}

		// 启动清理任务
		go behaviorAnalyzer.cleanupTask()
	})
	return behaviorAnalyzer
}

// RecordBehavior 记录用户行为
func (ba *BehaviorAnalyzer) RecordBehavior(c *gin.Context, behaviorType string, details string) {
	clientIP := c.ClientIP()
	userID := "anonymous"

	// 获取用户ID（如果已登录）
	if uid, exists := c.Get("userID"); exists {
		if id, ok := uid.(uint); ok {
			userID = fmt.Sprintf("%d", id)
		}
	}

	// 创建行为记录
	behavior := UserBehavior{
		UserID:    userID,
		IP:        clientIP,
		Behavior:  behaviorType,
		Timestamp: time.Now(),
		Details:   details,
	}

	// 存储到Redis
	if ba.redis != nil {
		ctx := context.Background()
		key := fmt.Sprintf("%s%s:%s", BehaviorKeyPrefix, clientIP, behaviorType)

		// 存储行为记录
		_, err := ba.redis.ZAdd(ctx, key, redis.Z{
			Score:  float64(behavior.Timestamp.Unix()),
			Member: fmt.Sprintf("%s|%s|%s", behavior.UserID, behavior.Timestamp.Format(time.RFC3339), behavior.Details),
		}).Result()

		if err != nil {
			Logger.Warn("存储行为记录失败", zap.Error(err))
		}

		// 设置过期时间
		ba.redis.Expire(ctx, key, BehaviorWindowLong)
	}

	// 存储到本地缓存
	cacheKey := fmt.Sprintf("%s:%s", clientIP, behaviorType)
	ba.mu.Lock()
	ba.behaviorCache[cacheKey] = append(ba.behaviorCache[cacheKey], behavior)
	// 限制缓存大小
	if len(ba.behaviorCache[cacheKey]) > 100 {
		ba.behaviorCache[cacheKey] = ba.behaviorCache[cacheKey][len(ba.behaviorCache[cacheKey])-100:]
	}
	ba.mu.Unlock()
}

// AnalyzeBehavior 分析用户行为，返回风险评分（0-100）
func (ba *BehaviorAnalyzer) AnalyzeBehavior(c *gin.Context) int {
	clientIP := c.ClientIP()
	userID := "anonymous"

	// 获取用户ID（如果已登录）
	if uid, exists := c.Get("userID"); exists {
		if id, ok := uid.(uint); ok {
			userID = fmt.Sprintf("%d", id)
		}
	}

	riskScore := 0

	// 1. 检查短时间内的请求频率
	riskScore += ba.checkRequestFrequency(clientIP, userID, BehaviorWindowShort)

	// 2. 检查行为模式异常
	riskScore += ba.checkBehaviorPattern(clientIP, userID)

	// 3. 检查秒杀行为异常
	riskScore += ba.checkSeckillBehavior(clientIP, userID)

	// 4. 检查IP异常（如代理IP、异常地理位置等）
	riskScore += ba.checkIPAnomaly(clientIP)

	// 限制风险评分范围
	if riskScore > 100 {
		riskScore = 100
	}

	return riskScore
}

// checkRequestFrequency 检查请求频率
func (ba *BehaviorAnalyzer) checkRequestFrequency(ip, userID string, window time.Duration) int {
	score := 0
	ctx := context.Background()

	// 检查IP级别的请求频率
	ipKey := fmt.Sprintf("%s%s:*", BehaviorKeyPrefix, ip)
	if ba.redis != nil {
		// 查找所有相关的行为键
		keys, err := ba.redis.Keys(ctx, ipKey).Result()
		if err == nil {
			totalCount := 0
			for _, key := range keys {
				// 计算时间窗口内的行为数量
				minScore := float64(time.Now().Add(-window).Unix())
				count, _ := ba.redis.ZCount(ctx, key, fmt.Sprintf("%f", minScore), "+inf").Result()
				totalCount += int(count)
			}

			// 根据请求频率计算风险分数（提高阈值）
			if totalCount > 30 {
				score += 60 // 提高风险分数
			} else if totalCount > 20 {
				score += 45
			} else if totalCount > 10 {
				score += 30
			} else if totalCount > 5 {
				score += 15
			}
		}
	}

	return score
}

// checkBehaviorPattern 检查行为模式
func (ba *BehaviorAnalyzer) checkBehaviorPattern(ip, userID string) int {
	score := 0
	ctx := context.Background()

	// 检查是否有异常的行为序列
	// 例如：直接访问秒杀接口而没有浏览商品页面
	pageViewKey := fmt.Sprintf("%s%s:%s", BehaviorKeyPrefix, ip, BehaviorTypePageView)
	seckillKey := fmt.Sprintf("%s%s:%s", BehaviorKeyPrefix, ip, BehaviorTypeSeckill)

	if ba.redis != nil {
		// 检查是否有秒杀行为
		seckillCount, _ := ba.redis.ZCard(ctx, seckillKey).Result()
		if seckillCount > 0 {
			// 检查是否有页面浏览行为
			pageViewCount, _ := ba.redis.ZCard(ctx, pageViewKey).Result()
			if pageViewCount == 0 {
				// 没有浏览页面直接秒杀，可能是脚本
				score += 50 // 提高风险分数
			} else if pageViewCount < seckillCount {
				// 浏览页面次数少于秒杀次数，可能是脚本
				score += 30
			}
		}
	}

	return score
}

// checkSeckillBehavior 检查秒杀行为
func (ba *BehaviorAnalyzer) checkSeckillBehavior(ip, userID string) int {
	score := 0
	ctx := context.Background()

	// 检查秒杀行为的时间分布
	seckillKey := fmt.Sprintf("%s%s:%s", BehaviorKeyPrefix, ip, BehaviorTypeSeckill)

	if ba.redis != nil {
		// 获取最近的秒杀行为
		behaviors, err := ba.redis.ZRevRange(ctx, seckillKey, 0, 9).Result()
		if err == nil && len(behaviors) > 0 {
			// 检查秒杀行为的时间间隔
			var timestamps []time.Time
			for _, behavior := range behaviors {
				parts := splitBehaviorString(behavior)
				if len(parts) >= 2 {
					timestamp, err := time.Parse(time.RFC3339, parts[1])
					if err == nil {
						timestamps = append(timestamps, timestamp)
					}
				}
			}

			// 检查时间间隔是否异常（过于规律）
			if len(timestamps) > 1 {
				intervals := make([]time.Duration, 0)
				for i := 0; i < len(timestamps)-1; i++ {
					interval := timestamps[i].Sub(timestamps[i+1])
					if interval > 0 {
						intervals = append(intervals, interval)
					}
				}

				// 计算时间间隔的标准差
				if len(intervals) > 1 {
					stdDev := calculateStdDev(intervals)
					// 如果时间间隔非常规律，可能是脚本
					if stdDev < time.Millisecond*100 {
						score += 35
					}
				}
			}
		}
	}

	return score
}

// checkIPAnomaly 检查IP异常
func (ba *BehaviorAnalyzer) checkIPAnomaly(ip string) int {
	score := 0

	// 这里可以添加IP异常检测逻辑
	// 例如：检查是否是代理IP、是否来自异常地理位置等
	// 暂时简单实现

	return score
}

// GetRiskLevel 根据风险评分获取风险等级
func (ba *BehaviorAnalyzer) GetRiskLevel(score int) string {
	if score >= RiskScoreHigh {
		return "high"
	} else if score >= RiskScoreMedium {
		return "medium"
	} else if score >= RiskScoreLow {
		return "low"
	}
	return "very_low"
}

// cleanupTask 定期清理过期的行为记录
func (ba *BehaviorAnalyzer) cleanupTask() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		ba.mu.Lock()
		// 清理本地缓存中过期的行为记录
		now := time.Now()
		for key, behaviors := range ba.behaviorCache {
			filtered := make([]UserBehavior, 0)
			for _, behavior := range behaviors {
				if now.Sub(behavior.Timestamp) < BehaviorWindowLong {
					filtered = append(filtered, behavior)
				}
			}
			if len(filtered) == 0 {
				delete(ba.behaviorCache, key)
			} else {
				ba.behaviorCache[key] = filtered
			}
		}
		ba.mu.Unlock()
	}
}

// 辅助函数：分割行为字符串
func splitBehaviorString(s string) []string {
	parts := make([]string, 0)
	current := ""
	for i := 0; i < len(s); i++ {
		if s[i] == '|' {
			parts = append(parts, current)
			current = ""
		} else {
			current += string(s[i])
		}
	}
	if current != "" {
		parts = append(parts, current)
	}
	return parts
}

// 辅助函数：计算时间间隔的标准差
func calculateStdDev(intervals []time.Duration) time.Duration {
	if len(intervals) == 0 {
		return 0
	}

	// 计算平均值
	total := time.Duration(0)
	for _, interval := range intervals {
		total += interval
	}
	mean := total / time.Duration(len(intervals))

	// 计算方差
	variance := time.Duration(0)
	for _, interval := range intervals {
		diff := interval - mean
		variance += diff * diff
	}
	variance /= time.Duration(len(intervals))

	// 计算标准差（简化计算）
	stdDev := time.Duration(int64(variance) / 2)
	return stdDev
}

// BehaviorAnalysisMiddleware 行为分析中间件
func BehaviorAnalysisMiddleware() gin.HandlerFunc {
	ba := GetBehaviorAnalyzer()

	return func(c *gin.Context) {
		// 记录页面访问行为
		if c.Request.Method == "GET" {
			ba.RecordBehavior(c, BehaviorTypePageView, c.Request.URL.Path)
		}

		// 记录秒杀行为
		if c.Request.Method == "POST" && c.Request.URL.Path == "/api/product/seckill" {
			ba.RecordBehavior(c, BehaviorTypeSeckill, c.Request.URL.Path)
		}

		// 分析用户行为，计算风险评分
		riskScore := ba.AnalyzeBehavior(c)
		riskLevel := ba.GetRiskLevel(riskScore)

		// 将风险信息存储到上下文
		c.Set("risk_score", riskScore)
		c.Set("risk_level", riskLevel)

		// 对于高风险用户，直接拒绝请求
		if riskLevel == "high" {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":      "请求过于频繁，请稍后再试",
				"risk_level": riskLevel,
				"risk_score": riskScore,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
