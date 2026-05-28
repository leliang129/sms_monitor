package main

import (
	"encoding/json"
	"os"
	"sync"
)

// ConfigManager 配置管理器
type ConfigManager struct {
	mu       sync.RWMutex
	config   Config
	filePath string
}

// NewConfigManager 创建配置管理器
func NewConfigManager(filePath string) (*ConfigManager, error) {
	cm := &ConfigManager{filePath: filePath}
	if err := cm.Load(); err != nil {
		return nil, err
	}
	return cm, nil
}

// Load 加载配置
func (cm *ConfigManager) Load() error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	data, err := os.ReadFile(cm.filePath)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &cm.config)
}

// Save 保存配置
func (cm *ConfigManager) Save() error {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	data, err := json.MarshalIndent(cm.config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(cm.filePath, data, 0644)
}

// Get 获取配置副本
func (cm *ConfigManager) Get() Config {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.config
}

// Update 更新配置
func (cm *ConfigManager) Update(cfg Config) error {
	cm.mu.Lock()
	cm.config = cfg
	cm.mu.Unlock()
	return cm.Save()
}

// GetAPIKey 获取API Key
func (cm *ConfigManager) GetAPIKey() string {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.config.APIKey
}

// GetBaseURL 获取Base URL
func (cm *ConfigManager) GetBaseURL() string {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.config.BaseURL
}
