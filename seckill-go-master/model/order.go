package model

import (
	"gorm.io/gorm"
	"time"
)

// OrderStatus 订单状态枚举
type OrderStatus int

const (
	OrderStatusPending  OrderStatus = 0 // 待支付
	OrderStatusPaid     OrderStatus = 1 // 已支付
	OrderStatusCanceled OrderStatus = 2 // 已取消
)

// Order 订单表模型
type Order struct {
	ID          uint           `gorm:"primaryKey" json:"id"`                            // 主键ID
	OrderNo     string         `gorm:"size:64;uniqueIndex;not null" json:"order_no"`    // 订单号
	UserID      uint           `gorm:"not null;index" json:"user_id"`                   // 关联用户ID
	ProductID   uint           `gorm:"not null;index" json:"product_id"`                // 关联商品ID
	OrderAmount float64        `gorm:"type:decimal(10,2);not null" json:"order_amount"` // 订单金额
	Status      OrderStatus    `gorm:"not null;default:0" json:"status"`                // 订单状态
	CreatedAt   time.Time      `json:"created_at"`                                      // 创建时间
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`                                  // 软删除标记
}
