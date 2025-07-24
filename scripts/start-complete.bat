@echo off
chcp 65001 >nul
echo ========================================
echo IDaaS 完整开发环境启动脚本
echo ========================================

echo.
echo 1. 检查Docker是否运行...
docker version >nul 2>&1
if %ERRORLEVEL% NEQ 0 (
    echo 错误: Docker未运行或未安装
    echo 请先启动Docker Desktop
    pause
    exit /b 1
)
echo Docker运行正常

echo.
echo 2. 启动PostgreSQL数据库...
docker run -d --name postgres-idaas ^
  -e POSTGRES_DB=idaas_dev ^
  -e POSTGRES_USER=postgres ^
  -e POSTGRES_PASSWORD=password ^
  -p 5432:5432 ^
  postgres:15-alpine

if %ERRORLEVEL% NEQ 0 (
    echo 数据库启动失败，可能已存在同名容器
    echo 尝试删除旧容器...
    docker rm -f postgres-idaas >nul 2>&1
    echo 重新启动数据库...
    docker run -d --name postgres-idaas ^
      -e POSTGRES_DB=idaas_dev ^
      -e POSTGRES_USER=postgres ^
      -e POSTGRES_PASSWORD=password ^
      -p 5432:5432 ^
      postgres:15-alpine
)

echo 等待数据库启动...
timeout /t 5 /nobreak >nul

echo.
echo 3. 执行数据库迁移...
echo 请手动执行以下SQL到数据库:
echo migrations/0001_initial_schema.up.sql
echo.
echo 或者使用数据库客户端连接到:
echo Host: localhost
echo Port: 5432
echo Database: idaas_dev
echo Username: postgres
echo Password: password

echo.
echo 4. 启动IDaaS应用...
echo 环境变量已设置:
set DATABASE_URL=postgresql://postgres:password@localhost:5432/idaas_dev
set JWT_SECRET=your-development-secret-key
set PORT=8080
set GO_ENV=development

echo DATABASE_URL=%DATABASE_URL%
echo JWT_SECRET=%JWT_SECRET%
echo PORT=%PORT%
echo GO_ENV=%GO_ENV%

echo.
echo 启动应用...
go run cmd/server/main.go

echo.
echo 应用已停止
echo 清理数据库容器...
docker rm -f postgres-idaas >nul 2>&1
echo 清理完成 