package database

import (
	"time"
)

// ModelInsertable 预定义带有自动添加创建时间戳模型
type ModelInsertable struct {
	Model
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}
