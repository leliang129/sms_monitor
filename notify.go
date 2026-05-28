package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/smtp"
	"strconv"
	"strings"
	"time"
)

// Notifier 通知管理器
type Notifier struct {
	client *http.Client
}

// NewNotifier 创建通知管理器
func NewNotifier() *Notifier {
	return &Notifier{
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

// SendHitNotification 发送命中通知
func (n *Notifier) SendHitNotification(channels []Channel, hit HitRecord) {
	for _, ch := range channels {
		if !ch.Enabled {
			continue
		}

		// 检查策略
		if !n.matchPolicy(ch, "hit", "normal") {
			continue
		}

		content := fmt.Sprintf(
			"🎯 %s 号码命中\n国家：%s\n💰 当前价格：$%.4f\n📦 可用数量：%d\n🎚 目标价格：≤ $%.4f\n✅ 可用号码数：%d（物理卡 %d）\n🕒 检查时间：%s",
			hit.ServiceName, hit.CountryName, hit.Cost,
			hit.Count, hit.MaxPrice, hit.Count, hit.PhysicalCount, hit.Timestamp,
		)

		go n.send(ch, content)
	}
}

// SendCheckNotification 发送检查结果通知（命中或未命中）
func (n *Notifier) SendCheckNotification(channels []Channel, hit HitRecord, isHit bool) {
	for _, ch := range channels {
		if !ch.Enabled {
			continue
		}

		// 检查策略
		if !n.matchPolicy(ch, "hit", "normal") {
			continue
		}

		var content string
		if isHit {
			content = fmt.Sprintf(
				"🎯 %s 号码命中\n国家：%s\n💰 当前价格：$%.4f\n📦 可用数量：%d\n🎚 目标价格：≤ $%.4f\n✅ 可用号码数：%d（物理卡 %d）\n🕒 检查时间：%s",
				hit.ServiceName, hit.CountryName, hit.Cost,
				hit.Count, hit.MaxPrice, hit.Count, hit.PhysicalCount, hit.Timestamp,
			)
		} else {
			content = fmt.Sprintf(
				"❌ %s 未命中\n国家：%s\n💰 当前价格：$%.4f\n📦 可用数量：%d\n🎚 目标价格：≤ $%.4f\n🕒 检查时间：%s",
				hit.ServiceName, hit.CountryName, hit.Cost,
				hit.Count, hit.MaxPrice, hit.Timestamp,
			)
		}

		go n.send(ch, content)
	}
}

// SendAlertNotification 发送告警通知
func (n *Notifier) SendAlertNotification(channels []Channel, level, message string) {
	for _, ch := range channels {
		if !ch.Enabled {
			continue
		}
		if !n.matchPolicy(ch, "alert", level) {
			continue
		}
		go n.send(ch, "⚠️ 系统告警\n"+message)
	}
}

// SendBalanceAlert 发送余额告警
func (n *Notifier) SendBalanceAlert(channels []Channel, balance float64) {
	msg := fmt.Sprintf("💳 余额告警\n当前余额：$%.4f\n请及时充值！", balance)
	for _, ch := range channels {
		if !ch.Enabled {
			continue
		}
		if !n.matchPolicy(ch, "balance_low", "high") {
			continue
		}
		go n.send(ch, msg)
	}
}

// matchPolicy 匹配通知策略
func (n *Notifier) matchPolicy(ch Channel, event, level string) bool {
	if len(ch.Policies) == 0 {
		return true // 无策略限制，默认允许
	}
	for _, p := range ch.Policies {
		if p.Event == event && p.Enabled {
			if p.Level == "" || p.Level == level {
				return true
			}
		}
	}
	return false
}

// send 发送通知
func (n *Notifier) send(ch Channel, content string) {
	var err error
	switch ch.Type {
	case "feishu":
		err = n.sendFeishu(ch.Config, content)
	case "dingtalk":
		err = n.sendDingTalk(ch.Config, content)
	case "wechat":
		err = n.sendWechat(ch.Config, content)
	case "telegram":
		err = n.sendTelegram(ch.Config, content)
	case "email":
		err = n.sendEmail(ch.Config, content)
	}
	if err != nil {
		fmt.Printf("通知发送失败 [%s]: %v\n", ch.Name, err)
	}
}

// sendFeishu 发送飞书
func (n *Notifier) sendFeishu(cfg ChannelConfig, content string) error {
	body := map[string]interface{}{
		"msg_type": "text",
		"content":  map[string]string{"text": content},
	}

	if cfg.Sign != "" {
		timestamp := time.Now().Unix()
		body["timestamp"] = strconv.FormatInt(timestamp, 10)
		body["sign"] = GenerateSign(cfg.Sign, timestamp)
	}

	return n.postJSON(cfg.WebhookURL, body)
}

// sendDingTalk 发送钉钉
func (n *Notifier) sendDingTalk(cfg ChannelConfig, content string) error {
	body := map[string]interface{}{
		"msgtype": "text",
		"text":    map[string]string{"content": content},
	}
	return n.postJSON(cfg.WebhookURL, body)
}

// sendWechat 发送企业微信
func (n *Notifier) sendWechat(cfg ChannelConfig, content string) error {
	body := map[string]interface{}{
		"msgtype": "text",
		"text":    map[string]string{"content": content},
	}
	return n.postJSON(cfg.WebhookURL, body)
}

// sendTelegram 发送Telegram
func (n *Notifier) sendTelegram(cfg ChannelConfig, content string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", cfg.BotToken)
	body := map[string]interface{}{
		"chat_id":    cfg.ChatID,
		"text":       content,
		"parse_mode": "HTML",
	}
	return n.postJSON(url, body)
}

// sendEmail 发送邮件
func (n *Notifier) sendEmail(cfg ChannelConfig, content string) error {
	if cfg.SMTPHost == "" {
		return fmt.Errorf("SMTP未配置")
	}

	addr := fmt.Sprintf("%s:%d", cfg.SMTPHost, cfg.SMTPPort)
	auth := smtp.PlainAuth("", cfg.SMTPUser, cfg.SMTPPass, cfg.SMTPHost)

	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: HeroSMS Monitor 通知\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s",
		cfg.SMTPUser, cfg.ToEmail, content)

	return smtp.SendMail(addr, auth, cfg.SMTPUser, strings.Split(cfg.ToEmail, ","), []byte(msg))
}

// postJSON POST JSON请求
func (n *Notifier) postJSON(url string, body interface{}) error {
	data, err := json.Marshal(body)
	if err != nil {
		return err
	}

	resp, err := n.client.Post(url, "application/json", bytes.NewReader(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}

// GenerateSign 生成签名
func GenerateSign(secret string, timestamp int64) string {
	stringToSign := fmt.Sprintf("%d\n%s", timestamp, secret)
	h := hmac.New(sha256.New, []byte(stringToSign))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
