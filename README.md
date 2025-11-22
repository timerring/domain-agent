# Domain Agent

AI-powered domain discovery platform

## 项目简介

Domain Agent 是一个基于 AI 的域名查询和推荐系统，通过对话式交互帮助用户快速找到合适的域名。

### 为什么需要 Domain Agent？

- 短域名价值高
- 与品牌强绑定

**痛点：**
- 好的域名大多已被注册，不清楚还有哪些域名可用
- 查询多个扩展名很麻烦，手动查询效率低
- 想到有记忆点的域名不容易，多数人没有创意
- 市面上的 AI 查询工具不够智能

### 核心功能

- AI 理解你的需求(已有初步选定域名 or 仅有意向)
- 智能生成创意域名
   1. 越短越好
   2. 有意思优先
   3. 可以考虑母语文字
   4. ...
- 基于 Golang 并发批量查询，并且覆盖尽可能多的域名，.com、.cn、.ai、.io 等。
- 多种方式验证可用性
   1. WHOIS
   2. DNS
   3. SSL 证书
- 反馈验证结果 

### 技术栈

**前端：**
- React 18 + TypeScript
- TailwindCSS + @tailwindcss/typography
- Google Fonts (Inter + JetBrains Mono)
- Vite
- Modern minimalist design

**后端：**
- Go 1.22+
- Gin Web Framework
- Vibecoding API (GPT-4)
- Gorilla WebSocket
- 集成 domain-scanner

## 快速开始

修改 `backend/.env.example` 文件中的配置，配置模型 API key，然后复制到 `backend/.env`。

### 一键启动

```bash
./start.sh
```

### 分别启动

**后端：**
```bash
cd backend
go run cmd/server/main.go
```

**前端：**
```bash
cd frontend
npm install
npm run dev
```

### 访问应用

- 前端界面：http://localhost:5173
- 后端 API：http://localhost:8080
- 健康检查：http://localhost:8080/health

## 项目结构

```
domain-agent/
├── backend/              # Go 后端服务
│   ├── cmd/server/      # 主程序入口
│   └── internal/        # 内部包
│       ├── agent/       # Agent 逻辑引擎
│       ├── api/         # HTTP handlers
│       ├── scanner/     # 域名扫描集成
│       └── types/       # 类型定义
│
└── frontend/            # React 前端应用
    └── src/
        ├── components/  # UI 组件
        ├── services/    # API 服务
        └── App.tsx      # 主应用
```

## 贡献

欢迎提交 Issue 和 Pull Request！

