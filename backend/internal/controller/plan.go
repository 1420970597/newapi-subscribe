package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"newapi-subscribe/internal/dto"
	"newapi-subscribe/internal/model"
	"newapi-subscribe/internal/service"
)

// GetPlans 获取套餐列表
func GetPlans(c *gin.Context) {
	// 检查是否需要登录
	if model.GetSetting(model.SettingRequireLogin) == "1" {
		if _, exists := c.Get("user"); !exists {
			c.JSON(http.StatusUnauthorized, dto.Response{
				Success: false,
				Message: "需要登录后查看",
			})
			return
		}
	}

	var plans []model.Plan
	query := model.DB.Where("status = ?", model.PlanStatusOn).Order("sort_order ASC, id ASC")

	if err := query.Find(&plans).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Message: "获取套餐列表失败",
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Success: true,
		Data:    plans,
	})
}

// GetPlan 获取套餐详情
func GetPlan(c *gin.Context) {
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

	c.JSON(http.StatusOK, dto.Response{
		Success: true,
		Data:    plan,
	})
}

// GetPlanModels 获取套餐可用模型
func GetPlanModels(c *gin.Context) {
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

	// 从 new-api 获取分组下的模型
	client := service.NewNewAPIClient()
	models, err := client.GetGroupModels(plan.NewAPIGroup)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Message: "获取模型列表失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Success: true,
		Data:    models,
	})
}
