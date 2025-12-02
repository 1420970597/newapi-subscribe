package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"newapi-subscribe/internal/config"
	"newapi-subscribe/internal/dto"
	"newapi-subscribe/internal/model"
)

type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     int    `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken 生成 JWT Token
func GenerateToken(user *model.User) (string, error) {
	claims := Claims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(jwt.NewNumericDate(nil).Add(24 * 7 * 3600 * 1e9)), // 7天
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.Cfg.JWTSecret))
}

// AuthMiddleware 认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := extractToken(c)
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, dto.Response{
				Success: false,
				Message: "未授权访问",
			})
			c.Abort()
			return
		}

		claims, err := parseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, dto.Response{
				Success: false,
				Message: "Token 无效或已过期",
			})
			c.Abort()
			return
		}

		// 查询用户
		var user model.User
		if err := model.DB.First(&user, claims.UserID).Error; err != nil {
			c.JSON(http.StatusUnauthorized, dto.Response{
				Success: false,
				Message: "用户不存在",
			})
			c.Abort()
			return
		}

		if user.Status != model.StatusEnabled {
			c.JSON(http.StatusForbidden, dto.Response{
				Success: false,
				Message: "账号已被禁用",
			})
			c.Abort()
			return
		}

		c.Set("user", &user)
		c.Set("userID", user.ID)
		c.Next()
	}
}

// AdminMiddleware 管理员中间件
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, dto.Response{
				Success: false,
				Message: "未授权访问",
			})
			c.Abort()
			return
		}

		u := user.(*model.User)
		if !u.IsAdmin() {
			c.JSON(http.StatusForbidden, dto.Response{
				Success: false,
				Message: "需要管理员权限",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// OptionalAuthMiddleware 可选认证中间件
func OptionalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := extractToken(c)
		if tokenString == "" {
			c.Next()
			return
		}

		claims, err := parseToken(tokenString)
		if err != nil {
			c.Next()
			return
		}

		var user model.User
		if err := model.DB.First(&user, claims.UserID).Error; err == nil && user.Status == model.StatusEnabled {
			c.Set("user", &user)
			c.Set("userID", user.ID)
		}

		c.Next()
	}
}

func extractToken(c *gin.Context) string {
	// 从 Authorization header 获取
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
			return parts[1]
		}
	}

	// 从 query 参数获取
	if token := c.Query("token"); token != "" {
		return token
	}

	return ""
}

func parseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Cfg.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}

// GetCurrentUser 获取当前用户
func GetCurrentUser(c *gin.Context) *model.User {
	if user, exists := c.Get("user"); exists {
		return user.(*model.User)
	}
	return nil
}
