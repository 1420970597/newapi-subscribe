package model

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Username  string         `gorm:"uniqueIndex;size:64;not null" json:"username"`
	Password  string         `gorm:"size:255" json:"-"`
	Email     string         `gorm:"size:128" json:"email"`
	Role      int            `gorm:"default:1" json:"role"` // 1=普通用户, 10=管理员
	Status    int            `gorm:"default:1" json:"status"` // 1=启用, 2=禁用

	// new-api 绑定信息
	NewAPIUserID   int    `gorm:"column:newapi_user_id" json:"newapi_user_id"`
	NewAPIUsername string `gorm:"column:newapi_username;size:64" json:"newapi_username"`
	NewAPIBound    int    `gorm:"column:newapi_bound;default:0" json:"newapi_bound"` // 0=未绑定, 1=已绑定

	// 邮件提醒设置
	EmailRemind int `gorm:"default:1" json:"email_remind"` // 是否开启邮件提醒
	RemindDays  int `gorm:"default:3" json:"remind_days"`  // 提前几天提醒

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// SetPassword 设置密码
func (u *User) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword 校验密码
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// IsAdmin 是否为管理员
func (u *User) IsAdmin() bool {
	return u.Role >= 10
}

const (
	RoleUser  = 1
	RoleAdmin = 10

	StatusEnabled  = 1
	StatusDisabled = 2
)
