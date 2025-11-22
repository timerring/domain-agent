# Domain Agent Backend

Go 后端服务

## 快速开始

### 1. 安装依赖

```bash
go mod download
```

### 2. 配置环境变量

```bash
# 复制环境变量模板
cp .env.example .env

# 编辑 .env 文件，添加你的 Vibecoding API Key
# VIBECODING_API_KEY=your_vibecoding_api_key_here
```

### 3. 运行服务

```bash
go run cmd/server/main.go
```

服务将在 `http://localhost:8080` 启动

## AI 功能说明

### 智能域名生成

系统集成了 Vibecoding API (GPT-4)，提供以下 AI 功能：

1. **智能意图识别** - 自动识别用户是想查询域名还是需要创意建议
2. **创意域名生成** - 根据用户描述生成个性化的域名建议
3. **智能评分排序** - 基于多个维度评估域名价值

### 使用示例

```bash
# 测试 AI 域名生成
curl -X POST http://localhost:8080/api/agent/chat \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "test",
    "message": "我想要一个科技公司的域名，简短有记忆点"
  }'
```

### 环境变量配置

| 变量名 | 说明 | 必需 |
|--------|------|------|
| VIBECODING_API_KEY | Vibecoding API 密钥 | 是 |
| PORT | 服务端口 | 否 (默认 8080) |
| GIN_MODE | 运行模式 | 否 (默认 debug) |

**注意**: 如果没有配置 `VIBECODING_API_KEY`，系统会回退到基于规则的简单响应。

## API 文档

### Agent 相关

- `POST /api/agent/chat` - 发送对话消息
- `GET /api/agent/session/:id` - 获取会话信息
- `GET /api/agent/stream` - WebSocket 流式连接

### 域名相关

- `POST /api/domains/check` - 批量检查域名
- `POST /api/domains/suggest` - 生成域名建议

## 项目结构

```
backend/
├── cmd/
│   └── server/          # 主程序入口
├── internal/
│   ├── agent/           # Agent 逻辑
│   ├── api/             # HTTP handlers
│   ├── scanner/         # 域名扫描
│   └── types/           # 类型定义
└── go.mod
```
