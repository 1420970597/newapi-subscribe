package config

import (
	"os"
	"strconv"
)

type Config struct {
	// 服务配置
	Port      string
	JWTSecret string

	// 数据库
	DBPath string

	// new-api 配置
	NewAPIURL       string
	NewAPIAdminUser string
	NewAPIAdminPass string

	// 易支付配置
	EpayURL string
	EpayPID string
	EpayKey string

	// SMTP 配置
	SMTPServer string
	SMTPPort   int
	SMTPUser   string
	SMTPPass   string
	SMTPFrom   string

	// 定时任务
	CronEnabled  bool
	CronSchedule string
}

var Cfg *Config

func Load() {
	Cfg = &Config{
		Port:      getEnv("PORT", "8080"),
		JWTSecret: getEnv("JWT_SECRET", "change-me-in-production"),
		DBPath:    getEnv("DB_PATH", "./data/subscribe.db"),

		NewAPIURL:       getEnv("NEWAPI_URL", ""),
		NewAPIAdminUser: getEnv("NEWAPI_ADMIN_USER", ""),
		NewAPIAdminPass: getEnv("NEWAPI_ADMIN_PASS", ""),

		EpayURL: getEnv("EPAY_URL", ""),
		EpayPID: getEnv("EPAY_PID", ""),
		EpayKey: getEnv("EPAY_KEY", ""),

		SMTPServer: getEnv("SMTP_SERVER", ""),
		SMTPPort:   getEnvInt("SMTP_PORT", 587),
		SMTPUser:   getEnv("SMTP_USER", ""),
		SMTPPass:   getEnv("SMTP_PASS", ""),
		SMTPFrom:   getEnv("SMTP_FROM", ""),

		CronEnabled:  getEnvBool("CRON_ENABLED", true),
		CronSchedule: getEnv("CRON_SCHEDULE", "0 0 * * *"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolVal, err := strconv.ParseBool(value); err == nil {
			return boolVal
		}
	}
	return defaultValue
}
