package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"newapi-subscribe/internal/dto"
	"newapi-subscribe/internal/middleware"
	"newapi-subscribe/internal/model"
	"newapi-subscribe/internal/service"
)

// GetOrders 获取订单列表
func GetOrders(c *gin.Context) {
	user := middleware.GetCurrentUser(c)

	var pagination dto.PaginationQuery
	if err := c.ShouldBindQuery(&pagination); err != nil {
		pagination.Page = 1
		pagination.PerPage = 20
	}

	var orders []model.Order
	var total int64

	model.DB.Model(&model.Order{}).Where("user_id = ?", user.ID).Count(&total)
	model.DB.Preload("Plan").
		Where("user_id = ?", user.ID).
		Order("created_at DESC").
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

// GetOrder 获取订单详情
func GetOrder(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "无效的订单 ID",
		})
		return
	}

	var order model.Order
	if err := model.DB.Preload("Plan").
		Where("id = ? AND user_id = ?", id, user.ID).
		First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, dto.Response{
			Success: false,
			Message: "订单不存在",
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Success: true,
		Data:    order,
	})
}

// CreatePayment 发起支付
func CreatePayment(c *gin.Context) {
	var req dto.PayRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "参数错误",
		})
		return
	}

	user := middleware.GetCurrentUser(c)

	var order model.Order
	if err := model.DB.Where("id = ? AND user_id = ?", req.OrderID, user.ID).First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, dto.Response{
			Success: false,
			Message: "订单不存在",
		})
		return
	}

	if order.Status != model.OrderStatusPending {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "订单状态异常",
		})
		return
	}

	// 调用易支付
	epay := service.NewEpayService()
	payURL, err := epay.CreatePayment(order.OrderNo, order.Amount, req.PaymentMethod, "订阅套餐")
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Message: "创建支付失败: " + err.Error(),
		})
		return
	}

	// 更新支付方式
	order.PaymentMethod = req.PaymentMethod
	model.DB.Save(&order)

	c.JSON(http.StatusOK, dto.Response{
		Success: true,
		Data: gin.H{
			"pay_url": payURL,
		},
	})
}

// PaymentNotify 支付回调
func PaymentNotify(c *gin.Context) {
	epay := service.NewEpayService()

	// 验证签名
	if !epay.VerifyNotify(c) {
		c.String(http.StatusBadRequest, "fail")
		return
	}

	orderNo := c.Query("out_trade_no")
	tradeNo := c.Query("trade_no")
	tradeStatus := c.Query("trade_status")

	if tradeStatus != "TRADE_SUCCESS" {
		c.String(http.StatusOK, "success")
		return
	}

	// 查找订单
	var order model.Order
	if err := model.DB.Where("order_no = ?", orderNo).First(&order).Error; err != nil {
		c.String(http.StatusOK, "fail")
		return
	}

	if order.Status != model.OrderStatusPending {
		c.String(http.StatusOK, "success")
		return
	}

	// 处理订单完成
	if err := service.CompleteOrder(&order, tradeNo); err != nil {
		c.String(http.StatusOK, "fail")
		return
	}

	c.String(http.StatusOK, "success")
}
