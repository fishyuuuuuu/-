package utils

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// 安全相关常量
const (
	// 限流配置
	GlobalRateLimitCapacity = 10000 // 全局令牌桶容量
	GlobalRateLimitRate     = 5000  // 全局每秒生成令牌数
	UserRateLimitCapacity   = 100   // 用户级令牌桶容量
	UserRateLimitRate       = 10    // 用户级每秒生成令牌数
	IPRateLimitCapacity     = 200   // IP级令牌桶容量
	IPRateLimitRate         = 20    // IP级每秒生成令牌数

	// 黑名单配置
	BlacklistKeyPrefix     = "security:blacklist:"
	BlacklistExpireBase    = 10 * time.Minute // 首次封禁10分钟
	BlacklistExpireStep    = 20 * time.Minute // 每次增加20分钟
	BlacklistMaxExpire     = 2 * time.Hour    // 最大封禁2小时
	BlacklistFailThreshold = 3                // 触发黑名单的失败次数阈值

	// 异常行为阈值
	IPFailThreshold      = 50 // IP 5分钟内失败次数阈值
	IPRepeatThreshold    = 20 // IP 1分钟内重复请求次数阈值
	UserCaptchaThreshold = 3  // 用户连续验证码错误次数阈值
	UserOrderThreshold   = 3  // 用户10分钟内订单创建失败次数阈值
)

// 全局Redis客户端（由main.go初始化时设置）
var securityRedis *redis.Client

// SetSecurityRedis 设置安全模块使用的Redis客户端
func SetSecurityRedis(client *redis.Client) {
	securityRedis = client
}

// SecurityManager 安全管理器
type SecurityManager struct {
	redis         *redis.Client
	globalLimiter *GlobalRateLimiter
	ipLimiters    map[string]*IPRateLimiter
	userLimiters  map[uint]*UserRateLimiter
	mu            sync.RWMutex
}

var securityManager *SecurityManager
var securityManagerOnce sync.Once
var suspiciousInputPattern = regexp.MustCompile(`(?i)(\b(select|union|insert|update|delete|drop|truncate|alter|exec|sleep|benchmark)\b|--|/\*|\*/|;)`)

// GetSecurityManager 获取安全管理器实例
func GetSecurityManager() *SecurityManager {
	securityManagerOnce.Do(func() {
		securityManager = &SecurityManager{
			redis:         securityRedis,
			globalLimiter: NewGlobalRateLimiter(GlobalRateLimitCapacity, GlobalRateLimitRate, securityRedis),
			ipLimiters:    make(map[string]*IPRateLimiter),
			userLimiters:  make(map[uint]*UserRateLimiter),
		}

		// 启动清理任务
		go securityManager.cleanupTask()
	})
	return securityManager
}

// SecurityCheck 执行完整的安全检查
func (sm *SecurityManager) SecurityCheck(c *gin.Context) (bool, *SecurityError) {
	// 1. 检查黑名单
	if blocked, reason := sm.IsBlocked(c); blocked {
		return false, &SecurityError{
			Code:    http.StatusForbidden,
			Message: "当前操作异常，请稍后再试",
			Reason:  reason,
		}
	}

	// 2. 全局限流检查
	if !sm.globalLimiter.Allow() {
		sm.recordLimitEvent(c, "global")
		return false, &SecurityError{
			Code:    http.StatusTooManyRequests,
			Message: "请求过于频繁，请稍后再试",
			Reason:  "global_rate_limit",
		}
	}

	// 3. IP限流检查
	clientIP := c.ClientIP()
	ipLimiter := sm.getIPLimiter(clientIP)
	if !ipLimiter.Allow(clientIP) {
		sm.recordLimitEvent(c, "ip")
		return false, &SecurityError{
			Code:    http.StatusTooManyRequests,
			Message: "请求过于频繁，请稍后再试",
			Reason:  "ip_rate_limit",
		}
	}

	// 4. 用户限流检查（如果已登录）
	if uid, ok := getContextUserID(c); ok {
		userLimiter := sm.getUserLimiter(uid)
		if !userLimiter.Allow(uid) {
			sm.recordLimitEvent(c, "user")
			return false, &SecurityError{
				Code:    http.StatusTooManyRequests,
				Message: "请求过于频繁，请稍后再试",
				Reason:  "user_rate_limit",
			}
		}
	}

	return true, nil
}

// SecurityError 安全错误
type SecurityError struct {
	Code    int
	Message string
	Reason  string
}

func (se *SecurityError) Error() string {
	return se.Message
}

// IsBlocked 检查是否被封禁
func (sm *SecurityManager) IsBlocked(c *gin.Context) (bool, string) {
	ctx := context.Background()

	// 检查IP黑名单
	clientIP := c.ClientIP()
	ipKey := BlacklistKeyPrefix + "ip:" + clientIP
	if sm.redis != nil {
		exists, _ := sm.redis.Exists(ctx, ipKey).Result()
		if exists > 0 {
			return true, "ip_blacklist"
		}
	}

	// 检查用户黑名单
	if uid, ok := getContextUserID(c); ok {
		userKey := BlacklistKeyPrefix + "user:" + strconv.Itoa(int(uid))
		if sm.redis != nil {
			exists, _ := sm.redis.Exists(ctx, userKey).Result()
			if exists > 0 {
				return true, "user_blacklist"
			}
		}
	}

	return false, ""
}

// AddToBlacklist 添加到黑名单
func (sm *SecurityManager) AddToBlacklist(ctx context.Context, targetType string, targetID string, reason string) error {
	if sm.redis == nil {
		Logger.Warn("Redis未连接，无法添加到黑名单",
			zap.String("targetType", targetType),
			zap.String("targetID", targetID))
		return nil
	}

	key := BlacklistKeyPrefix + targetType + ":" + targetID

	// 获取当前封禁次数
	countKey := key + ":count"
	count, _ := sm.redis.Get(ctx, countKey).Int()
	count++

	// 计算封禁时长（阶梯式）
	expireTime := BlacklistExpireBase + time.Duration(count-1)*BlacklistExpireStep
	if expireTime > BlacklistMaxExpire {
		expireTime = BlacklistMaxExpire
	}

	// 添加到黑名单
	pipe := sm.redis.Pipeline()
	pipe.Set(ctx, key, reason, expireTime)
	pipe.Set(ctx, countKey, count, expireTime)
	_, err := pipe.Exec(ctx)

	if err != nil {
		Logger.Error("添加到黑名单失败",
			zap.String("targetType", targetType),
			zap.String("targetID", targetID),
			zap.Error(err))
		return err
	}

	Logger.Info("已添加到黑名单",
		zap.String("targetType", targetType),
		zap.String("targetID", targetID),
		zap.String("reason", reason),
		zap.Duration("expireTime", expireTime),
		zap.Int("count", count))

	// 记录安全事件
	sm.recordSecurityEvent("blacklist_add", targetType, targetID, reason)

	return nil
}

// RemoveFromBlacklist 从黑名单移除
func (sm *SecurityManager) RemoveFromBlacklist(ctx context.Context, targetType string, targetID string) error {
	if sm.redis == nil {
		return nil
	}

	key := BlacklistKeyPrefix + targetType + ":" + targetID
	countKey := key + ":count"

	pipe := sm.redis.Pipeline()
	pipe.Del(ctx, key)
	pipe.Del(ctx, countKey)
	_, err := pipe.Exec(ctx)

	if err != nil {
		Logger.Error("从黑名单移除失败",
			zap.String("targetType", targetType),
			zap.String("targetID", targetID),
			zap.Error(err))
		return err
	}

	Logger.Info("已从黑名单移除",
		zap.String("targetType", targetType),
		zap.String("targetID", targetID))

	sm.recordSecurityEvent("blacklist_remove", targetType, targetID, "")

	return nil
}

// GetBlacklistInfo 获取黑名单信息
func (sm *SecurityManager) GetBlacklistInfo(ctx context.Context, targetType string, targetID string) (map[string]interface{}, error) {
	if sm.redis == nil {
		return nil, errors.New("Redis未连接")
	}

	key := BlacklistKeyPrefix + targetType + ":" + targetID
	countKey := key + ":count"

	ttl, err := sm.redis.TTL(ctx, key).Result()
	if err != nil || ttl < 0 {
		return nil, errors.New("不在黑名单中")
	}

	reason, _ := sm.redis.Get(ctx, key).Result()
	count, _ := sm.redis.Get(ctx, countKey).Int()

	return map[string]interface{}{
		"target_type": targetType,
		"target_id":   targetID,
		"reason":      reason,
		"ttl_seconds": int(ttl.Seconds()),
		"count":       count,
	}, nil
}

// RecordFailedRequest 记录失败请求
func (sm *SecurityManager) RecordFailedRequest(ctx context.Context, c *gin.Context, failType string) {
	clientIP := c.ClientIP()

	// IP维度统计
	ipFailKey := fmt.Sprintf("security:fail:ip:%s:%s", clientIP, failType)
	ipWindowKey := fmt.Sprintf("security:fail:ip:%s:window", clientIP)

	if sm.redis != nil {
		pipe := sm.redis.Pipeline()
		pipe.Incr(ctx, ipFailKey)
		pipe.Expire(ctx, ipFailKey, 5*time.Minute)
		pipe.Incr(ctx, ipWindowKey)
		pipe.Expire(ctx, ipWindowKey, 1*time.Minute)
		pipe.Exec(ctx)

		// 检查是否需要封禁IP
		failCount, _ := sm.redis.Get(ctx, ipFailKey).Int()
		windowCount, _ := sm.redis.Get(ctx, ipWindowKey).Int()

		if failCount >= IPFailThreshold || windowCount >= IPRepeatThreshold {
			sm.AddToBlacklist(ctx, "ip", clientIP, fmt.Sprintf("异常请求: %s", failType))
		}
	}

	// 用户维度统计
	if uid, ok := getContextUserID(c); ok {
		userFailKey := fmt.Sprintf("security:fail:user:%d:%s", uid, failType)

		if sm.redis != nil {
			count, _ := sm.redis.Incr(ctx, userFailKey).Result()
			sm.redis.Expire(ctx, userFailKey, 10*time.Minute)

			if failType == "captcha" && count >= int64(UserCaptchaThreshold) {
				sm.AddToBlacklist(ctx, "user", strconv.Itoa(int(uid)), "验证码错误次数过多")
			}
			if failType == "order" && count >= int64(UserOrderThreshold) {
				sm.AddToBlacklist(ctx, "user", strconv.Itoa(int(uid)), "订单创建失败次数过多")
			}
		}
	}
}

func getContextUserID(c *gin.Context) (uint, bool) {
	if value, exists := c.Get("user_id"); exists {
		if uid, ok := value.(uint); ok && uid > 0 {
			return uid, true
		}
	}
	if value, exists := c.Get("userID"); exists {
		if uid, ok := value.(uint); ok && uid > 0 {
			return uid, true
		}
	}
	return 0, false
}

// getIPLimiter 获取IP限流器
func (sm *SecurityManager) getIPLimiter(ip string) *IPRateLimiter {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	limiter, exists := sm.ipLimiters[ip]
	if !exists {
		limiter = NewIPRateLimiter(IPRateLimitCapacity, IPRateLimitRate, sm.redis)
		sm.ipLimiters[ip] = limiter
	}

	return limiter
}

// getUserLimiter 获取用户限流器
func (sm *SecurityManager) getUserLimiter(userID uint) *UserRateLimiter {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	limiter, exists := sm.userLimiters[userID]
	if !exists {
		limiter = NewUserRateLimiter(UserRateLimitCapacity, UserRateLimitRate, sm.redis)
		sm.userLimiters[userID] = limiter
	}

	return limiter
}

// recordLimitEvent 记录限流事件
func (sm *SecurityManager) recordLimitEvent(c *gin.Context, limitType string) {
	Logger.Warn("触发限流",
		zap.String("type", limitType),
		zap.String("ip", c.ClientIP()),
		zap.String("path", c.Request.URL.Path))

	sm.recordSecurityEvent("rate_limit", limitType, c.ClientIP(), c.Request.URL.Path)
}

// recordSecurityEvent 记录安全事件
func (sm *SecurityManager) recordSecurityEvent(eventType, targetType, targetID, detail string) {
	Logger.Info("安全事件",
		zap.String("event_type", eventType),
		zap.String("target_type", targetType),
		zap.String("target_id", targetID),
		zap.String("detail", detail),
		zap.Time("timestamp", time.Now()))
}

// cleanupTask 清理任务
func (sm *SecurityManager) cleanupTask() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		sm.mu.Lock()

		// 清理过期的IP限流器
		for ip, limiter := range sm.ipLimiters {
			// 检查是否长时间未使用
			if limiter != nil {
				delete(sm.ipLimiters, ip)
			}
		}

		// 清理过期的用户限流器
		for userID, limiter := range sm.userLimiters {
			if limiter != nil {
				delete(sm.userLimiters, userID)
			}
		}

		sm.mu.Unlock()

		Logger.Info("安全限流器清理完成")
	}
}

// GetSecurityStats 获取安全统计信息
func (sm *SecurityManager) GetSecurityStats(ctx context.Context) map[string]interface{} {
	stats := map[string]interface{}{
		"timestamp": time.Now().Format(time.RFC3339),
	}

	// 获取黑名单数量
	if sm.redis != nil {
		blacklistCount, _ := sm.redis.Keys(ctx, BlacklistKeyPrefix+"*").Result()
		stats["blacklist_count"] = len(blacklistCount) / 2 // 每个黑名单有key和count两个键
	}

	// 限流器数量
	sm.mu.RLock()
	stats["ip_limiter_count"] = len(sm.ipLimiters)
	stats["user_limiter_count"] = len(sm.userLimiters)
	sm.mu.RUnlock()

	return stats
}

// SecurityMiddleware 安全中间件
func SecurityMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/api/captcha" {
			c.Next()
			return
		}

		if hasSuspiciousInput(c) {
			sm := GetSecurityManager()
			sm.RecordFailedRequest(context.Background(), c, "sql_injection")
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "请求参数包含非法内容",
			})
			c.Abort()
			return
		}

		// 获取安全管理器
		sm := GetSecurityManager()

		// 执行安全检查
		allowed, secErr := sm.SecurityCheck(c)
		if !allowed {
			c.JSON(secErr.Code, gin.H{
				"error":  secErr.Message,
				"reason": secErr.Reason,
				"code":   secErr.Code,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func hasSuspiciousInput(c *gin.Context) bool {
	rawQuery := c.Request.URL.RawQuery
	if suspiciousInputPattern.MatchString(rawQuery) {
		return true
	}
	if suspiciousInputPattern.MatchString(c.Request.URL.Path) {
		return true
	}
	for _, values := range c.Request.URL.Query() {
		for _, value := range values {
			if suspiciousInputPattern.MatchString(value) {
				return true
			}
		}
	}
	if err := c.Request.ParseForm(); err == nil {
		for _, values := range c.Request.PostForm {
			for _, value := range values {
				if suspiciousInputPattern.MatchString(strings.TrimSpace(value)) {
					return true
				}
			}
		}
	}
	return false
}

// CaptchaSecurityMiddleware 验证码安全中间件
func CaptchaSecurityMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		sm := GetSecurityManager()
		ba := GetBehaviorAnalyzer()

		// 检查是否需要验证码（根据失败次数和行为分析）
		clientIP := c.ClientIP()
		ctx := context.Background()

		// 获取行为分析的风险等级
		riskLevel := "low"
		if level, exists := c.Get("risk_level"); exists {
			if levelStr, ok := level.(string); ok {
				riskLevel = levelStr
			}
		} else {
			// 如果没有风险等级，进行行为分析
			riskScore := ba.AnalyzeBehavior(c)
			riskLevel = ba.GetRiskLevel(riskScore)
		}

		// 根据风险等级设置验证码难度
		switch riskLevel {
		case "high":
			c.Set("captcha_difficulty", "hard")
		case "medium":
			c.Set("captcha_difficulty", "medium")
		default:
			// 检查失败次数
			if sm.redis != nil {
				failKey := fmt.Sprintf("security:fail:ip:%s:captcha", clientIP)
				failCount, _ := sm.redis.Get(ctx, failKey).Int()

				// 如果失败次数超过阈值，增加验证难度
				if failCount >= UserCaptchaThreshold {
					c.Set("captcha_difficulty", "hard")
				}
			}
		}

		c.Next()
	}
}

// RecordFailedCaptcha 记录验证码失败
func RecordFailedCaptcha(c *gin.Context) {
	sm := GetSecurityManager()
	sm.RecordFailedRequest(context.Background(), c, "captcha")
}

// RecordFailedOrder 记录订单失败
func RecordFailedOrder(c *gin.Context) {
	sm := GetSecurityManager()
	sm.RecordFailedRequest(context.Background(), c, "order")
}

// RecordFailedRequest 记录请求失败
func RecordFailedRequest(c *gin.Context, failType string) {
	sm := GetSecurityManager()
	sm.RecordFailedRequest(context.Background(), c, failType)
}
