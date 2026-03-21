package model

import (
	"time"

	"gorm.io/gorm"
)

// User 用户表模型
type User struct {
	ID        uint64         `gorm:"primaryKey" json:"id"`                         // 主键ID
	Username  string         `gorm:"size:50;uniqueIndex;not null" json:"username"` // 用户名(唯一)
	Password  string         `gorm:"size:100;not null" json:"-"`                   // 密码(密文，返回JSON时隐藏)
	Phone     string         `gorm:"size:20;uniqueIndex" json:"phone"`             // 手机号(唯一)
	CreatedAt time.Time      `json:"created_at"`                                   // 创建时间
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`                               // 软删除标记
}
