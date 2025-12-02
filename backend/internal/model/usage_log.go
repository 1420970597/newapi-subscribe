package model

import (
	"time"
)

// UsageLog 使用日志模型（本地缓存）
type UsageLog struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	SubscriptionID uint      `gorm:"not null;index" json:"subscription_id"`
	LogDate        time.Time `gorm:"type:date;not null;index" json:"log_date"`

	// 汇总数据
	TotalQuota   int `gorm:"default:0" json:"total_quota"`
	RequestCount int `gorm:"default:0" json:"request_count"`

	// 模型分布 (JSON)
	ModelUsage string `gorm:"type:text" json:"model_usage"` // {"gpt-4": 1000, "gpt-3.5": 500}

	CreatedAt time.Time `json:"created_at"`
}
