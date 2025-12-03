package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"newapi-subscribe/internal/dto"
	"newapi-subscribe/internal/model"
	"newapi-subscribe/internal/service"
)

// AdminGetUsers 获取用户列表
func AdminGetUsers(c *gin.Context) {
	var pagination dto.PaginationQuery
	if err := c.ShouldBindQuery(&pagination); err != nil {
		pagination.Page = 1
		pagination.PerPage = 20
	}

	keyword := c.Query("keyword")

	var users []model.User
	var total int64

	query := model.DB.Model(&model.User{})
	if keyword != "" {
		query = query.Where("username LIKE ? OR email LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	query.Count(&total)
	query.Order("id DESC").Offset(pagination.Offset()).Limit(pagination.PerPage).Find(&users)

	// 获取每个用户的订阅状态
	type UserWithSubscription struct {
		model.User
		Subscription *model.Subscription `json:"subscription"`
	}

	result := make([]UserWithSubscription, len(users))
	for i, u := range users {
		result[i].User = u
		var sub model.Subscription
		if err := model.DB.Preload("Plan").
			Where("user_id = ? AND status = ?", u.ID, model.SubscriptionStatusActive).
			First(&sub).Error; err == nil {
			result[i].Subscription = &sub
		}
	}

	c.JSON(http.StatusOK, dto.PaginatedResponse{
		Success: true,
		Data:    result,
		Total:   total,
		Page:    pagination.Page,
		PerPage: pagination.PerPage,
	})
}

// AdminGetUser 获取用户详情
func AdminGetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "无效的用户 ID",
		})
		return
	}

	var user model.User
	if err := model.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, dto.Response{
			Success: false,
			Message: "用户不存在",
		})
		return
	}

	// 获取订阅
	var subscription model.Subscription
	model.DB.Preload("Plan").
		Where("user_id = ? AND status = ?", user.ID, model.SubscriptionStatusActive).
		First(&subscription)

	// 获取 new-api 余额
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
			"user":          user,
			"subscription":  subscription,
			"current_quota": currentQuota,
		},
	})
}

// AdminGetUserUsage 获取用户使用分析
func AdminGetUserUsage(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "无效的用户 ID",
		})
		return
	}

	var user model.User
	if err := model.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, dto.Response{
			Success: false,
			Message: "用户不存在",
		})
		return
	}

	if user.NewAPIBound != 1 {
		c.JSON(http.StatusOK, dto.Response{
			Success: true,
			Data:    nil,
		})
		return
	}

	startDate := c.DefaultQuery("start_date", "")
	endDate := c.DefaultQuery("end_date", "")

	client := service.NewNewAPIClient()
	logs, err := client.GetUserLogs(user.NewAPIUserID, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Message: "获取使用记录失败",
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Success: true,
		Data:    logs,
	})
}

// AdminUpdateUser 更新用户信息
func AdminUpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "无效的用户 ID",
		})
		return
	}

	var user model.User
	if err := model.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, dto.Response{
			Success: false,
			Message: "用户不存在",
		})
		return
	}

	var req struct {
		Email       string `json:"email"`
		Status      int    `json:"status"`
		Role        int    `json:"role"`
		EmailRemind int    `json:"email_remind"`
		RemindDays  int    `json:"remind_days"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "参数错误",
		})
		return
	}

	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Status > 0 {
		user.Status = req.Status
	}
	if req.Role > 0 {
		user.Role = req.Role
	}
	user.EmailRemind = req.EmailRemind
	if req.RemindDays > 0 {
		user.RemindDays = req.RemindDays
	}

	if err := model.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Message: "更新失败",
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Success: true,
		Data:    user,
	})
}

// AdminGetSubscriptions 获取所有订阅
func AdminGetSubscriptions(c *gin.Context) {
	var pagination dto.PaginationQuery
	if err := c.ShouldBindQuery(&pagination); err != nil {
		pagination.Page = 1
		pagination.PerPage = 20
	}

	status := c.Query("status")

	var subscriptions []model.Subscription
	var total int64

	query := model.DB.Model(&model.Subscription{})
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)
	query.Preload("User").Preload("Plan").
		Order("id DESC").
		Offset(pagination.Offset()).
		Limit(pagination.PerPage).
		Find(&subscriptions)

	c.JSON(http.StatusOK, dto.PaginatedResponse{
		Success: true,
		Data:    subscriptions,
		Total:   total,
		Page:    pagination.Page,
		PerPage: pagination.PerPage,
	})
}

// AdminGetOrders 获取所有订单
func AdminGetOrders(c *gin.Context) {
	var pagination dto.PaginationQuery
	if err := c.ShouldBindQuery(&pagination); err != nil {
		pagination.Page = 1
		pagination.PerPage = 20
	}

	status := c.Query("status")

	var orders []model.Order
	var total int64

	query := model.DB.Model(&model.Order{})
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)
	query.Preload("User").Preload("Plan").
		Order("id DESC").
		Offset(pagination.Offset()).
		Limit(pagination.PerPage).
		Find(&orders)

	c.JSON(http.StatusOK, dto.PaginatedResponse{
		Success: true,
		Data:    orders,
		Total:   total,
		Page:    pagination.Page,
		PerPage: pagination.PerPage,
	})
}

// AdminCreatePlan 创建套餐
func AdminCreatePlan(c *gin.Context) {
	var req dto.CreatePlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	plan := &model.Plan{
		Name:         req.Name,
		Description:  req.Description,
		PeriodType:   req.PeriodType,
		PeriodDays:   req.PeriodDays,
		DailyQuota:   req.DailyQuota,
		CarryOver:    req.CarryOver,
		MaxCarryOver: req.MaxCarryOver,
		PriceType:    req.PriceType,
		Price:        req.Price,
		NewAPIGroup:  req.NewAPIGroup,
		Status:       req.Status,
		SortOrder:    req.SortOrder,
	}

	if err := model.DB.Create(plan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Message: "创建失败",
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Success: true,
		Data:    plan,
	})
}

// AdminUpdatePlan 更新套餐
func AdminUpdatePlan(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "无效的套餐 ID",
		})
		return
	}

	var plan model.Plan
	if err := model.DB.First(&plan, id).Error; err != nil {
		c.JSON(http.StatusNotFound, dto.Response{
			Success: false,
			Message: "套餐不存在",
		})
		return
	}

	var req dto.UpdatePlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "参数错误",
		})
		return
	}

	if req.Name != "" {
		plan.Name = req.Name
	}
	if req.Description != "" {
		plan.Description = req.Description
	}
	if req.PeriodType != "" {
		plan.PeriodType = req.PeriodType
	}
	if req.PeriodDays > 0 {
		plan.PeriodDays = req.PeriodDays
	}
	if req.DailyQuota > 0 {
		plan.DailyQuota = req.DailyQuota
	}
	plan.CarryOver = req.CarryOver
	plan.MaxCarryOver = req.MaxCarryOver
	if req.PriceType != "" {
		plan.PriceType = req.PriceType
	}
	if req.Price > 0 {
		plan.Price = req.Price
	}
	if req.NewAPIGroup != "" {
		plan.NewAPIGroup = req.NewAPIGroup
	}
	plan.Status = req.Status
	plan.SortOrder = req.SortOrder

	if err := model.DB.Save(&plan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Message: "更新失败",
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Success: true,
		Data:    plan,
	})
}

// AdminDeletePlan 删除套餐
func AdminDeletePlan(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "无效的套餐 ID",
		})
		return
	}

	// 检查是否有活跃订阅使用该套餐
	var count int64
	model.DB.Model(&model.Subscription{}).
		Where("plan_id = ? AND status = ?", id, model.SubscriptionStatusActive).
		Count(&count)

	if count > 0 {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "该套餐有活跃订阅，无法删除",
		})
		return
	}

	if err := model.DB.Delete(&model.Plan{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Message: "删除失败",
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: "删除成功",
	})
}

// AdminGetSettings 获取系统设置
func AdminGetSettings(c *gin.Context) {
	var settings []model.Setting
	model.DB.Find(&settings)

	settingsMap := make(map[string]string)
	for _, s := range settings {
		settingsMap[s.Key] = s.Value
	}

	c.JSON(http.StatusOK, dto.Response{
		Success: true,
		Data:    settingsMap,
	})
}

// AdminUpdateSettings 更新系统设置
func AdminUpdateSettings(c *gin.Context) {
	var req map[string]string
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "参数错误",
		})
		return
	}

	for key, value := range req {
		model.SetSetting(key, value)
	}

	c.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: "更新成功",
	})
}

// AdminTriggerSync 手动触发同步
func AdminTriggerSync(c *gin.Context) {
	go service.SyncAllSubscriptions()

	c.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: "同步任务已启动",
	})
}

// AdminGetNewAPIGroups 获取 new-api 分组
func AdminGetNewAPIGroups(c *gin.Context) {
	client := service.NewNewAPIClient()
	groups, err := client.GetGroups()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Message: "获取分组失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Success: true,
		Data:    groups,
	})
}

// AdminCompleteOrder 手动补单
func AdminCompleteOrder(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "无效的订单 ID",
		})
		return
	}

	var order model.Order
	if err := model.DB.First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, dto.Response{
			Success: false,
			Message: "订单不存在",
		})
		return
	}

	if order.Status != model.OrderStatusPending {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "只能补单待支付状态的订单",
		})
		return
	}

	// 手动补单，交易号设为 MANUAL
	if err := service.CompleteOrder(&order, "MANUAL_"+order.OrderNo); err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Message: "补单失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: "补单成功",
	})
}
