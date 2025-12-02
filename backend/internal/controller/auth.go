package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"newapi-subscribe/internal/config"
	"newapi-subscribe/internal/dto"
	"newapi-subscribe/internal/middleware"
	"newapi-subscribe/internal/model"
	"newapi-subscribe/internal/service"
)

// Register 用户注册
func Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	// 检查是否允许注册
	if model.GetSetting(model.SettingAllowRegister) != "1" {
		c.JSON(http.StatusForbidden, dto.Response{
			Success: false,
			Message: "注册功能已关闭",
		})
		return
	}

	// 检查用户名是否已存在
	var existingUser model.User
	if err := model.DB.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "用户名已存在",
		})
		return
	}

	// 创建用户
	user := &model.User{
		Username: req.Username,
		Email:    req.Email,
		Role:     model.RoleUser,
		Status:   model.StatusEnabled,
	}
	if err := user.SetPassword(req.Password); err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Message: "密码加密失败",
		})
		return
	}

	if err := model.DB.Create(user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Message: "创建用户失败",
		})
		return
	}

	// 生成 token
	token, err := generateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Message: "生成 Token 失败",
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Success: true,
		Data: gin.H{
			"token": token,
			"user":  user,
		},
	})
}

// Login 用户登录
func Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "参数错误",
		})
		return
	}

	var user model.User
	if err := model.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, dto.Response{
			Success: false,
			Message: "用户名或密码错误",
		})
		return
	}

	if !user.CheckPassword(req.Password) {
		c.JSON(http.StatusUnauthorized, dto.Response{
			Success: false,
			Message: "用户名或密码错误",
		})
		return
	}

	if user.Status != model.StatusEnabled {
		c.JSON(http.StatusForbidden, dto.Response{
			Success: false,
			Message: "账号已被禁用",
		})
		return
	}

	token, err := generateToken(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Message: "生成 Token 失败",
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Success: true,
		Data: gin.H{
			"token": token,
			"user":  user,
		},
	})
}

// NewAPILogin 使用 new-api 账号登录
func NewAPILogin(c *gin.Context) {
	var req dto.NewAPILoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "参数错误",
		})
		return
	}

	// 检查是否启用 new-api 登录
	if model.GetSetting(model.SettingNewAPILoginEnabled) != "1" {
		c.JSON(http.StatusForbidden, dto.Response{
			Success: false,
			Message: "new-api 登录功能已关闭",
		})
		return
	}

	// 验证 new-api 账号
	client := service.NewNewAPIClient()
	newAPIUser, err := client.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.Response{
			Success: false,
			Message: "new-api 账号验证失败: " + err.Error(),
		})
		return
	}

	// 查找或创建本地用户
	var user model.User
	if err := model.DB.Where("newapi_user_id = ?", newAPIUser.ID).First(&user).Error; err != nil {
		// 创建新用户
		user = model.User{
			Username:       req.Username,
			Role:           model.RoleUser,
			Status:         model.StatusEnabled,
			NewAPIUserID:   newAPIUser.ID,
			NewAPIUsername: newAPIUser.Username,
			NewAPIBound:    1,
		}
		if err := model.DB.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, dto.Response{
				Success: false,
				Message: "创建用户失败",
			})
			return
		}
	}

	if user.Status != model.StatusEnabled {
		c.JSON(http.StatusForbidden, dto.Response{
			Success: false,
			Message: "账号已被禁用",
		})
		return
	}

	token, err := generateToken(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Message: "生成 Token 失败",
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Success: true,
		Data: gin.H{
			"token": token,
			"user":  user,
		},
	})
}

// Logout 登出
func Logout(c *gin.Context) {
	c.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: "登出成功",
	})
}

// GetCurrentUser 获取当前用户信息
func GetCurrentUser(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, dto.Response{
			Success: false,
			Message: "未登录",
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Success: true,
		Data:    user,
	})
}

func generateToken(user *model.User) (string, error) {
	claims := middleware.Claims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.Cfg.JWTSecret))
}
