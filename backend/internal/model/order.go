package model

import (
	"time"

	"gorm.io/gorm"
)

// Order 订单模型
type Order struct {
	ID     uint `gorm:"primaryKey" json:"id"`
	OrderNo string `gorm:"uniqueIndex;size:64;not null" json:"order_no"`
	UserID  uint   `gorm:"not null;index" json:"user_id"`
	PlanID  uint   `gorm:"not null" json:"plan_id"`

	// 订单信息
	OrderType  string  `gorm:"size:16;not null" json:"order_type"` // new=新购, renew=续费
	PeriodDays int     `gorm:"not null" json:"period_days"`
	Amount     float64 `gorm:"type:decimal(10,2);not null" json:"amount"`

	// 支付信息
	PaymentMethod string `gorm:"size:32" json:"payment_method"` // alipay/wxpay
	TradeNo       string `gorm:"size:128" json:"trade_no"`

	// 状态
	Status string     `gorm:"size:16;not null" json:"status"` // pending/paid/cancelled/refunded
	PaidAt *time.Time `json:"paid_at"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Plan *Plan `gorm:"foreignKey:PlanID" json:"plan,omitempty"`
}

const (
	OrderTypeNew   = "new"
	OrderTypeRenew = "renew"

	OrderStatusPending   = "pending"
	OrderStatusPaid      = "paid"
	OrderStatusCancelled = "cancelled"
	OrderStatusRefunded  = "refunded"
)
