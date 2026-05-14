package utils

import (
	"go.uber.org/zap"
)

// AuditLog 记录审计日志
func AuditLog(action string, userID string, details map[string]interface{}) {
	// 构建日志字段
	fields := []zap.Field{
		zap.String("action", action),
		zap.String("user_id", userID),
	}

	// 添加详细信息
	for key, value := range details {
		fields = append(fields, zap.Any(key, value))
	}

	// 记录审计日志
	Logger.Info("审计日志", fields...)
}

// AuditLogin 记录登录审计日志
func AuditLogin(userID string, success bool, ip string) {
	AuditLog("login", userID, map[string]interface{}{
		"success": success,
		"ip":      ip,
	})
}

// AuditRegister 记录注册审计日志
func AuditRegister(userID string, success bool, ip string) {
	AuditLog("register", userID, map[string]interface{}{
		"success": success,
		"ip":      ip,
	})
}

// AuditSeckill 记录秒杀审计日志
func AuditSeckill(userID string, productID string, success bool, ip string) {
	AuditLog("seckill", userID, map[string]interface{}{
		"product_id": productID,
		"success":    success,
		"ip":         ip,
	})
}

// AuditProductManagement 记录商品管理审计日志
func AuditProductManagement(userID string, action string, productID string, ip string) {
	AuditLog("product_management", userID, map[string]interface{}{
		"sub_action": action,
		"product_id": productID,
		"ip":         ip,
	})
}

// AuditEventManagement 记录活动管理审计日志
func AuditEventManagement(userID string, action string, eventID string, ip string) {
	AuditLog("event_management", userID, map[string]interface{}{
		"sub_action": action,
		"event_id":   eventID,
		"ip":         ip,
	})
}

// AuditBlacklistManagement 记录账号黑名单管理审计日志
func AuditBlacklistManagement(operator string, action string, targetUserID uint, targetUsername string, reason string, ip string) {
	AuditLog("user_blacklist_management", operator, map[string]interface{}{
		"sub_action":      action,
		"target_user_id":  targetUserID,
		"target_username": targetUsername,
		"reason":          reason,
		"ip":              ip,
	})
}
