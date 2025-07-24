@echo off
chcp 65001 >nul
echo ========================================
echo IDaaS 开发环境启动
echo ========================================

echo.
echo 注意: 请确保PostgreSQL数据库已启动
echo 如果没有数据库，请运行: scripts\start-complete.bat

echo.
echo 设置环境变量...
set DATABASE_URL=postgresql://postgres:password@47.95.195.62:5432/idaas_dev?sslmode=disable
set JWT_SECRET=your-development-secret-key
set JWT_USER_SECRET_KEY=your-user-secret-key
set JWT_SERVICE_SECRET_KEY=your-service-secret-key
set USER_TOKEN_EXPIRATION=3600
set SERVICE_TOKEN_EXPIRATION=300
set PORT=8080
set GO_ENV=development

echo 环境变量已设置:
echo DATABASE_URL=%DATABASE_URL%
echo JWT_SECRET=%JWT_SECRET%
echo PORT=%PORT%
echo GO_ENV=%GO_ENV%
echo JWT_USER_SECRET_KEY=%JWT_USER_SECRET_KEY%
echo JWT_SERVICE_SECRET_KEY=%JWT_SERVICE_SECRET_KEY%
echo USER_TOKEN_EXPIRATION=%USER_TOKEN_EXPIRATION%
echo SERVICE_TOKEN_EXPIRATION=%SERVICE_TOKEN_EXPIRATION%

echo.
echo 启动应用...
go run cmd/server/main.go 