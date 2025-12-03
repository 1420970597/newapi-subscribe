package service

import (
	"log"
	"time"

	"newapi-subscribe/internal/model"
)

// SyncAllSubscriptions 同步所有订阅的额度
func SyncAllSubscriptions() {
	log.Println("开始执行订阅额度同步...")

	client := NewNewAPIClient()
	if err := client.AdminLogin(); err != nil {
		log.Printf("管理员登录失败: %v", err)
		return
	}

	today := time.Now().Truncate(24 * time.Hour)

	// 获取所有活跃订阅
	var subscriptions []model.Subscription
	model.DB.Preload("User").
		Where("status = ?", model.SubscriptionStatusActive).
		Find(&subscriptions)

	for _, sub := range subscriptions {
		if err := syncSubscription(client, &sub, today); err != nil {
			log.Printf("同步订阅 %d 失败: %v", sub.ID, err)
		}
	}

	// 发送到期提醒
	sendExpirationReminders()

	log.Println("订阅额度同步完成")
}

// syncSubscription 同步单个订阅
func syncSubscription(client *NewAPIClient, sub *model.Subscription, today time.Time) error {
	// 检查是否过期
	if sub.EndDate.Before(today) {
		sub.Status = model.SubscriptionStatusExpired
		model.DB.Save(sub)

		// 清零 new-api 余额
		if sub.User != nil && sub.User.NewAPIBound == 1 {
			newAPIUser, err := client.GetUser(sub.User.NewAPIUserID)
			if err == nil {
				newAPIUser.Quota = 0
				client.UpdateUser(newAPIUser)
			}
		}
		log.Printf("订阅 %d 已过期", sub.ID)
		return nil
	}

	// 获取用户
	if sub.User == nil || sub.User.NewAPIBound != 1 {
		return nil
	}

	// 获取 new-api 当前余额（昨日剩余）
	newAPIUser, err := client.GetUser(sub.User.NewAPIUserID)
	if err != nil {
		return err
	}
	yesterdayRemaining := newAPIUser.Quota

	// 计算新的每日额度
	var newQuota int
	var carriedQuota int

	if sub.CarryOver == 1 && yesterdayRemaining > 0 {
		// 计算结转额度
		carriedQuota = yesterdayRemaining
		if sub.MaxCarryOver > 0 && carriedQuota > sub.MaxCarryOver {
			carriedQuota = sub.MaxCarryOver
		}
		newQuota = sub.DailyQuota + carriedQuota
	} else {
		newQuota = sub.DailyQuota
		carriedQuota = 0
	}

	// 更新 new-api 用户余额和分组
	newAPIUser.Quota = newQuota
	newAPIUser.Group = sub.NewAPIGroup
	if err := client.UpdateUser(newAPIUser); err != nil {
		return err
	}

	// 更新本地记录
	sub.TodayQuota = newQuota
	sub.CarriedQuota = carriedQuota
	sub.LastSyncDate = &today
	model.DB.Save(sub)

	log.Printf("订阅 %d 同步完成: 每日额度=%d, 结转=%d, 新额度=%d",
		sub.ID, sub.DailyQuota, carriedQuota, newQuota)

	return nil
}

// sendExpirationReminders 发送到期提醒
func sendExpirationReminders() {
	today := time.Now().Truncate(24 * time.Hour)

	var subscriptions []model.Subscription
	model.DB.Preload("User").Preload("Plan").
		Where("status = ?", model.SubscriptionStatusActive).
		Find(&subscriptions)

	for _, sub := range subscriptions {
		if sub.User == nil || sub.User.EmailRemind != 1 || sub.User.Email == "" {
			continue
		}

		daysRemaining := int(sub.EndDate.Sub(today).Hours() / 24)
		if daysRemaining <= sub.User.RemindDays && daysRemaining >= 0 {
			// 发送提醒邮件
			planName := ""
			if sub.Plan != nil {
				planName = sub.Plan.Name
			}
			SendExpirationReminder(sub.User.Email, sub.User.Username, planName, daysRemaining)
		}
	}
}

// CompleteOrder 完成订单
func CompleteOrder(order *model.Order, tradeNo string) error {
	// 更新订单状态
	now := time.Now()
	order.Status = model.OrderStatusPaid
	order.TradeNo = tradeNo
	order.PaidAt = &now

	if err := model.DB.Save(order).Error; err != nil {
		return err
	}

	// 获取用户和套餐
	var user model.User
	var plan model.Plan
	model.DB.First(&user, order.UserID)
	model.DB.First(&plan, order.PlanID)

	// 处理 new-api 账号
	client := NewNewAPIClient()

	if user.NewAPIBound != 1 {
		// 创建新账号
		newAPIUser, err := client.CreateUser(user.Username, "temp_password", plan.NewAPIGroup)
		if err != nil {
			log.Printf("创建 new-api 账号失败: %v", err)
		} else {
			user.NewAPIUserID = newAPIUser.ID
			user.NewAPIUsername = newAPIUser.Username
			user.NewAPIBound = 1
			model.DB.Save(&user)
		}
	}

	// 创建或更新订阅
	var subscription model.Subscription
	today := time.Now().Truncate(24 * time.Hour)

	// 计算到期时间的辅助函数
	calcEndDate := func(baseDate time.Time, periodDays int) time.Time {
		if plan.PeriodType == model.PeriodTypeMonth {
			// 按月订阅：计算月数
			months := periodDays / 30
			if months < 1 {
				months = 1
			}
			return baseDate.AddDate(0, months, 0)
		} else if plan.PeriodType == model.PeriodTypeWeek {
			// 按周订阅：计算周数
			weeks := periodDays / 7
			if weeks < 1 {
				weeks = 1
			}
			return baseDate.AddDate(0, 0, weeks*7)
		}
		// 按天或自定义
		return baseDate.AddDate(0, 0, periodDays)
	}

	if order.OrderType == model.OrderTypeRenew {
		// 续费：延长现有订阅
		model.DB.Where("user_id = ? AND status = ?", user.ID, model.SubscriptionStatusActive).
			First(&subscription)

		if subscription.ID > 0 {
			subscription.EndDate = calcEndDate(subscription.EndDate, order.PeriodDays)
			model.DB.Save(&subscription)
		}
	} else {
		// 新购：创建新订阅
		// 先将旧订阅设为过期
		model.DB.Model(&model.Subscription{}).
			Where("user_id = ? AND status = ?", user.ID, model.SubscriptionStatusActive).
			Update("status", model.SubscriptionStatusExpired)

		subscription = model.Subscription{
			UserID:       user.ID,
			PlanID:       plan.ID,
			Status:       model.SubscriptionStatusActive,
			StartDate:    today,
			EndDate:      calcEndDate(today, order.PeriodDays),
			TodayQuota:   plan.DailyQuota,
			DailyQuota:   plan.DailyQuota,
			CarryOver:    plan.CarryOver,
			MaxCarryOver: plan.MaxCarryOver,
			NewAPIGroup:  plan.NewAPIGroup,
			LastSyncDate: &today,
		}
		model.DB.Create(&subscription)
	}

	// 设置 new-api 初始额度
	if user.NewAPIBound == 1 {
		newAPIUser, err := client.GetUser(user.NewAPIUserID)
		if err == nil {
			newAPIUser.Quota = plan.DailyQuota
			newAPIUser.Group = plan.NewAPIGroup
			client.UpdateUser(newAPIUser)
		}
	}

	log.Printf("订单 %s 完成，用户 %d 订阅已激活", order.OrderNo, user.ID)
	return nil
}
