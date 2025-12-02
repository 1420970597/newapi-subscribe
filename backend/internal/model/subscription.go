package model

import (
	"time"

	"gorm.io/gorm"
)

// Subscription 用户订阅模型
type Subscription struct {
	ID     uint `gorm:"primaryKey" json:"id"`
	UserID uint `gorm:"not null;index" json:"user_id"`
	PlanID uint `gorm:"not null" json:"plan_id"`

	// 订阅状态
	Status string `gorm:"size:16;not null" json:"status"` // active/expired/cancelled

	// 时间信息
	StartDate time.Time `gorm:"type:date;not null" json:"start_date"`
	EndDate   time.Time `gorm:"type:date;not null" json:"end_date"`

	// 当日额度信息
	TodayQuota    int        `gorm:"not null" json:"today_quota"`
	CarriedQuota  int        `gorm:"default:0" json:"carried_quota"`
	LastSyncDate  *time.Time `gorm:"type:date" json:"last_sync_date"`

	// 配置快照（购买时的套餐配置）
	DailyQuota   int    `gorm:"not null" json:"daily_quota"`
	CarryOver    int    `gorm:"not null" json:"carry_over"`
	MaxCarryOver int    `gorm:"default:0" json:"max_carry_over"`
	NewAPIGroup  string `gorm:"column:newapi_group;size:64;not null" json:"newapi_group"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Plan *Plan `gorm:"foreignKey:PlanID" json:"plan,omitempty"`
}

const (
	SubscriptionStatusActive    = "active"
	SubscriptionStatusExpired   = "expired"
	SubscriptionStatusCancelled = "cancelled"
)

// IsActive 是否有效
func (s *Subscription) IsActive() bool {
	return s.Status == SubscriptionStatusActive && time.Now().Before(s.EndDate.AddDate(0, 0, 1))
}

// DaysRemaining 剩余天数
func (s *Subscription) DaysRemaining() int {
	if !s.IsActive() {
		return 0
	}
	days := int(time.Until(s.EndDate).Hours()/24) + 1
	if days < 0 {
		return 0
	}
	return days
}
