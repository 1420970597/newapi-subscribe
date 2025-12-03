package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"newapi-subscribe/internal/dto"
	"newapi-subscribe/internal/middleware"
	"newapi-subscribe/internal/model"
	"newapi-subscribe/internal/service"
)

// GetCurrentSubscription 获取当前订阅状态
func GetCurrentSubscription(c *gin.Context) {
	user := middleware.GetCurrentUser(c)

	var subscription model.Subscription
	err := model.DB.Preload("Plan").
		Where("user_id = ? AND status = ?", user.ID, model.SubscriptionStatusActive).
		First(&subscription).Error

	if err != nil {
		c.JSON(http.StatusOK, dto.Response{
			Success: true,
			Data:    nil, // 没有订阅
		})
		return
	}

	// 获取 new-api 当前余额
	var currentQuota int
	if user.NewAPIBound == 1 {
		client := service.NewNewAPIClient()
		if newAPIUser, err := client.GetUser(user.NewAPIUserID); err == nil {
			currentQuota = newAPIUser.Quota
		}
	}

	c.JSON(http.StatusOK, dto.Response{
		Success: true,
		Data: gin.H{
			"subscription":  subscription,
			"current_quota": currentQuota,
			"days_remaining": subscription.DaysRemaining(),
		},
	})
}

// PurchaseSubscription 购买订阅
func PurchaseSubscription(c *gin.Context) {
	var req dto.PurchaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	user := middleware.GetCurrentUser(c)

	// 获取套餐
	var plan model.Plan
	if err := model.DB.First(&plan, req.PlanID).Error; err != nil {
		c.JSON(http.StatusNotFound, dto.Response{
			Success: false,
			Message: "套餐不存在",
		})
		return
	}

	if plan.Status != model.PlanStatusOn {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "套餐已下架",
		})
		return
	}

	// 确定购买天数
	periodDays := plan.PeriodDays
	if plan.PeriodType == model.PeriodTypeCustom && req.PeriodDays > 0 {
		periodDays = req.PeriodDays
	}

	// 检查是否已有相同套餐的活跃订阅
	var existingSub model.Subscription
	if err := model.DB.Where("user_id = ? AND plan_id = ? AND status = ?",
		user.ID, plan.ID, model.SubscriptionStatusActive).First(&existingSub).Error; err == nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "您已有该套餐的有效订阅，请选择续费或其他套餐",
		})
		return
	}

	// 处理 new-api 账号
	client := service.NewNewAPIClient()
	switch req.NewAPIAction {
	case "bind_existing":
		// 验证并绑定现有账号
		newAPIUser, err := client.Login(req.NewAPIUsername, req.NewAPIPassword)
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.Response{
				Success: false,
				Message: "new-api 账号验证失败",
			})
			return
		}
		user.NewAPIUserID = newAPIUser.ID
		user.NewAPIUsername = newAPIUser.Username
		user.NewAPIBound = 1
		model.DB.Save(user)

	case "create_new":
		// 检查是否已绑定
		if user.NewAPIBound == 1 {
			c.JSON(http.StatusBadRequest, dto.Response{
				Success: false,
				Message: "您已绑定 new-api 账号，请选择其他操作",
			})
			return
		}
		// 将在支付成功后创建

	case "overwrite":
		if user.NewAPIBound != 1 {
			c.JSON(http.StatusBadRequest, dto.Response{
				Success: false,
				Message: "您未绑定 new-api 账号",
			})
			return
		}
		// 将在支付成功后覆盖
	}

	// 计算价格
	amount := plan.CalculatePrice(periodDays)

	// 创建订单
	order := &model.Order{
		OrderNo:    generateOrderNo(user.ID),
		UserID:     user.ID,
		PlanID:     plan.ID,
		OrderType:  model.OrderTypeNew,
		PeriodDays: periodDays,
		Amount:     amount,
		Status:     model.OrderStatusPending,
	}

	// 检查是否有其他活跃订阅（如果有则为续费/切换套餐）
	var anyActiveSub model.Subscription
	if err := model.DB.Where("user_id = ? AND status = ?", user.ID, model.SubscriptionStatusActive).
		First(&anyActiveSub).Error; err == nil {
		order.OrderType = model.OrderTypeRenew
	}

	if err := model.DB.Create(order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Message: "创建订单失败",
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Success: true,
		Data: gin.H{
			"order":        order,
			"newapi_action": req.NewAPIAction,
		},
	})
}

// RenewSubscription 续费订阅
func RenewSubscription(c *gin.Context) {
	var req dto.RenewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "参数错误",
		})
		return
	}

	user := middleware.GetCurrentUser(c)

	// 获取当前订阅
	var subscription model.Subscription
	if err := model.DB.Preload("Plan").
		Where("user_id = ? AND status = ?", user.ID, model.SubscriptionStatusActive).
		First(&subscription).Error; err != nil {
		c.JSON(http.StatusNotFound, dto.Response{
			Success: false,
			Message: "没有可续费的订阅",
		})
		return
	}

	// 计算价格
	var plan model.Plan
	model.DB.First(&plan, subscription.PlanID)
	amount := plan.CalculatePrice(req.PeriodDays)

	// 创建续费订单
	order := &model.Order{
		OrderNo:    generateOrderNo(user.ID),
		UserID:     user.ID,
		PlanID:     plan.ID,
		OrderType:  model.OrderTypeRenew,
		PeriodDays: req.PeriodDays,
		Amount:     amount,
		Status:     model.OrderStatusPending,
	}

	if err := model.DB.Create(order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Message: "创建订单失败",
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Success: true,
		Data:    order,
	})
}

// GetUsageLogs 获取使用日志
func GetUsageLogs(c *gin.Context) {
	user := middleware.GetCurrentUser(c)

	var pagination dto.PaginationQuery
	if err := c.ShouldBindQuery(&pagination); err != nil {
		pagination.Page = 1
		pagination.PerPage = 20
	}

	// 获取用户订阅
	var subscription model.Subscription
	if err := model.DB.Where("user_id = ?", user.ID).Order("created_at DESC").First(&subscription).Error; err != nil {
		c.JSON(http.StatusOK, dto.PaginatedResponse{
			Success: true,
			Data:    []model.UsageLog{},
			Total:   0,
			Page:    pagination.Page,
			PerPage: pagination.PerPage,
		})
		return
	}

	var logs []model.UsageLog
	var total int64

	model.DB.Model(&model.UsageLog{}).Where("subscription_id = ?", subscription.ID).Count(&total)
	model.DB.Where("subscription_id = ?", subscription.ID).
		Order("log_date DESC").
		Offset(pagination.Offset()).
		Limit(pagination.PerPage).
		Find(&logs)

	c.JSON(http.StatusOK, dto.PaginatedResponse{
		Success: true,
		Data:    logs,
		Total:   total,
		Page:    pagination.Page,
		PerPage: pagination.PerPage,
	})
}

// GetUsageDetail 获取详细使用日志（从 new-api）
func GetUsageDetail(c *gin.Context) {
	user := middleware.GetCurrentUser(c)

	if user.NewAPIBound != 1 {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "请先绑定 new-api 账号，或前往 new-api 站点查询",
		})
		return
	}

	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	client := service.NewNewAPIClient()
	logs, err := client.GetUserLogs(user.NewAPIUserID, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Message: "获取日志失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Success: true,
		Data:    logs,
	})
}

// GetTodayUsage 获取当日用量
func GetTodayUsage(c *gin.Context) {
	user := middleware.GetCurrentUser(c)

	if user.NewAPIBound != 1 {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "请先绑定 new-api 账号",
		})
		return
	}

	client := service.NewNewAPIClient()
	todayUsed, err := client.GetUserQuotaUsedToday(user.NewAPIUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Message: "获取今日用量失败: " + err.Error(),
		})
		return
	}

	// 获取当前订阅信息
	var subscription model.Subscription
	var dailyQuota int
	if err := model.DB.Where("user_id = ? AND status = ?", user.ID, model.SubscriptionStatusActive).
		First(&subscription).Error; err == nil {
		dailyQuota = subscription.DailyQuota
	}

	// 获取当前余额
	var currentQuota int
	if newAPIUser, err := client.GetUser(user.NewAPIUserID); err == nil {
		currentQuota = newAPIUser.Quota
	}

	c.JSON(http.StatusOK, dto.Response{
		Success: true,
		Data: gin.H{
			"today_used":    todayUsed,
			"daily_quota":   dailyQuota,
			"current_quota": currentQuota,
		},
	})
}

func generateOrderNo(userID uint) string {
	return fmt.Sprintf("SUB%d%d", userID, time.Now().UnixNano())
}
