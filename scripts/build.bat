@echo off
echo 构建IDaaS应用...

REM 设置环境变量
set GOOS=windows
set GOARCH=amd64
set CGO_ENABLED=0

echo 编译平台: %GOOS/%GOARCH
echo CGO: %CGO_ENABLED%

REM 清理旧的构建文件
if exist server.exe del server.exe

REM 构建应用
echo 开始构建...
go build -ldflags="-s -w" -o server.exe ./cmd/server

if %ERRORLEVEL% EQU 0 (
    echo 构建成功!
    echo 可执行文件: server.exe
    echo 文件大小:
    dir server.exe | find "server.exe"
) else (
    echo 构建失败!
    exit /b 1
)

echo.
echo 运行方式:
echo 1. 直接运行: server.exe
echo 2. 使用脚本: scripts\start-dev.bat
echo 3. 开发工具: scripts\dev-tools.bat 