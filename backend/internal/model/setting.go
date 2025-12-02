package model

import (
	"time"
)

// Setting 系统设置模型
type Setting struct {
	Key       string    `gorm:"primaryKey;size:64" json:"key"`
	Value     string    `gorm:"type:text" json:"value"`
	UpdatedAt time.Time `json:"updated_at"`
}

// 默认设置键
const (
	SettingSiteName          = "site_name"
	SettingSiteDescription   = "site_description"
	SettingRequireLogin      = "require_login"
	SettingAllowRegister     = "allow_register"
	SettingNewAPILoginEnabled = "newapi_login_enabled"
)

// DefaultSettings 默认设置
var DefaultSettings = map[string]string{
	SettingSiteName:          "订阅中心",
	SettingSiteDescription:   "AI 模型订阅服务",
	SettingRequireLogin:      "0",
	SettingAllowRegister:     "1",
	SettingNewAPILoginEnabled: "1",
}
