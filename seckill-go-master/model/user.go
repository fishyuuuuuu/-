package model

import (
	"time"

	"gorm.io/gorm"
)

// User 用户表模型
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`                         // 主键ID
	Username          string         `gorm:"size:50;uniqueIndex;not null" json:"username"` // 用户名(唯一)
	Password          string         `gorm:"size:100;not null" json:"password"`            // 密码(密文)
	Phone             string         `gorm:"size:20;uniqueIndex" json:"phone"`             // 手机号(唯一)
	Status            string         `gorm:"type:varchar(20);default:'正常'" json:"status"`  // 状态: 正常/已拉黑
	BlacklistReason   string         `gorm:"type:varchar(255)" json:"blacklist_reason"`    // 拉黑原因
	BlacklistTime     *time.Time     `json:"blacklist_time"`                               // 拉黑时间
	BlacklistOperator string         `gorm:"type:varchar(50)" json:"blacklist_operator"`   // 拉黑操作人
	CreatedAt         time.Time      `json:"created_at"`                                   // 创建时间
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`                               // 软删除标记
	Roles             []Role         `gorm:"many2many:user_roles;" json:"-"`               // 用户角色关联
}
