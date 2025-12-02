package service

import (
	"fmt"
	"log"
	"net/smtp"

	"newapi-subscribe/internal/config"
	"newapi-subscribe/internal/model"
)

// SendEmail 发送邮件
func SendEmail(to, subject, body string) error {
	cfg := config.Cfg
	if cfg.SMTPServer == "" {
		return fmt.Errorf("SMTP 未配置")
	}

	from := cfg.SMTPFrom
	if from == "" {
		from = cfg.SMTPUser
	}

	// 构建邮件内容
	msg := fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/html; charset=UTF-8\r\n"+
		"\r\n%s", from, to, subject, body)

	auth := smtp.PlainAuth("", cfg.SMTPUser, cfg.SMTPPass, cfg.SMTPServer)
	addr := fmt.Sprintf("%s:%d", cfg.SMTPServer, cfg.SMTPPort)

	return smtp.SendMail(addr, auth, from, []string{to}, []byte(msg))
}

// SendExpirationReminder 发送到期提醒
func SendExpirationReminder(email, username, planName string, daysRemaining int) {
	siteName := model.GetSetting(model.SettingSiteName)

	var subject string
	var body string

	if daysRemaining == 0 {
		subject = fmt.Sprintf("[%s] 您的订阅今日到期", siteName)
		body = fmt.Sprintf(`
			<div style="font-family: sans-serif; max-width: 600px; margin: 0 auto;">
				<h2>订阅到期提醒</h2>
				<p>亲爱的 %s：</p>
				<p>您的 <strong>%s</strong> 订阅将于今日到期。</p>
				<p>为了不影响您的正常使用，请及时续费。</p>
				<p style="margin-top: 30px; color: #666;">
					—— %s
				</p>
			</div>
		`, username, planName, siteName)
	} else {
		subject = fmt.Sprintf("[%s] 您的订阅将于 %d 天后到期", siteName, daysRemaining)
		body = fmt.Sprintf(`
			<div style="font-family: sans-serif; max-width: 600px; margin: 0 auto;">
				<h2>订阅到期提醒</h2>
				<p>亲爱的 %s：</p>
				<p>您的 <strong>%s</strong> 订阅将于 <strong>%d 天后</strong>到期。</p>
				<p>为了不影响您的正常使用，请及时续费。</p>
				<p style="margin-top: 30px; color: #666;">
					—— %s
				</p>
			</div>
		`, username, planName, daysRemaining, siteName)
	}

	if err := SendEmail(email, subject, body); err != nil {
		log.Printf("发送到期提醒邮件失败: %v", err)
	} else {
		log.Printf("已发送到期提醒邮件给 %s", email)
	}
}
