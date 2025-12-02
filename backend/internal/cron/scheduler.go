package cron

import (
	"log"

	"github.com/robfig/cron/v3"
	"newapi-subscribe/internal/config"
	"newapi-subscribe/internal/service"
)

var scheduler *cron.Cron

// Start 启动定时任务
func Start() {
	if !config.Cfg.CronEnabled {
		log.Println("定时任务已禁用")
		return
	}

	scheduler = cron.New()

	// 每日额度同步任务
	_, err := scheduler.AddFunc(config.Cfg.CronSchedule, func() {
		log.Println("执行定时任务: 同步订阅额度")
		service.SyncAllSubscriptions()
	})

	if err != nil {
		log.Printf("添加定时任务失败: %v", err)
		return
	}

	scheduler.Start()
	log.Printf("定时任务已启动，调度: %s", config.Cfg.CronSchedule)
}

// Stop 停止定时任务
func Stop() {
	if scheduler != nil {
		scheduler.Stop()
		log.Println("定时任务已停止")
	}
}
