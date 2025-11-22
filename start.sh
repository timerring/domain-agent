#!/bin/bash

echo "Starting Domain Agent..."

# 加载环境变量
if [ -f "backend/.env" ]; then
    export $(grep -v '^#' backend/.env | xargs)
    echo "✅ Environment variables loaded from backend/.env"
fi

# 启动后端
echo "Starting backend..."
cd backend
go run cmd/server/main.go &
BACKEND_PID=$!

# 等待后端启动
sleep 3

# 启动前端
echo "Starting frontend..."
cd ../frontend
npm run dev &
FRONTEND_PID=$!

echo ""
echo "Domain Agent is running!"
echo "   Backend:  http://localhost:8080"
echo "   Frontend: http://localhost:5173"
echo ""
echo "Press Ctrl+C to stop all services"

# 捕获退出信号
trap "kill $BACKEND_PID $FRONTEND_PID; exit" INT TERM

# 等待
wait
