package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

// Handler HTTP处理器
type Handler struct {
	monitor *Monitor
	config  *ConfigManager
	storage *Storage
}

// NewHandler 创建处理器
func NewHandler(monitor *Monitor, config *ConfigManager, storage *Storage) *Handler {
	return &Handler{monitor: monitor, config: config, storage: storage}
}

// RegisterRoutes 注册路由
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	// 仪表盘
	mux.HandleFunc("/api/dashboard", h.handleDashboard)
	mux.HandleFunc("/api/status", h.handleStatus)

	// 监控设置
	mux.HandleFunc("/api/config", h.handleConfig)
	mux.HandleFunc("/api/watchlist", h.handleWatchlist)
	mux.HandleFunc("/api/watchlist/add", h.handleWatchlistAdd)
	mux.HandleFunc("/api/watchlist/delete", h.handleWatchlistDelete)
	mux.HandleFunc("/api/watchlist/update", h.handleWatchlistUpdate)
	mux.HandleFunc("/api/services", h.handleServices)
	mux.HandleFunc("/api/countries", h.handleCountries)

	// 监控控制
	mux.HandleFunc("/api/monitor/start", h.handleStart)
	mux.HandleFunc("/api/monitor/stop", h.handleStop)
	mux.HandleFunc("/api/check", h.handleCheck)

	// 日志
	mux.HandleFunc("/api/logs", h.handleLogs)
	mux.HandleFunc("/api/hits", h.handleHits)
	mux.HandleFunc("/api/checks", h.handleChecks)
	mux.HandleFunc("/api/logs/export", h.handleExportLogs)

	// 通知渠道
	mux.HandleFunc("/api/channels", h.handleChannels)
	mux.HandleFunc("/api/channels/add", h.handleChannelAdd)
	mux.HandleFunc("/api/channels/update", h.handleChannelUpdate)
	mux.HandleFunc("/api/channels/delete", h.handleChannelDelete)

	// 财务
	mux.HandleFunc("/api/finance", h.handleFinance)

	// 系统配置
	mux.HandleFunc("/api/system/proxy", h.handleProxy)
	mux.HandleFunc("/api/system/users", h.handleUsers)
	mux.HandleFunc("/api/system/export", h.handleExport)

	// SSE
	mux.HandleFunc("/api/events", h.handleEvents)
}

// ========== 辅助函数 ==========

func jsonResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func errorResponse(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

// ========== 仪表盘 ==========

func (h *Handler) handleDashboard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorResponse(w, 405, "Method not allowed")
		return
	}

	data := DashboardData{
		Status:     h.monitor.GetStatus(),
		RecentHits: h.storage.GetHits(50),
		RecentLogs: h.storage.GetLogs("", 100),
		Alerts:     h.storage.GetLogs("error", 20),
	}
	jsonResponse(w, data)
}

func (h *Handler) handleStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorResponse(w, 405, "Method not allowed")
		return
	}
	jsonResponse(w, h.monitor.GetStatus())
}

// ========== 监控设置 ==========

func (h *Handler) handleConfig(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		jsonResponse(w, h.config.Get())
	case http.MethodPut:
		var cfg Config
		if err := json.NewDecoder(r.Body).Decode(&cfg); err != nil {
			errorResponse(w, 400, "Invalid JSON")
			return
		}
		if err := h.config.Update(cfg); err != nil {
			errorResponse(w, 500, err.Error())
			return
		}
		h.storage.AddLog("info", "system", "配置已更新", "")
		jsonResponse(w, map[string]string{"status": "ok"})
	default:
		errorResponse(w, 405, "Method not allowed")
	}
}

func (h *Handler) handleWatchlist(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorResponse(w, 405, "Method not allowed")
		return
	}
	jsonResponse(w, h.config.Get().Watchlist)
}

func (h *Handler) handleWatchlistAdd(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errorResponse(w, 405, "Method not allowed")
		return
	}

	var item WatchItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		errorResponse(w, 400, "Invalid JSON")
		return
	}

	cfg := h.config.Get()
	item.ID = GenerateID()
	item.Enabled = true
	cfg.Watchlist = append(cfg.Watchlist, item)

	if err := h.config.Update(cfg); err != nil {
		errorResponse(w, 500, err.Error())
		return
	}

	h.storage.AddLog("info", "system", "添加监控项", fmt.Sprintf("%s @ %s", item.ServiceName, item.CountryName))
	jsonResponse(w, item)
}

func (h *Handler) handleWatchlistDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errorResponse(w, 405, "Method not allowed")
		return
	}

	var req struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorResponse(w, 400, "Invalid JSON")
		return
	}

	cfg := h.config.Get()
	for i, item := range cfg.Watchlist {
		if item.ID == req.ID {
			cfg.Watchlist = append(cfg.Watchlist[:i], cfg.Watchlist[i+1:]...)
			break
		}
	}

	if err := h.config.Update(cfg); err != nil {
		errorResponse(w, 500, err.Error())
		return
	}

	jsonResponse(w, map[string]string{"status": "ok"})
}

func (h *Handler) handleWatchlistUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errorResponse(w, 405, "Method not allowed")
		return
	}

	var item WatchItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		errorResponse(w, 400, "Invalid JSON")
		return
	}

	cfg := h.config.Get()
	for i, w := range cfg.Watchlist {
		if w.ID == item.ID {
			cfg.Watchlist[i] = item
			break
		}
	}

	if err := h.config.Update(cfg); err != nil {
		errorResponse(w, 500, err.Error())
		return
	}

	jsonResponse(w, item)
}

func (h *Handler) handleServices(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorResponse(w, 405, "Method not allowed")
		return
	}
	jsonResponse(w, h.monitor.GetServices())
}

func (h *Handler) handleCountries(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorResponse(w, 405, "Method not allowed")
		return
	}

	cfg := h.config.Get()
	url := fmt.Sprintf("%s?action=getCountries&api_key=%s", cfg.BaseURL, cfg.APIKey)
	resp, err := h.monitor.client.Get(url)
	if err != nil {
		errorResponse(w, 500, err.Error())
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	io.Copy(w, resp.Body)
}

// ========== 监控控制 ==========

func (h *Handler) handleStart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errorResponse(w, 405, "Method not allowed")
		return
	}
	h.monitor.Start()
	jsonResponse(w, map[string]string{"status": "started"})
}

func (h *Handler) handleStop(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errorResponse(w, 405, "Method not allowed")
		return
	}
	h.monitor.Stop()
	jsonResponse(w, map[string]string{"status": "stopped"})
}

func (h *Handler) handleCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errorResponse(w, 405, "Method not allowed")
		return
	}

	hits, err := h.monitor.RunCheck()
	if err != nil {
		errorResponse(w, 500, err.Error())
		return
	}

	jsonResponse(w, map[string]interface{}{
		"hits":  hits,
		"count": len(hits),
	})
}

// ========== 日志 ==========

func (h *Handler) handleLogs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorResponse(w, 405, "Method not allowed")
		return
	}

	logType := r.URL.Query().Get("type")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 {
		limit = 100
	}

	jsonResponse(w, h.storage.GetLogs(logType, limit))
}

func (h *Handler) handleHits(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorResponse(w, 405, "Method not allowed")
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 {
		limit = 100
	}

	jsonResponse(w, h.storage.GetHits(limit))
}

func (h *Handler) handleChecks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorResponse(w, 405, "Method not allowed")
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 {
		limit = 100
	}

	jsonResponse(w, h.storage.GetChecks(limit))
}

func (h *Handler) handleExportLogs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorResponse(w, 405, "Method not allowed")
		return
	}

	filename, err := h.storage.ExportHitsCSV()
	if err != nil {
		errorResponse(w, 500, err.Error())
		return
	}

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	http.ServeFile(w, r, filename)
}

// ========== 通知渠道 ==========

func (h *Handler) handleChannels(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorResponse(w, 405, "Method not allowed")
		return
	}
	jsonResponse(w, h.config.Get().Channels)
}

func (h *Handler) handleChannelAdd(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errorResponse(w, 405, "Method not allowed")
		return
	}

	var ch Channel
	if err := json.NewDecoder(r.Body).Decode(&ch); err != nil {
		errorResponse(w, 400, "Invalid JSON")
		return
	}

	cfg := h.config.Get()
	ch.ID = GenerateID()
	cfg.Channels = append(cfg.Channels, ch)

	if err := h.config.Update(cfg); err != nil {
		errorResponse(w, 500, err.Error())
		return
	}

	h.storage.AddLog("info", "system", "添加通知渠道", ch.Name)
	jsonResponse(w, ch)
}

func (h *Handler) handleChannelUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errorResponse(w, 405, "Method not allowed")
		return
	}

	var ch Channel
	if err := json.NewDecoder(r.Body).Decode(&ch); err != nil {
		errorResponse(w, 400, "Invalid JSON")
		return
	}

	cfg := h.config.Get()
	for i, c := range cfg.Channels {
		if c.ID == ch.ID {
			cfg.Channels[i] = ch
			break
		}
	}

	if err := h.config.Update(cfg); err != nil {
		errorResponse(w, 500, err.Error())
		return
	}

	jsonResponse(w, ch)
}

func (h *Handler) handleChannelDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errorResponse(w, 405, "Method not allowed")
		return
	}

	var req struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorResponse(w, 400, "Invalid JSON")
		return
	}

	cfg := h.config.Get()
	for i, ch := range cfg.Channels {
		if ch.ID == req.ID {
			cfg.Channels = append(cfg.Channels[:i], cfg.Channels[i+1:]...)
			break
		}
	}

	if err := h.config.Update(cfg); err != nil {
		errorResponse(w, 500, err.Error())
		return
	}

	jsonResponse(w, map[string]string{"status": "ok"})
}

// ========== 财务 ==========

func (h *Handler) handleFinance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorResponse(w, 405, "Method not allowed")
		return
	}

	days, _ := strconv.Atoi(r.URL.Query().Get("days"))
	if days <= 0 {
		days = 7
	}

	data := FinanceData{
		Balance:      h.monitor.balance,
		TotalCost:    0,
		DailyCosts:   h.storage.GetCosts(days),
		ServiceCosts: h.storage.GetServiceCosts(days),
	}

	for _, c := range data.DailyCosts {
		data.TotalCost += c.Amount
	}

	jsonResponse(w, data)
}

// ========== 系统配置 ==========

func (h *Handler) handleProxy(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		jsonResponse(w, h.config.Get().Proxy)
	case http.MethodPut:
		var proxy ProxyConfig
		if err := json.NewDecoder(r.Body).Decode(&proxy); err != nil {
			errorResponse(w, 400, "Invalid JSON")
			return
		}
		cfg := h.config.Get()
		cfg.Proxy = proxy
		h.config.Update(cfg)
		jsonResponse(w, map[string]string{"status": "ok"})
	default:
		errorResponse(w, 405, "Method not allowed")
	}
}

func (h *Handler) handleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// 返回用户列表（隐藏密码）
		users := h.config.Get().Users
		safeUsers := make([]User, len(users))
		for i, u := range users {
			safeUsers[i] = User{
				ID:       u.ID,
				Username: u.Username,
				Role:     u.Role,
				APIKey:   u.APIKey,
			}
		}
		jsonResponse(w, safeUsers)
	default:
		errorResponse(w, 405, "Method not allowed")
	}
}

func (h *Handler) handleExport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorResponse(w, 405, "Method not allowed")
		return
	}

	format := r.URL.Query().Get("format")
	if format == "csv" {
		filename, err := h.storage.ExportHitsCSV()
		if err != nil {
			errorResponse(w, 500, err.Error())
			return
		}
		w.Header().Set("Content-Type", "text/csv")
		w.Header().Set("Content-Disposition", "attachment; filename="+filename)
		http.ServeFile(w, r, filename)
	} else {
		// 默认导出JSON
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Disposition", "attachment; filename=export.json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"config":  h.config.Get(),
			"hits":    h.storage.GetHits(0),
			"logs":    h.storage.GetLogs("", 0),
		})
	}
}

// ========== SSE ==========

func (h *Handler) handleEvents(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorResponse(w, 405, "Method not allowed")
		return
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		errorResponse(w, 500, "Streaming not supported")
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	eventCh := h.monitor.Subscribe()
	defer h.monitor.Unsubscribe(eventCh)

	// 发送初始状态
	status := h.monitor.GetStatus()
	data, _ := json.Marshal(status)
	fmt.Fprintf(w, "event: status\ndata: %s\n\n", data)
	flusher.Flush()

	for {
		select {
		case event := <-eventCh:
			data, _ := json.Marshal(event.Data)
			fmt.Fprintf(w, "event: %s\ndata: %s\n\n", event.Type, data)
			flusher.Flush()
		case <-r.Context().Done():
			return
		}
	}
}
