@echo off
chcp 65001 >nul
echo 测试环境变量设置...

REM 设置环境变量
set DATABASE_URL=postgresql://postgres:password@localhost:5432/idaas_dev
set JWT_SECRET=your-development-secret-key
set PORT=8080
set GO_ENV=development

echo 环境变量已设置:
echo DATABASE_URL="%DATABASE_URL%"
echo JWT_SECRET="%JWT_SECRET%"
echo PORT="%PORT%"
echo GO_ENV="%GO_ENV%"

echo.
echo 测试Go程序编译...
go build -o test-server.exe cmd/server/main.go

if %ERRORLEVEL% EQU 0 (
    echo 编译成功!
    echo 测试运行...
    test-server.exe
) else (
    echo 编译失败!
)

REM 清理测试文件
if exist test-server.exe del test-server.exe 