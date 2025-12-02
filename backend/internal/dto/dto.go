package dto

// 认证相关
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=64"`
	Password string `json:"password" binding:"required,min=6"`
	Email    string `json:"email" binding:"omitempty,email"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type NewAPILoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// 套餐相关
type CreatePlanRequest struct {
	Name         string  `json:"name" binding:"required"`
	Description  string  `json:"description"`
	PeriodType   string  `json:"period_type" binding:"required,oneof=day week month custom"`
	PeriodDays   int     `json:"period_days" binding:"required,min=1"`
	DailyQuota   int     `json:"daily_quota" binding:"required,min=1"`
	CarryOver    int     `json:"carry_over" binding:"oneof=0 1"`
	MaxCarryOver int     `json:"max_carry_over" binding:"min=0"`
	PriceType    string  `json:"price_type" binding:"required,oneof=fixed daily"`
	Price        float64 `json:"price" binding:"required,min=0"`
	NewAPIGroup  string  `json:"newapi_group" binding:"required"`
	Status       int     `json:"status" binding:"oneof=0 1"`
	SortOrder    int     `json:"sort_order"`
}

type UpdatePlanRequest struct {
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	PeriodType   string  `json:"period_type" binding:"omitempty,oneof=day week month custom"`
	PeriodDays   int     `json:"period_days" binding:"omitempty,min=1"`
	DailyQuota   int     `json:"daily_quota" binding:"omitempty,min=1"`
	CarryOver    int     `json:"carry_over" binding:"omitempty,oneof=0 1"`
	MaxCarryOver int     `json:"max_carry_over" binding:"omitempty,min=0"`
	PriceType    string  `json:"price_type" binding:"omitempty,oneof=fixed daily"`
	Price        float64 `json:"price" binding:"omitempty,min=0"`
	NewAPIGroup  string  `json:"newapi_group"`
	Status       int     `json:"status" binding:"omitempty,oneof=0 1"`
	SortOrder    int     `json:"sort_order"`
}

// 订阅相关
type PurchaseRequest struct {
	PlanID      uint   `json:"plan_id" binding:"required"`
	PeriodDays  int    `json:"period_days"` // 自定义天数（可选）
	NewAPIAction string `json:"newapi_action" binding:"required,oneof=bind_existing create_new overwrite"`
	// bind_existing 时需要
	NewAPIUsername string `json:"newapi_username"`
	NewAPIPassword string `json:"newapi_password"`
}

type RenewRequest struct {
	PeriodDays int `json:"period_days" binding:"required,min=1"`
}

// 支付相关
type PayRequest struct {
	OrderID       uint   `json:"order_id" binding:"required"`
	PaymentMethod string `json:"payment_method" binding:"required,oneof=alipay wxpay"`
}

// 用户相关
type UpdateProfileRequest struct {
	Email string `json:"email" binding:"omitempty,email"`
}

type BindNewAPIRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UpdateEmailSettingsRequest struct {
	EmailRemind int `json:"email_remind" binding:"oneof=0 1"`
	RemindDays  int `json:"remind_days" binding:"min=1,max=30"`
}

// 通用响应
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type PaginatedResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Total   int64       `json:"total"`
	Page    int         `json:"page"`
	PerPage int         `json:"per_page"`
}

// 分页参数
type PaginationQuery struct {
	Page    int `form:"page,default=1" binding:"min=1"`
	PerPage int `form:"per_page,default=20" binding:"min=1,max=100"`
}

func (p *PaginationQuery) Offset() int {
	return (p.Page - 1) * p.PerPage
}
