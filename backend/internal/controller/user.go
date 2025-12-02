package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"newapi-subscribe/internal/dto"
	"newapi-subscribe/internal/middleware"
	"newapi-subscribe/internal/model"
	"newapi-subscribe/internal/service"
)

// UpdateProfile 更新个人信息
func UpdateProfile(c *gin.Context) {
	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "参数错误",
		})
		return
	}

	user := middleware.GetCurrentUser(c)
	user.Email = req.Email

	if err := model.DB.Save(user).Error; err != nil {
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

// BindNewAPI 绑定 new-api 账号
func BindNewAPI(c *gin.Context) {
	var req dto.BindNewAPIRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "参数错误",
		})
		return
	}

	user := middleware.GetCurrentUser(c)

	if user.NewAPIBound == 1 {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "您已绑定 new-api 账号",
		})
		return
	}

	// 验证 new-api 账号
	client := service.NewNewAPIClient()
	newAPIUser, err := client.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "new-api 账号验证失败",
		})
		return
	}

	// 检查是否已被其他用户绑定
	var existingUser model.User
	if err := model.DB.Where("newapi_user_id = ? AND id != ?", newAPIUser.ID, user.ID).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "该 new-api 账号已被其他用户绑定",
		})
		return
	}

	user.NewAPIUserID = newAPIUser.ID
	user.NewAPIUsername = newAPIUser.Username
	user.NewAPIBound = 1

	if err := model.DB.Save(user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Message: "绑定失败",
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: "绑定成功",
		Data:    user,
	})
}

// UpdateEmailSettings 更新邮件提醒设置
func UpdateEmailSettings(c *gin.Context) {
	var req dto.UpdateEmailSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "参数错误",
		})
		return
	}

	user := middleware.GetCurrentUser(c)
	user.EmailRemind = req.EmailRemind
	user.RemindDays = req.RemindDays

	if err := model.DB.Save(user).Error; err != nil {
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
