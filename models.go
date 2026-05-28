package main

import "time"

// ========== 核心数据结构 ==========

// Config 应用配置
type Config struct {
	APIKey          string           `json:"api_key"`
	BaseURL         string           `json:"base_url"`
	IntervalSeconds int              `json:"interval_seconds"`
	Webhook         WebhookConfig    `json:"webhook"`
	Watchlist       []WatchItem      `json:"watchlist"`
	Channels        []Channel        `json:"channels"`
	Proxy           ProxyConfig      `json:"proxy"`
	BalanceAlert    BalanceAlert     `json:"balance_alert"`
	Users           []User           `json:"users"`
}

// WatchItem 监控项
type WatchItem struct {
	ID          string  `json:"id"`
	Service     string  `json:"service"`
	ServiceName string  `json:"service_name"`
	Country     int     `json:"country"`
	CountryName string  `json:"country_name"`
	MaxPrice    float64 `json:"max_price"`
	Enabled     bool    `json:"enabled"`
	CreatedAt   string  `json:"created_at"`
}

// ========== 通知渠道 ==========

// Channel 通知渠道
type Channel struct {
	ID       string         `json:"id"`
	Type     string         `json:"type"` // feishu, dingtalk, wechat, telegram, email
	Name     string         `json:"name"`
	Enabled  bool           `json:"enabled"`
	Config   ChannelConfig  `json:"config"`
	Policies []NotifyPolicy `json:"policies"`
}

// ChannelConfig 渠道配置
type ChannelConfig struct {
	WebhookURL string `json:"webhook_url,omitempty"`
	Sign       string `json:"sign,omitempty"`
	BotToken   string `json:"bot_token,omitempty"`
	ChatID     string `json:"chat_id,omitempty"`
	SMTPHost   string `json:"smtp_host,omitempty"`
	SMTPPort   int    `json:"smtp_port,omitempty"`
	SMTPUser   string `json:"smtp_user,omitempty"`
	SMTPPass   string `json:"smtp_pass,omitempty"`
	ToEmail    string `json:"to_email,omitempty"`
}

// NotifyPolicy 通知策略
type NotifyPolicy struct {
	Event   string `json:"event"`   // hit, balance_low, error
	Level   string `json:"level"`   // high, normal, low
	Enabled bool   `json:"enabled"`
}

// WebhookConfig 飞书Webhook配置（兼容旧版）
type WebhookConfig struct {
	Enabled bool   `json:"enabled"`
	URL     string `json:"url"`
	Sign    string `json:"sign"`
}

// ========== 日志 ==========

// LogEntry 日志条目
type LogEntry struct {
	ID        string `json:"id"`
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"` // info, hit, error, warn
	Type      string `json:"type"` // check, notify, system, error
	Message   string `json:"message"`
	Details   string `json:"details,omitempty"`
}

// HitRecord 命中记录
type HitRecord struct {
	ID            string  `json:"id"`
	Timestamp     string  `json:"timestamp"`
	Service       string  `json:"service"`
	ServiceName   string  `json:"service_name"`
	Country       int     `json:"country"`
	CountryName   string  `json:"country_name"`
	Cost          float64 `json:"cost"`
	Count         int     `json:"count"`
	PhysicalCount int     `json:"physical_count"`
	MaxPrice      float64 `json:"max_price"`
	PhoneNumber   string  `json:"phone_number,omitempty"`
	SMSCode       string  `json:"sms_code,omitempty"`
}

// CheckRecord 检查记录（包括无命中的）
type CheckRecord struct {
	ID        string `json:"id"`
	Timestamp string `json:"timestamp"`
	Status    string `json:"status"` // success, error, no_hit
	HitCount  int    `json:"hit_count"`
	Error     string `json:"error,omitempty"`
	Duration  int64  `json:"duration_ms"`
}

// ========== 财务 ==========

// BalanceAlert 余额预警
type BalanceAlert struct {
	Enabled  bool    `json:"enabled"`
	Threshold float64 `json:"threshold"`
}

// CostRecord 消耗记录
type CostRecord struct {
	ID        string  `json:"id"`
	Timestamp string  `json:"timestamp"`
	Service   string  `json:"service"`
	Country   int     `json:"country"`
	Amount    float64 `json:"amount"`
	Type      string  `json:"type"` // purchase, refund
}

// ========== 系统 ==========

// ProxyConfig 代理配置
type ProxyConfig struct {
	Enabled bool   `json:"enabled"`
	URL     string `json:"url"`
}

// User 用户
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
	Role     string `json:"role"` // admin, viewer
	APIKey   string `json:"api_key,omitempty"`
}

// ========== API 响应 ==========

// MonitorStatus 监控状态
type MonitorStatus struct {
	Running       bool    `json:"running"`
	Balance       float64 `json:"balance"`
	CountryCount  int     `json:"country_count"`
	LastCheck     string  `json:"last_check"`
	TotalHits     int     `json:"total_hits"`
	TotalChecks   int     `json:"total_checks"`
	ErrorCount    int     `json:"error_count"`
	WatchCount    int     `json:"watch_count"`
	UptimeSeconds int64   `json:"uptime_seconds"`
}

// DashboardData 仪表盘数据
type DashboardData struct {
	Status      MonitorStatus `json:"status"`
	RecentHits  []HitRecord   `json:"recent_hits"`
	RecentLogs  []LogEntry    `json:"recent_logs"`
	Alerts      []LogEntry    `json:"alerts"`
}

// FinanceData 财务数据
type FinanceData struct {
	Balance      float64      `json:"balance"`
	TotalCost    float64      `json:"total_cost"`
	DailyCosts   []CostRecord `json:"daily_costs"`
	ServiceCosts map[string]float64 `json:"service_costs"`
}

// ========== 辅助函数 ==========

// GenerateID 生成简单ID
func GenerateID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(6)
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
		time.Sleep(1)
	}
	return string(b)
}
