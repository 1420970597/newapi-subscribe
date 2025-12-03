package router

import (
	"github.com/gin-gonic/gin"
	"newapi-subscribe/internal/controller"
	"newapi-subscribe/internal/middleware"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 中间件
	r.Use(middleware.CORSMiddleware())

	// API 路由组
	api := r.Group("/api")
	{
		// 认证接口
		auth := api.Group("/auth")
		{
			auth.POST("/register", controller.Register)
			auth.POST("/login", controller.Login)
			auth.POST("/login/newapi", controller.NewAPILogin)
			auth.POST("/logout", middleware.AuthMiddleware(), controller.Logout)
			auth.GET("/me", middleware.AuthMiddleware(), controller.GetCurrentUser)
		}

		// 套餐接口
		plans := api.Group("/plans")
		{
			plans.GET("", middleware.OptionalAuthMiddleware(), controller.GetPlans)
			plans.GET("/:id", middleware.OptionalAuthMiddleware(), controller.GetPlan)
			plans.GET("/:id/models", controller.GetPlanModels)
		}

		// 订阅接口（需要登录）
		subscriptions := api.Group("/subscriptions")
		subscriptions.Use(middleware.AuthMiddleware())
		{
			subscriptions.GET("/current", controller.GetCurrentSubscription)
			subscriptions.POST("/purchase", controller.PurchaseSubscription)
			subscriptions.POST("/renew", controller.RenewSubscription)
			subscriptions.GET("/usage", controller.GetUsageLogs)
			subscriptions.GET("/usage/detail", controller.GetUsageDetail)
		}

		// 订单接口
		orders := api.Group("/orders")
		{
			orders.GET("/notify", controller.PaymentNotify) // 支付回调（公开）
			orders.Use(middleware.AuthMiddleware())
			orders.GET("", controller.GetOrders)
			orders.GET("/:id", controller.GetOrder)
			orders.POST("/pay", controller.CreatePayment)
		}

		// 用户接口（需要登录）
		user := api.Group("/user")
		user.Use(middleware.AuthMiddleware())
		{
			user.PUT("/profile", controller.UpdateProfile)
			user.POST("/bind-newapi", controller.BindNewAPI)
			user.PUT("/email-settings", controller.UpdateEmailSettings)
		}

		// 管理接口（需要管理员权限）
		admin := api.Group("/admin")
		admin.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
		{
			// 用户管理
			admin.GET("/users", controller.AdminGetUsers)
			admin.GET("/users/:id", controller.AdminGetUser)
			admin.GET("/users/:id/usage", controller.AdminGetUserUsage)
			admin.PUT("/users/:id", controller.AdminUpdateUser)

			// 订阅管理
			admin.GET("/subscriptions", controller.AdminGetSubscriptions)

			// 订单管理
			admin.GET("/orders", controller.AdminGetOrders)
			admin.POST("/orders/:id/complete", controller.AdminCompleteOrder)

			// 套餐管理
			admin.POST("/plans", controller.AdminCreatePlan)
			admin.PUT("/plans/:id", controller.AdminUpdatePlan)
			admin.DELETE("/plans/:id", controller.AdminDeletePlan)

			// 系统设置
			admin.GET("/settings", controller.AdminGetSettings)
			admin.PUT("/settings", controller.AdminUpdateSettings)

			// 同步操作
			admin.POST("/sync/trigger", controller.AdminTriggerSync)

			// new-api 信息
			admin.GET("/newapi/groups", controller.AdminGetNewAPIGroups)
		}
	}

	// 静态文件（前端）
	r.Static("/assets", "./static/assets")
	r.StaticFile("/vite.svg", "./static/vite.svg")
	r.NoRoute(func(c *gin.Context) {
		c.File("./static/index.html")
	})

	return r
}
