package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"newapi-subscribe/internal/config"
	"newapi-subscribe/internal/cron"
	"newapi-subscribe/internal/model"
	"newapi-subscribe/internal/router"
)

func main() {
	// 加载 .env 文件
	if err := godotenv.Load(); err != nil {
		log.Println("未找到 .env 文件，使用环境变量")
	}

	// 加载配置
	config.Load()

	// 初始化数据库
	if err := model.InitDB(config.Cfg.DBPath); err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}

	// 启动定时任务
	cron.Start()

	// 设置路由
	r := router.SetupRouter()

	// 优雅关闭
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		log.Println("正在关闭服务...")
		cron.Stop()
		os.Exit(0)
	}()

	// 启动服务
	log.Printf("服务启动在 :%s", config.Cfg.Port)
	if err := r.Run(":" + config.Cfg.Port); err != nil {
		log.Fatalf("启动服务失败: %v", err)
	}
}
