# HeroSMS Monitor

HeroSMS 号码价格监控系统 - Go 版本

[![Go Report Card](https://goreportcard.com/badge/github.com/leliang129/sms_monitor)](https://goreportcard.com/report/github.com/leliang129/sms_monitor)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

## 功能特性

- 实时监控指定服务/国家的号码价格
- 价格低于目标时自动告警
- 多渠道通知：飞书、钉钉、企业微信、Telegram、邮件
- 精美 Web UI 管理面板
- 历史日志查看与 CSV 导出
- 财务看板与消耗明细
- 系统配置在线管理
- SSE 实时状态推送

## 快速开始

### 方式一：Docker 部署（推荐）

```bash
# 使用 docker-compose
docker compose up -d

# 或者直接运行
docker run -d \
  --name sms-monitor \
  -p 8080:8080 \
  -v ./data:/app/data \
  -v ./config.json:/app/config.json \
  ghcr.io/leliang129/sms_monitor:latest
```

### 方式二：下载预编译版本

前往 [Releases](https://github.com/leliang129/sms_monitor/releases) 下载对应平台的可执行文件。

支持平台：Linux (amd64/arm64), macOS (amd64/arm64), Windows (amd64)

#### 1. 下载并解压

```bash
# Linux amd64
wget https://github.com/leliang129/sms_monitor/releases/latest/download/sms_monitor_linux_amd64
chmod +x sms_monitor_linux_amd64

# Linux arm64
wget https://github.com/leliang129/sms_monitor/releases/latest/download/sms_monitor_linux_arm64
chmod +x sms_monitor_linux_arm64

# macOS arm64 (Apple Silicon)
wget https://github.com/leliang129/sms_monitor/releases/latest/download/sms_monitor_darwin_arm64
chmod +x sms_monitor_darwin_arm64
```

#### 2. 创建配置文件

```bash
cp config.example.json config.json
# 编辑 config.json，填入你的 API Key
```

#### 3. 运行

```bash
# 前台运行
./sms_monitor_linux_amd64

# 后台运行
nohup ./sms_monitor_linux_amd64 > sms_monitor.log 2>&1 &
```

#### 4. 设置开机自启（Linux systemd）

创建服务文件 `/etc/systemd/system/sms-monitor.service`：

```ini
[Unit]
Description=HeroSMS Monitor
After=network.target

[Service]
Type=simple
WorkingDirectory=/opt/sms_monitor
ExecStart=/opt/sms_monitor/sms_monitor_linux_amd64
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

启用并启动服务：

```bash
sudo systemctl daemon-reload
sudo systemctl enable sms-monitor
sudo systemctl start sms-monitor

# 查看状态
sudo systemctl status sms-monitor

# 查看日志
sudo journalctl -u sms-monitor -f
```

启动后访问 http://localhost:8080

### 方式三：从源码编译

```bash
git clone https://github.com/leliang129/sms_monitor.git
cd sms_monitor
go build -o sms_monitor .
./sms_monitor
```

启动后访问 http://localhost:8080

## 配置说明

首次运行会自动生成 `config.json`，也可以通过 Web UI 在线配置。

```json
{
  "api_key": "YOUR_HEROSMS_API_KEY",
  "base_url": "https://hero-sms.com/stubs/handler_api.php",
  "interval_seconds": 600,
  "proxy": {
    "enabled": false,
    "url": ""
  }
}
```

| 参数 | 说明 | 默认值 |
|------|------|--------|
| `api_key` | HeroSMS API 密钥 | 必填 |
| `base_url` | API 接口地址 | hero-sms.com |
| `interval_seconds` | 检查间隔(秒) | 600 (10分钟) |
| `proxy.url` | 代理地址 | 空(直连) |

### 通知渠道配置

支持以下通知渠道：

- **飞书**: Webhook URL + 签名
- **钉钉**: Webhook URL + 签名
- **企业微信**: Webhook URL
- **Telegram**: Bot Token + Chat ID
- **邮件**: SMTP 服务器配置

## 镜像版本

| 标签 | 说明 |
|------|------|
| `latest` | 最新稳定版 |
| `v1.x.x` | 指定版本 |
| `v1.x` | 指定主次版本 |

镜像支持架构：`linux/amd64`, `linux/arm64`

## 项目结构

```
sms_monitor/
├── main.go              # 程序入口
├── config.go            # 配置管理
├── models.go            # 数据模型
├── monitor.go           # 核心监控逻辑
├── handler.go           # HTTP API 处理
├── storage.go           # 数据存储
├── notify.go            # 通知发送
├── Dockerfile           # Docker 构建文件
├── docker-compose.yml   # Docker Compose 配置
├── web/
│   ├── index.html       # 前端页面
│   └── style.css        # 样式文件
├── .github/workflows/   # GitHub Actions
├── data/                # 运行数据(自动创建)
└── config.json          # 配置文件(自动创建)
```

## API 接口

| 接口 | 方法 | 说明 |
|------|------|------|
| `/api/dashboard` | GET | 仪表盘数据 |
| `/api/status` | GET | 监控状态 |
| `/api/config` | GET/PUT | 配置管理 |
| `/api/monitor/start` | POST | 启动监控 |
| `/api/monitor/stop` | POST | 停止监控 |
| `/api/check` | POST | 手动检查 |
| `/api/watchlist` | GET | 监控列表 |
| `/api/logs` | GET | 日志查询 |
| `/api/hits` | GET | 命中记录 |
| `/api/channels` | GET | 通知渠道 |
| `/api/finance` | GET | 财务数据 |
| `/api/events` | GET | SSE 事件流 |

## 开发

```bash
# 运行
go run .

# 构建
go build -o sms_monitor .

# 交叉编译
GOOS=linux GOARCH=amd64 go build -o sms_monitor_linux_amd64 .
GOOS=linux GOARCH=arm64 go build -o sms_monitor_linux_arm64 .
GOOS=windows GOARCH=amd64 go build -o sms_monitor_windows_amd64.exe .
GOOS=darwin GOARCH=arm64 go build -o sms_monitor_darwin_arm64 .

# Docker 构建
docker build -t sms_monitor .
```

## CI/CD

推送到 `main` 分支时自动运行 CI 检查。

打 tag 时自动构建并发布：
```bash
git tag v1.0.0
git push origin v1.0.0
```

自动发布内容：
- 多平台二进制文件 (Linux/macOS/Windows, amd64/arm64)
- Docker 镜像 (ghcr.io/leliang129/sms_monitor)
- GitHub Release

## License

[MIT](LICENSE)
