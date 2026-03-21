package model

import (
	"gorm.io/gorm"
	"time"
)

// Event 秒杀活动表模型
type Event struct {
	ID        uint           `gorm:"primaryKey" json:"id"`             // 主键ID
	ProductID uint           `gorm:"not null;index" json:"product_id"` // 关联商品ID
	StartTime time.Time      `gorm:"not null" json:"start_time"`       // 活动开始时间
	EndTime   time.Time      `gorm:"not null" json:"end_time"`         // 活动结束时间
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`                   // 软删除标记
}
