package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// CountryInfo 国家信息
type CountryInfo struct {
	ID      int    `json:"id"`
	CHN     string `json:"chn"`
	ENG     string `json:"eng"`
	Visible int    `json:"visible"`
}

// PriceData 价格数据
type PriceData struct {
	Cost          float64 `json:"cost"`
	Count         int     `json:"count"`
	PhysicalCount int     `json:"physicalCount"`
}

// Monitor 监控器
type Monitor struct {
	mu          sync.RWMutex
	config      *ConfigManager
	storage     *Storage
	notifier    *Notifier
	client      *http.Client
	countries   map[int]string
	services    map[string]string
	running     bool
	balance     float64
	lastCheck   string
	startTime   time.Time
	stopCh      chan struct{}
	eventCh     chan Event
	subscribers []chan Event
	subMu       sync.Mutex
}

// Event SSE事件
type Event struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// NewMonitor 创建监控器
func NewMonitor(config *ConfigManager, storage *Storage) *Monitor {
	return &Monitor{
		config:    config,
		storage:   storage,
		notifier:  NewNotifier(),
		client:    &http.Client{Timeout: 30 * time.Second},
		countries: make(map[int]string),
		services:  make(map[string]string),
		stopCh:    make(chan struct{}),
		eventCh:   make(chan Event, 100),
		startTime: time.Now(),
	}
}

// Init 初始化
func (m *Monitor) Init() {
	m.fetchCountries()
	m.fetchServices()
	m.storage.AddLog("info", "system", "系统启动", fmt.Sprintf("加载 %d 个国家", len(m.countries)))
}

// fetchCountries 获取国家列表
func (m *Monitor) fetchCountries() {
	cfg := m.config.Get()
	url := fmt.Sprintf("%s?action=getCountries&api_key=%s", cfg.BaseURL, cfg.APIKey)

	resp, err := m.client.Get(url)
	if err != nil {
		m.storage.AddLog("error", "system", "获取国家列表失败", err.Error())
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var countries map[string]CountryInfo
	if err := json.Unmarshal(body, &countries); err != nil {
		return
	}

	m.mu.Lock()
	for _, info := range countries {
		if info.Visible == 1 {
			name := info.CHN
			if name == "" {
				name = info.ENG
			}
			m.countries[info.ID] = name
		}
	}
	m.mu.Unlock()
}

// fetchServices 获取服务列表
func (m *Monitor) fetchServices() {
	cfg := m.config.Get()
	url := fmt.Sprintf("%s?action=getServicesList&api_key=%s&lang=cn", cfg.BaseURL, cfg.APIKey)

	resp, err := m.client.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result struct {
		Services []struct {
			Code string `json:"code"`
			Name string `json:"name"`
		} `json:"services"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return
	}

	m.mu.Lock()
	for _, s := range result.Services {
		m.services[s.Code] = s.Name
	}
	m.mu.Unlock()
}

// GetCountryCount 获取国家数量
func (m *Monitor) GetCountryCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.countries)
}

// GetCountryName 获取国家名称
func (m *Monitor) GetCountryName(id int) string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if name, ok := m.countries[id]; ok {
		return name
	}
	return fmt.Sprintf("国家#%d", id)
}

// GetServices 获取服务列表
func (m *Monitor) GetServices() map[string]string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	result := make(map[string]string)
	for k, v := range m.services {
		result[k] = v
	}
	return result
}

// FetchBalance 获取余额
func (m *Monitor) FetchBalance() (float64, error) {
	cfg := m.config.Get()
	url := fmt.Sprintf("%s?action=getBalance&api_key=%s", cfg.BaseURL, cfg.APIKey)

	resp, err := m.client.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	text := string(body)

	var balance float64
	if _, err := fmt.Sscanf(text, "ACCESS_BALANCE:%f", &balance); err != nil {
		return 0, fmt.Errorf("unexpected response: %s", text)
	}

	m.mu.Lock()
	m.balance = balance
	m.mu.Unlock()

	// 检查余额告警
	cfgAlert := m.config.Get().BalanceAlert
	if cfgAlert.Enabled && balance < cfgAlert.Threshold {
		m.notifier.SendBalanceAlert(m.config.Get().Channels, balance)
	}

	return balance, nil
}

// FetchPrices 获取价格数据
func (m *Monitor) FetchPrices() (map[string]map[string]PriceData, error) {
	cfg := m.config.Get()
	url := fmt.Sprintf("%s?action=getPrices&api_key=%s", cfg.BaseURL, cfg.APIKey)

	resp, err := m.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var prices map[string]map[string]PriceData
	if err := json.Unmarshal(body, &prices); err != nil {
		return nil, err
	}

	return prices, nil
}

// RunCheck 执行一次检查
func (m *Monitor) RunCheck() ([]HitRecord, error) {
	start := time.Now()

	// 获取余额
	balance, err := m.FetchBalance()
	if err != nil {
		m.storage.AddLog("error", "check", "获取余额失败", err.Error())
		m.storage.AddCheck(CheckRecord{Status: "error", Error: err.Error(), Duration: time.Since(start).Milliseconds()})
		return nil, err
	}

	// 获取价格
	prices, err := m.FetchPrices()
	if err != nil {
		m.storage.AddLog("error", "check", "获取价格失败", err.Error())
		m.storage.AddCheck(CheckRecord{Status: "error", Error: err.Error(), Duration: time.Since(start).Milliseconds()})
		return nil, err
	}

	// 检查命中
	cfg := m.config.Get()
	var hits []HitRecord
	checkedAt := time.Now().Format("2006-01-02 15:04:05")

	for _, item := range cfg.Watchlist {
		if !item.Enabled {
			continue
		}

		countryStr := strconv.Itoa(item.Country)
		countryData, ok := prices[countryStr]
		if !ok {
			continue
		}

		serviceData, ok := countryData[item.Service]
		if !ok {
			continue
		}

		countryName := item.CountryName
		if countryName == "" {
			countryName = m.GetCountryName(item.Country)
		}

		hit := HitRecord{
			Timestamp:     checkedAt,
			Service:       item.Service,
			ServiceName:   item.ServiceName,
			Country:       item.Country,
			CountryName:   countryName,
			Cost:          serviceData.Cost,
			Count:         serviceData.Count,
			PhysicalCount: serviceData.PhysicalCount,
			MaxPrice:      item.MaxPrice,
		}

		isHit := serviceData.Cost <= item.MaxPrice && serviceData.Count > 0

		if isHit {
			hits = append(hits, hit)

			// 保存命中记录
			m.storage.AddHit(hit)

			// 记录消耗
			m.storage.AddCost(CostRecord{
				Service: item.ServiceName,
				Country: item.Country,
				Amount:  serviceData.Cost,
				Type:    "purchase",
			})
		}

		// 发送通知（命中或未命中都通知）
		go m.notifier.SendCheckNotification(cfg.Channels, hit, isHit)
	}

	// 记录日志
	if len(hits) > 0 {
		m.storage.AddLog("hit", "check", fmt.Sprintf("命中 %d 个号码", len(hits)), "")
	} else {
		m.storage.AddLog("info", "check", "本轮检查无命中", fmt.Sprintf("余额: $%.4f", balance))
	}

	// 记录检查记录
	m.storage.AddCheck(CheckRecord{
		Status:   "success",
		HitCount: len(hits),
		Duration: time.Since(start).Milliseconds(),
	})

	// 更新状态
	m.mu.Lock()
	m.lastCheck = time.Now().Format("2006-01-02 15:04:05")
	m.mu.Unlock()

	// 发送SSE事件
	m.sendEvent(Event{
		Type: "check",
		Data: map[string]interface{}{
			"balance":   balance,
			"hits":      hits,
			"time":      m.lastCheck,
			"hit_count": len(hits),
		},
	})

	return hits, nil
}

// GetStatus 获取监控状态
func (m *Monitor) GetStatus() MonitorStatus {
	m.mu.RLock()
	defer m.mu.RUnlock()
	cfg := m.config.Get()

	totalHits, totalChecks, errorCount := m.storage.GetStats()

	return MonitorStatus{
		Running:       m.running,
		Balance:       m.balance,
		CountryCount:  len(m.countries),
		LastCheck:     m.lastCheck,
		TotalHits:     totalHits,
		TotalChecks:   totalChecks,
		ErrorCount:    errorCount,
		WatchCount:    len(cfg.Watchlist),
		UptimeSeconds: int64(time.Since(m.startTime).Seconds()),
	}
}

// Start 启动监控
func (m *Monitor) Start() {
	m.mu.Lock()
	if m.running {
		m.mu.Unlock()
		return
	}
	m.running = true
	m.stopCh = make(chan struct{})
	m.mu.Unlock()

	m.storage.AddLog("info", "system", "监控已启动", "")
	m.sendEvent(Event{Type: "status", Data: "started"})

	go func() {
		// 立即执行一次
		if _, err := m.RunCheck(); err != nil {
			m.storage.AddLog("error", "system", "检查失败", err.Error())
		}

		interval := m.config.Get().IntervalSeconds
		if interval < 1 {
			interval = 60
		}
		ticker := time.NewTicker(time.Duration(interval) * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-m.stopCh:
				m.mu.Lock()
				m.running = false
				m.mu.Unlock()
				m.storage.AddLog("info", "system", "监控已停止", "")
				m.sendEvent(Event{Type: "status", Data: "stopped"})
				return
			case <-ticker.C:
				if _, err := m.RunCheck(); err != nil {
					m.storage.AddLog("error", "system", "检查失败", err.Error())
				}
			}
		}
	}()
}

// Stop 停止监控
func (m *Monitor) Stop() {
	m.mu.RLock()
	if !m.running {
		m.mu.RUnlock()
		return
	}
	m.mu.RUnlock()
	close(m.stopCh)
}

// IsRunning 是否运行中
func (m *Monitor) IsRunning() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.running
}

// Subscribe 订阅事件
func (m *Monitor) Subscribe() chan Event {
	ch := make(chan Event, 100)
	m.subMu.Lock()
	m.subscribers = append(m.subscribers, ch)
	m.subMu.Unlock()
	return ch
}

// Unsubscribe 取消订阅
func (m *Monitor) Unsubscribe(ch chan Event) {
	m.subMu.Lock()
	defer m.subMu.Unlock()
	for i, sub := range m.subscribers {
		if sub == ch {
			m.subscribers = append(m.subscribers[:i], m.subscribers[i+1:]...)
			break
		}
	}
	close(ch)
}

// sendEvent 发送事件
func (m *Monitor) sendEvent(event Event) {
	m.subMu.Lock()
	defer m.subMu.Unlock()
	for _, ch := range m.subscribers {
		select {
		case ch <- event:
		default:
		}
	}
}
