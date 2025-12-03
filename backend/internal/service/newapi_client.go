package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"newapi-subscribe/internal/config"
)

// NewAPIClient new-api HTTP 客户端
type NewAPIClient struct {
	baseURL     string
	adminID     string
	httpClient  *http.Client
	cookies     []*http.Cookie
	currentUser *NewAPIUser // 当前登录用户信息
}

// NewAPIUser new-api 用户信息
type NewAPIUser struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Role      int    `json:"role"`
	Status    int    `json:"status"`
	Quota     int    `json:"quota"`
	UsedQuota int    `json:"used_quota"`
	Group     string `json:"group"`
}

// NewAPILog new-api 日志
type NewAPILog struct {
	ID               int    `json:"id"`
	UserID           int    `json:"user_id"`
	CreatedAt        int64  `json:"created_at"`
	Type             int    `json:"type"`
	Content          string `json:"content"`
	ModelName        string `json:"model_name"`
	Quota            int    `json:"quota"`
	PromptTokens     int    `json:"prompt_tokens"`
	CompletionTokens int    `json:"completion_tokens"`
}

// NewNewAPIClient 创建 new-api 客户端
func NewNewAPIClient() *NewAPIClient {
	jar, _ := cookiejar.New(nil)
	baseURL := config.Cfg.NewAPIURL
	// 移除末尾斜杠
	baseURL = strings.TrimSuffix(baseURL, "/")
	return &NewAPIClient{
		baseURL: baseURL,
		adminID: config.Cfg.NewAPIAdminID,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
			Jar:     jar,
		},
	}
}

// AdminLogin 管理员登录
func (c *NewAPIClient) AdminLogin() error {
	_, err := c.login(config.Cfg.NewAPIAdminUser, config.Cfg.NewAPIAdminPass)
	return err
}

// Login 用户登录并返回用户信息
func (c *NewAPIClient) Login(username, password string) (*NewAPIUser, error) {
	return c.login(username, password)
}

func (c *NewAPIClient) login(username, password string) (*NewAPIUser, error) {
	data := map[string]string{
		"username": username,
		"password": password,
	}
	body, _ := json.Marshal(data)

	resp, err := c.httpClient.Post(c.baseURL+"/api/user/login", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	var result struct {
		Success bool       `json:"success"`
		Message string     `json:"message"`
		Data    NewAPIUser `json:"data"`
	}

	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	if !result.Success {
		return nil, errors.New(result.Message)
	}

	c.cookies = resp.Cookies()
	c.currentUser = &result.Data

	return &result.Data, nil
}

// GetSelf 获取当前登录用户信息
func (c *NewAPIClient) GetSelf() (*NewAPIUser, error) {
	// 如果已经有用户信息，直接返回
	if c.currentUser != nil {
		return c.currentUser, nil
	}

	req, _ := http.NewRequest("GET", c.baseURL+"/api/user/self", nil)
	for _, cookie := range c.cookies {
		req.AddCookie(cookie)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	var result struct {
		Success bool       `json:"success"`
		Message string     `json:"message"`
		Data    NewAPIUser `json:"data"`
	}

	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v, body: %s", err, string(respBody))
	}

	if !result.Success {
		return nil, errors.New(result.Message)
	}

	return &result.Data, nil
}

// addAdminHeader 添加管理员 header
func (c *NewAPIClient) addAdminHeader(req *http.Request) {
	for _, cookie := range c.cookies {
		req.AddCookie(cookie)
	}
	if c.adminID != "" {
		req.Header.Set("New-Api-User", c.adminID)
	}
}

// GetUser 获取用户信息（需要管理员权限）
func (c *NewAPIClient) GetUser(userID int) (*NewAPIUser, error) {
	if err := c.AdminLogin(); err != nil {
		return nil, fmt.Errorf("管理员登录失败: %v", err)
	}

	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/api/user/%d", c.baseURL, userID), nil)
	c.addAdminHeader(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	var result struct {
		Success bool       `json:"success"`
		Message string     `json:"message"`
		Data    NewAPIUser `json:"data"`
	}

	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v, body: %s", err, string(respBody))
	}

	if !result.Success {
		return nil, errors.New(result.Message)
	}

	return &result.Data, nil
}

// UpdateUser 更新用户（需要管理员权限）
func (c *NewAPIClient) UpdateUser(user *NewAPIUser) error {
	if err := c.AdminLogin(); err != nil {
		return fmt.Errorf("管理员登录失败: %v", err)
	}

	body, _ := json.Marshal(user)
	req, _ := http.NewRequest("PUT", c.baseURL+"/api/user/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	c.addAdminHeader(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	var result struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}

	if err := json.Unmarshal(respBody, &result); err != nil {
		return fmt.Errorf("解析响应失败: %v, body: %s", err, string(respBody))
	}

	if !result.Success {
		return errors.New(result.Message)
	}

	return nil
}

// CreateUser 创建用户（需要管理员权限）
func (c *NewAPIClient) CreateUser(username, password, group string) (*NewAPIUser, error) {
	if err := c.AdminLogin(); err != nil {
		return nil, fmt.Errorf("管理员登录失败: %v", err)
	}

	data := map[string]interface{}{
		"username":     username,
		"password":     password,
		"display_name": username,
		"group":        group,
	}
	body, _ := json.Marshal(data)

	req, _ := http.NewRequest("POST", c.baseURL+"/api/user/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	c.addAdminHeader(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	var result struct {
		Success bool       `json:"success"`
		Message string     `json:"message"`
		Data    NewAPIUser `json:"data"`
	}

	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v, body: %s", err, string(respBody))
	}

	if !result.Success {
		return nil, fmt.Errorf("创建用户失败: %s", result.Message)
	}

	return &result.Data, nil
}

// GetGroups 获取分组列表（需要管理员权限）
func (c *NewAPIClient) GetGroups() ([]string, error) {
	if err := c.AdminLogin(); err != nil {
		return nil, fmt.Errorf("管理员登录失败: %v", err)
	}

	req, _ := http.NewRequest("GET", c.baseURL+"/api/group/", nil)
	c.addAdminHeader(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	var result struct {
		Success bool     `json:"success"`
		Message string   `json:"message"`
		Data    []string `json:"data"`
	}

	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v, body: %s", err, string(respBody))
	}

	if !result.Success {
		return nil, fmt.Errorf("获取分组失败: %s", result.Message)
	}

	return result.Data, nil
}

// GetGroupModels 获取分组下的模型
func (c *NewAPIClient) GetGroupModels(group string) (interface{}, error) {
	return map[string]string{"group": group}, nil
}

// GetUserLogs 获取用户日志（需要管理员权限）
func (c *NewAPIClient) GetUserLogs(userID int, startDate, endDate string) ([]NewAPILog, error) {
	if err := c.AdminLogin(); err != nil {
		return nil, fmt.Errorf("管理员登录失败: %v", err)
	}

	params := url.Values{}
	if startDate != "" {
		params.Set("start_timestamp", startDate)
	}
	if endDate != "" {
		params.Set("end_timestamp", endDate)
	}

	reqURL := fmt.Sprintf("%s/api/log?user_id=%d&%s", c.baseURL, userID, params.Encode())
	req, _ := http.NewRequest("GET", reqURL, nil)
	c.addAdminHeader(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	var result struct {
		Success bool        `json:"success"`
		Message string      `json:"message"`
		Data    []NewAPILog `json:"data"`
	}

	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	if !result.Success {
		return nil, errors.New(result.Message)
	}

	return result.Data, nil
}

// GetUserQuotaUsedToday 获取用户今日已用额度
func (c *NewAPIClient) GetUserQuotaUsedToday(userID int) (int, error) {
	if err := c.AdminLogin(); err != nil {
		return 0, fmt.Errorf("管理员登录失败: %v", err)
	}

	// 获取今日时间范围
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	params := url.Values{}
	params.Set("user_id", fmt.Sprintf("%d", userID))
	params.Set("start_timestamp", fmt.Sprintf("%d", startOfDay.Unix()))
	params.Set("end_timestamp", fmt.Sprintf("%d", endOfDay.Unix()))

	reqURL := fmt.Sprintf("%s/api/log/?%s", c.baseURL, params.Encode())
	req, _ := http.NewRequest("GET", reqURL, nil)
	c.addAdminHeader(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	var result struct {
		Success bool        `json:"success"`
		Message string      `json:"message"`
		Data    []NewAPILog `json:"data"`
	}

	if err := json.Unmarshal(respBody, &result); err != nil {
		return 0, fmt.Errorf("解析响应失败: %v", err)
	}

	if !result.Success {
		return 0, errors.New(result.Message)
	}

	// 计算今日已用额度
	var totalUsed int
	for _, log := range result.Data {
		totalUsed += log.Quota
	}

	return totalUsed, nil
}
