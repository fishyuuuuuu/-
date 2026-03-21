package model

import (
	"gorm.io/gorm"
)

// Product 商品表模型
type Product struct {
	ID    uint64  `gorm:"primaryKey" json:"id"`                     // 主键ID
	Name  string  `gorm:"size:100;not null" json:"name"`            // 商品名称
	Price float64 `gorm:"type:decimal(10,2);not null" json:"price"` // 商品价格
	Stock int     `gorm:"not null;default:0" json:"stock"`          // 总库存
	//Version   int            `gorm:"not null;default:0" json:"version"`        // 版本号(用于乐观锁)
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // 软删除标记
}
