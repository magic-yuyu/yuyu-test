@echo off
echo 执行数据库迁移...

REM 设置环境变量
set DATABASE_URL=postgresql://postgres:password@localhost:5432/idaas_dev

echo 连接到数据库: %DATABASE_URL%

REM 这里需要安装golang-migrate工具
REM 或者手动执行SQL文件
echo 请手动执行以下SQL文件:
echo migrations/0001_initial_schema.up.sql

echo.
echo 或者安装golang-migrate后运行:
echo migrate -database "%DATABASE_URL%" -path migrations up 