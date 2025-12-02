package model

import (
	"time"

	"gorm.io/gorm"
)

// Plan 订阅套餐模型
type Plan struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:128;not null" json:"name"`
	Description string         `gorm:"type:text" json:"description"`

	// 周期设置
	PeriodType string `gorm:"size:16;not null" json:"period_type"` // day/week/month/custom
	PeriodDays int    `gorm:"not null" json:"period_days"`

	// 额度设置
	DailyQuota   int `gorm:"not null" json:"daily_quota"`
	CarryOver    int `gorm:"default:0" json:"carry_over"`     // 0=不结转, 1=结转
	MaxCarryOver int `gorm:"default:0" json:"max_carry_over"` // 最大结转额度 (0=无限制)

	// 价格设置
	PriceType string  `gorm:"size:16;not null" json:"price_type"` // fixed=固定价格, daily=按天计价
	Price     float64 `gorm:"type:decimal(10,2);not null" json:"price"`

	// new-api 分组绑定
	NewAPIGroup string `gorm:"column:newapi_group;size:64;not null" json:"newapi_group"`

	// 状态
	Status    int `gorm:"default:1" json:"status"` // 1=上架, 0=下架
	SortOrder int `gorm:"default:0" json:"sort_order"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

const (
	PeriodTypeDay    = "day"
	PeriodTypeWeek   = "week"
	PeriodTypeMonth  = "month"
	PeriodTypeCustom = "custom"

	PriceTypeFixed = "fixed"
	PriceTypeDaily = "daily"

	PlanStatusOn  = 1
	PlanStatusOff = 0
)

// CalculatePrice 计算订单价格
func (p *Plan) CalculatePrice(days int) float64 {
	if p.PriceType == PriceTypeFixed {
		return p.Price
	}
	// 按天计价
	return p.Price * float64(days)
}
