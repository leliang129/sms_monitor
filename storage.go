package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Storage 数据存储
type Storage struct {
	mu       sync.RWMutex
	dataDir  string
	logs     []LogEntry
	hits     []HitRecord
	checks   []CheckRecord
	costs    []CostRecord
	maxLogs  int
}

// NewStorage 创建存储
func NewStorage(dataDir string) *Storage {
	s := &Storage{
		dataDir: dataDir,
		maxLogs: 10000,
	}
	os.MkdirAll(dataDir, 0755)
	s.loadAll()
	return s
}

// ========== 日志 ==========

// AddLog 添加日志
func (s *Storage) AddLog(level, logType, message, details string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry := LogEntry{
		ID:        GenerateID(),
		Timestamp: time.Now().Format("2006-01-02 15:04:05"),
		Level:     level,
		Type:      logType,
		Message:   message,
		Details:   details,
	}
	s.logs = append([]LogEntry{entry}, s.logs...)
	if len(s.logs) > s.maxLogs {
		s.logs = s.logs[:s.maxLogs]
	}
	go s.saveLogs()
}

// GetLogs 获取日志
func (s *Storage) GetLogs(logType string, limit int) []LogEntry {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []LogEntry
	for _, l := range s.logs {
		if logType == "" || l.Type == logType {
			result = append(result, l)
			if limit > 0 && len(result) >= limit {
				break
			}
		}
	}
	return result
}

// ========== 命中记录 ==========

// AddHit 添加命中
func (s *Storage) AddHit(hit HitRecord) {
	s.mu.Lock()
	defer s.mu.Unlock()

	hit.ID = GenerateID()
	hit.Timestamp = time.Now().Format("2006-01-02 15:04:05")
	s.hits = append([]HitRecord{hit}, s.hits...)
	if len(s.hits) > s.maxLogs {
		s.hits = s.hits[:s.maxLogs]
	}
	go s.saveHits()
}

// GetHits 获取命中记录
func (s *Storage) GetHits(limit int) []HitRecord {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if limit <= 0 || limit > len(s.hits) {
		limit = len(s.hits)
	}
	result := make([]HitRecord, limit)
	copy(result, s.hits[:limit])
	return result
}

// ========== 检查记录 ==========

// AddCheck 添加检查记录
func (s *Storage) AddCheck(record CheckRecord) {
	s.mu.Lock()
	defer s.mu.Unlock()

	record.ID = GenerateID()
	record.Timestamp = time.Now().Format("2006-01-02 15:04:05")
	s.checks = append([]CheckRecord{record}, s.checks...)
	if len(s.checks) > s.maxLogs {
		s.checks = s.checks[:s.maxLogs]
	}
	go s.saveChecks()
}

// GetChecks 获取检查记录
func (s *Storage) GetChecks(limit int) []CheckRecord {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if limit <= 0 || limit > len(s.checks) {
		limit = len(s.checks)
	}
	result := make([]CheckRecord, limit)
	copy(result, s.checks[:limit])
	return result
}

// ========== 消耗记录 ==========

// AddCost 添加消耗
func (s *Storage) AddCost(record CostRecord) {
	s.mu.Lock()
	defer s.mu.Unlock()

	record.ID = GenerateID()
	record.Timestamp = time.Now().Format("2006-01-02 15:04:05")
	s.costs = append([]CostRecord{record}, s.costs...)
	if len(s.costs) > s.maxLogs {
		s.costs = s.costs[:s.maxLogs]
	}
	go s.saveCosts()
}

// GetCosts 获取消耗记录
func (s *Storage) GetCosts(days int) []CostRecord {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if days <= 0 {
		return s.costs
	}

	cutoff := time.Now().AddDate(0, 0, -days).Format("2006-01-02")
	var result []CostRecord
	for _, c := range s.costs {
		if c.Timestamp >= cutoff {
			result = append(result, c)
		}
	}
	return result
}

// GetServiceCosts 按服务统计消耗
func (s *Storage) GetServiceCosts(days int) map[string]float64 {
	costs := s.GetCosts(days)
	result := make(map[string]float64)
	for _, c := range costs {
		result[c.Service] += c.Amount
	}
	return result
}

// ========== 导出 ==========

// ExportHitsCSV 导出命中记录为CSV
func (s *Storage) ExportHitsCSV() (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	filename := filepath.Join(s.dataDir, fmt.Sprintf("hits_%s.csv", time.Now().Format("20060102_150405")))
	f, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	// 写入 BOM
	f.Write([]byte{0xEF, 0xBB, 0xBF})

	w := csv.NewWriter(f)
	defer w.Flush()

	// 表头
	w.Write([]string{"时间", "服务", "服务代码", "国家", "国家ID", "价格", "数量", "物理卡", "目标价"})

	for _, h := range s.hits {
		w.Write([]string{
			h.Timestamp,
			h.ServiceName,
			h.Service,
			h.CountryName,
			fmt.Sprintf("%d", h.Country),
			fmt.Sprintf("%.4f", h.Cost),
			fmt.Sprintf("%d", h.Count),
			fmt.Sprintf("%d", h.PhysicalCount),
			fmt.Sprintf("%.4f", h.MaxPrice),
		})
	}

	return filename, nil
}

// ========== 持久化 ==========

func (s *Storage) dataFile(name string) string {
	return filepath.Join(s.dataDir, name+".json")
}

func (s *Storage) saveJSON(filename string, data interface{}) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(data)
}

func (s *Storage) loadJSON(filename string, v interface{}) error {
	f, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer f.Close()
	return json.NewDecoder(f).Decode(v)
}

func (s *Storage) saveLogs() {
	s.saveJSON(s.dataFile("logs"), s.logs)
}

func (s *Storage) saveHits() {
	s.saveJSON(s.dataFile("hits"), s.hits)
}

func (s *Storage) saveChecks() {
	s.saveJSON(s.dataFile("checks"), s.checks)
}

func (s *Storage) saveCosts() {
	s.saveJSON(s.dataFile("costs"), s.costs)
}

func (s *Storage) loadAll() {
	s.loadJSON(s.dataFile("logs"), &s.logs)
	s.loadJSON(s.dataFile("hits"), &s.hits)
	s.loadJSON(s.dataFile("checks"), &s.checks)
	s.loadJSON(s.dataFile("costs"), &s.costs)
}

// GetStats 统计信息
func (s *Storage) GetStats() (totalHits, totalChecks, errorCount int) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	totalHits = len(s.hits)
	totalChecks = len(s.checks)
	for _, c := range s.checks {
		if c.Status == "error" {
			errorCount++
		}
	}
	return
}
