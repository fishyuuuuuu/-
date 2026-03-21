package utils

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

// RateLimiter 令牌桶速率限制器
type RateLimiter struct {
	capacity       int        // 令牌桶容量
	rate           int        // 每秒生成令牌数
	tokens         int        // 当前令牌数
	lastRefillTime time.Time  // 上次填充时间
	mu             sync.Mutex // 互斥锁
}

// NewRateLimiter 创建一个新的令牌桶速率限制器
func NewRateLimiter(capacity, rate int) *RateLimiter {
	return &RateLimiter{
		capacity:       capacity,
		rate:           rate,
		tokens:         capacity,
		lastRefillTime: time.Now(),
	}
}

// Allow 检查是否允许请求
func (rl *RateLimiter) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	// 计算自上次填充以来的时间
	now := time.Now()
	timeElapsed := now.Sub(rl.lastRefillTime).Seconds()

	// 填充令牌
	tokensToAdd := int(timeElapsed * float64(rl.rate))
	if tokensToAdd > 0 {
		rl.tokens = min(rl.capacity, rl.tokens+tokensToAdd)
		rl.lastRefillTime = now
	}

	// 检查是否有足够的令牌
	if rl.tokens > 0 {
		rl.tokens--
		return true
	}

	return false
}

// min 返回两个整数中的最小值
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// IPRateLimiter IP速率限制器
type IPRateLimiter struct {
	limiters    map[string]*RateLimiter
	mu          sync.Mutex
	capacity    int
	rate        int
	redisClient *redis.Client
}

// NewIPRateLimiter 创建一个新的IP速率限制器
func NewIPRateLimiter(capacity, rate int, redisClient *redis.Client) *IPRateLimiter {
	return &IPRateLimiter{
		limiters:    make(map[string]*RateLimiter),
		capacity:    capacity,
		rate:        rate,
		redisClient: redisClient,
	}
}

// Allow 检查指定IP是否允许请求
func (iprl *IPRateLimiter) Allow(ip string) bool {
	// 首先尝试使用Redis进行分布式限流
	if iprl.redisClient != nil {
		if allowed, err := iprl.allowWithRedis(ip); err == nil {
			return allowed
		}
	}

	// 如果Redis不可用，使用本地限流
	iprl.mu.Lock()
	limiter, exists := iprl.limiters[ip]
	if !exists {
		limiter = NewRateLimiter(iprl.capacity, iprl.rate)
		iprl.limiters[ip] = limiter
	}
	iprl.mu.Unlock()

	return limiter.Allow()
}

// allowWithRedis 使用Redis进行分布式限流
func (iprl *IPRateLimiter) allowWithRedis(ip string) (bool, error) {
	ctx := context.Background()
	key := fmt.Sprintf("rate_limit:ip:%s", ip)

	// 使用Redis Lua脚本实现令牌桶算法
	luaScript := `
	local key = KEYS[1]
	local capacity = tonumber(ARGV[1])
	local rate = tonumber(ARGV[2])
	local now = tonumber(ARGV[3])
	local interval = 1.0

	local current = tonumber(redis.call('get', key)) or capacity
	local lastRefillTime = tonumber(redis.call('get', key .. ':last')) or now

	local timeElapsed = now - lastRefillTime
	local tokensToAdd = math.floor(timeElapsed * rate / 1000)

	if tokensToAdd > 0 then
		current = math.min(capacity, current + tokensToAdd)
		redis.call('set', key .. ':last', now)
	end

	if current > 0 then
		current = current - 1
		redis.call('set', key, current)
		redis.call('expire', key, 600) -- 10分钟过期
		redis.call('expire', key .. ':last', 600)
		return 1
	else
		return 0
	end
	`

	result, err := iprl.redisClient.Eval(ctx, luaScript, []string{key}, iprl.capacity, iprl.rate, time.Now().UnixMilli()).Result()
	if err != nil {
		return false, err
	}

	return result == int64(1), nil
}

// UserRateLimiter 用户速率限制器
type UserRateLimiter struct {
	limiters    map[uint]*RateLimiter
	mu          sync.Mutex
	capacity    int
	rate        int
	redisClient *redis.Client
}

// NewUserRateLimiter 创建一个新的用户速率限制器
func NewUserRateLimiter(capacity, rate int, redisClient *redis.Client) *UserRateLimiter {
	return &UserRateLimiter{
		limiters:    make(map[uint]*RateLimiter),
		capacity:    capacity,
		rate:        rate,
		redisClient: redisClient,
	}
}

// Allow 检查指定用户是否允许请求
func (url *UserRateLimiter) Allow(userID uint) bool {
	// 首先尝试使用Redis进行分布式限流
	if url.redisClient != nil {
		if allowed, err := url.allowWithRedis(userID); err == nil {
			return allowed
		}
	}

	// 如果Redis不可用，使用本地限流
	url.mu.Lock()
	limiter, exists := url.limiters[userID]
	if !exists {
		limiter = NewRateLimiter(url.capacity, url.rate)
		url.limiters[userID] = limiter
	}
	url.mu.Unlock()

	return limiter.Allow()
}

// allowWithRedis 使用Redis进行分布式限流
func (url *UserRateLimiter) allowWithRedis(userID uint) (bool, error) {
	ctx := context.Background()
	key := fmt.Sprintf("rate_limit:user:%d", userID)

	// 使用Redis Lua脚本实现令牌桶算法
	luaScript := `
	local key = KEYS[1]
	local capacity = tonumber(ARGV[1])
	local rate = tonumber(ARGV[2])
	local now = tonumber(ARGV[3])
	local interval = 1.0

	local current = tonumber(redis.call('get', key)) or capacity
	local lastRefillTime = tonumber(redis.call('get', key .. ':last')) or now

	local timeElapsed = now - lastRefillTime
	local tokensToAdd = math.floor(timeElapsed * rate / 1000)

	if tokensToAdd > 0 then
		current = math.min(capacity, current + tokensToAdd)
		redis.call('set', key .. ':last', now)
	end

	if current > 0 then
		current = current - 1
		redis.call('set', key, current)
		redis.call('expire', key, 600) -- 10分钟过期
		redis.call('expire', key .. ':last', 600)
		return 1
	else
		return 0
	end
	`

	result, err := url.redisClient.Eval(ctx, luaScript, []string{key}, url.capacity, url.rate, time.Now().UnixMilli()).Result()
	if err != nil {
		return false, err
	}

	return result == int64(1), nil
}

// Cleanup 清理过期的限流器
func (iprl *IPRateLimiter) Cleanup() {
	iprl.mu.Lock()
	defer iprl.mu.Unlock()

	for ip, limiter := range iprl.limiters {
		// 检查限流器是否长时间未使用
		if time.Since(limiter.lastRefillTime) > 5*time.Minute {
			delete(iprl.limiters, ip)
		}
	}
}

// GlobalRateLimiter 全局速率限制器
type GlobalRateLimiter struct {
	limiter     *RateLimiter
	redisClient *redis.Client
}

// NewGlobalRateLimiter 创建一个新的全局速率限制器
func NewGlobalRateLimiter(capacity, rate int, redisClient *redis.Client) *GlobalRateLimiter {
	return &GlobalRateLimiter{
		limiter:     NewRateLimiter(capacity, rate),
		redisClient: redisClient,
	}
}

// Allow 检查是否允许请求
func (grl *GlobalRateLimiter) Allow() bool {
	// 首先尝试使用Redis进行分布式限流
	if grl.redisClient != nil {
		if allowed, err := grl.allowWithRedis(); err == nil {
			return allowed
		}
	}

	// 如果Redis不可用，使用本地限流
	return grl.limiter.Allow()
}

// allowWithRedis 使用Redis进行分布式限流
func (grl *GlobalRateLimiter) allowWithRedis() (bool, error) {
	ctx := context.Background()
	key := "rate_limit:global"

	// 使用Redis Lua脚本实现令牌桶算法
	luaScript := `
	local key = KEYS[1]
	local capacity = tonumber(ARGV[1])
	local rate = tonumber(ARGV[2])
	local now = tonumber(ARGV[3])
	local interval = 1.0

	local current = tonumber(redis.call('get', key)) or capacity
	local lastRefillTime = tonumber(redis.call('get', key .. ':last')) or now

	local timeElapsed = now - lastRefillTime
	local tokensToAdd = math.floor(timeElapsed * rate / 1000)

	if tokensToAdd > 0 then
		current = math.min(capacity, current + tokensToAdd)
		redis.call('set', key .. ':last', now)
	end

	if current > 0 then
		current = current - 1
		redis.call('set', key, current)
		redis.call('expire', key, 600) -- 10分钟过期
		redis.call('expire', key .. ':last', 600)
		return 1
	else
		return 0
	end
	`

	result, err := grl.redisClient.Eval(ctx, luaScript, []string{key}, 1000, 100, time.Now().UnixMilli()).Result()
	if err != nil {
		return false, err
	}

	return result == int64(1), nil
}

// 全局限流器实例
var (
	ipRateLimiter     *IPRateLimiter
	userRateLimiter   *UserRateLimiter
	globalRateLimiter *GlobalRateLimiter
	ipLimiterOnce     sync.Once
	userLimiterOnce   sync.Once
	globalLimiterOnce sync.Once
	globalRedisClient *redis.Client
)

// SetRedisClient 设置全局Redis客户端
func SetRedisClient(client *redis.Client) {
	globalRedisClient = client
}

// GetIPRateLimiter 获取IP速率限制器实例
func GetIPRateLimiter() *IPRateLimiter {
	ipLimiterOnce.Do(func() {
		// 提高限流容量，支持高并发测试
		capacity := 10000
		rate := 5000

		ipRateLimiter = NewIPRateLimiter(capacity, rate, globalRedisClient)

		go func() {
			ticker := time.NewTicker(1 * time.Minute)
			defer ticker.Stop()
			for range ticker.C {
				ipRateLimiter.Cleanup()
			}
		}()
	})
	return ipRateLimiter
}

// GetUserRateLimiter 获取用户速率限制器实例
func GetUserRateLimiter() *UserRateLimiter {
	userLimiterOnce.Do(func() {
		capacity := 10
		rate := 5

		userRateLimiter = NewUserRateLimiter(capacity, rate, globalRedisClient)
	})
	return userRateLimiter
}

// GetGlobalRateLimiter 获取全局速率限制器实例
func GetGlobalRateLimiter() *GlobalRateLimiter {
	globalLimiterOnce.Do(func() {
		capacity := 100
		rate := 50

		globalRateLimiter = NewGlobalRateLimiter(capacity, rate, globalRedisClient)
	})
	return globalRateLimiter
}
