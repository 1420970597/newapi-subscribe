package service

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/url"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
	"newapi-subscribe/internal/config"
)

// EpayService 易支付服务
type EpayService struct {
	apiURL string
	pid    string
	key    string
}

// NewEpayService 创建易支付服务
func NewEpayService() *EpayService {
	return &EpayService{
		apiURL: config.Cfg.EpayURL,
		pid:    config.Cfg.EpayPID,
		key:    config.Cfg.EpayKey,
	}
}

// CreatePayment 创建支付
func (e *EpayService) CreatePayment(orderNo string, amount float64, payType, subject string) (string, error) {
	if e.apiURL == "" || e.pid == "" || e.key == "" {
		return "", fmt.Errorf("易支付未配置")
	}

	params := map[string]string{
		"pid":          e.pid,
		"type":         payType,
		"out_trade_no": orderNo,
		"notify_url":   config.Cfg.NewAPIURL + "/api/orders/notify", // 使用本系统地址
		"return_url":   config.Cfg.NewAPIURL + "/user/orders",
		"name":         subject,
		"money":        fmt.Sprintf("%.2f", amount),
	}

	// 生成签名
	params["sign"] = e.generateSign(params)
	params["sign_type"] = "MD5"

	// 构建支付 URL
	payURL := e.apiURL + "/submit.php?"
	values := url.Values{}
	for k, v := range params {
		values.Set(k, v)
	}
	payURL += values.Encode()

	return payURL, nil
}

// VerifyNotify 验证回调签名
func (e *EpayService) VerifyNotify(c *gin.Context) bool {
	params := make(map[string]string)
	for k, v := range c.Request.URL.Query() {
		if k != "sign" && k != "sign_type" && len(v) > 0 {
			params[k] = v[0]
		}
	}

	sign := c.Query("sign")
	expectedSign := e.generateSign(params)

	return sign == expectedSign
}

// generateSign 生成签名
func (e *EpayService) generateSign(params map[string]string) string {
	// 按键排序
	var keys []string
	for k := range params {
		if k != "sign" && k != "sign_type" && params[k] != "" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	// 拼接字符串
	var parts []string
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s=%s", k, params[k]))
	}
	signStr := strings.Join(parts, "&") + e.key

	// MD5
	hash := md5.Sum([]byte(signStr))
	return hex.EncodeToString(hash[:])
}
