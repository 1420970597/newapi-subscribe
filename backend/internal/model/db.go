package model

import (
	"log"
	"os"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDB 初始化数据库
func InitDB(dbPath string) error {
	// 确保目录存在
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	var err error
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return err
	}

	// 自动迁移
	if err := DB.AutoMigrate(
		&User{},
		&Plan{},
		&Subscription{},
		&Order{},
		&Setting{},
		&UsageLog{},
	); err != nil {
		return err
	}

	// 初始化默认设置
	initDefaultSettings()

	// 初始化管理员账号
	initAdminUser()

	return nil
}

// initDefaultSettings 初始化默认设置
func initDefaultSettings() {
	for key, value := range DefaultSettings {
		var setting Setting
		result := DB.Where("key = ?", key).First(&setting)
		if result.Error == gorm.ErrRecordNotFound {
			DB.Create(&Setting{Key: key, Value: value})
		}
	}
}

// initAdminUser 初始化管理员账号
func initAdminUser() {
	var count int64
	DB.Model(&User{}).Where("role >= ?", RoleAdmin).Count(&count)
	if count == 0 {
		admin := &User{
			Username: "admin",
			Role:     RoleAdmin,
			Status:   StatusEnabled,
		}
		admin.SetPassword("admin123")
		if err := DB.Create(admin).Error; err != nil {
			log.Printf("创建管理员账号失败: %v", err)
		} else {
			log.Println("已创建默认管理员账号: admin / admin123")
		}
	}
}

// GetSetting 获取设置值
func GetSetting(key string) string {
	var setting Setting
	if err := DB.Where("key = ?", key).First(&setting).Error; err != nil {
		return DefaultSettings[key]
	}
	return setting.Value
}

// SetSetting 设置值
func SetSetting(key, value string) error {
	return DB.Save(&Setting{Key: key, Value: value}).Error
}
