package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	configPath := "config.json"
	dataDir := "data"
	addr := ":8080"

	if len(os.Args) > 1 {
		configPath = os.Args[1]
	}
	if len(os.Args) > 2 {
		dataDir = os.Args[2]
	}

	// 加载配置
	config, err := NewConfigManager(configPath)
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 初始化存储
	storage := NewStorage(dataDir)

	// 创建监控器
	monitor := NewMonitor(config, storage)
	monitor.Init()

	// 创建处理器
	handler := NewHandler(monitor, config, storage)

	// 注册路由
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	// 静态文件 - 从磁盘读取，支持热更新
	mux.Handle("/", http.FileServer(http.Dir("web")))

	// 启动信息
	cfg := config.Get()
	fmt.Println("╔══════════════════════════════════════════╗")
	fmt.Println("║        HeroSMS Monitor (Go)              ║")
	fmt.Println("╠══════════════════════════════════════════╣")
	fmt.Printf("║  Web UI:    http://localhost%-13s║\n", addr)
	fmt.Printf("║  国家库:    %-3d 个国家%-17s║\n", monitor.GetCountryCount(), "")
	fmt.Printf("║  监控项:    %-3d 个%-22s║\n", len(cfg.Watchlist), "")
	fmt.Printf("║  通知渠道:  %-3d 个%-22s║\n", len(cfg.Channels), "")
	fmt.Println("╚══════════════════════════════════════════╝")
	fmt.Println()

	storage.AddLog("info", "system", "系统启动", fmt.Sprintf("监听 %s", addr))

	log.Printf("服务器启动在 %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
