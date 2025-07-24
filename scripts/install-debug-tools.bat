@echo off
chcp 65001 >nul
echo ========================================
echo IDaaS API调试工具安装脚本
echo ========================================
echo.

echo 正在检查系统环境...
echo.

REM 检查Go是否安装
go version >nul 2>&1
if %errorlevel% neq 0 (
    echo ❌ Go未安装，请先安装Go
    echo 下载地址: https://golang.org/dl/
    pause
    exit /b 1
) else (
    echo ✅ Go已安装
)

REM 检查VS Code是否安装
code --version >nul 2>&1
if %errorlevel% neq 0 (
    echo ❌ VS Code未安装，请先安装VS Code
    echo 下载地址: https://code.visualstudio.com/
    pause
    exit /b 1
) else (
    echo ✅ VS Code已安装
)

echo.
echo ========================================
echo 安装VS Code扩展
echo ========================================
echo.

echo 正在安装REST Client扩展...
code --install-extension humao.rest-client
if %errorlevel% equ 0 (
    echo ✅ REST Client扩展安装成功
) else (
    echo ❌ REST Client扩展安装失败
)

echo 正在安装Go扩展...
code --install-extension golang.go
if %errorlevel% equ 0 (
    echo ✅ Go扩展安装成功
) else (
    echo ❌ Go扩展安装失败
)

echo 正在安装JSON扩展...
code --install-extension ms-vscode.vscode-json
if %errorlevel% equ 0 (
    echo ✅ JSON扩展安装成功
) else (
    echo ❌ JSON扩展安装失败
)

echo 正在安装Docker扩展...
code --install-extension ms-vscode.vscode-docker
if %errorlevel% equ 0 (
    echo ✅ Docker扩展安装成功
) else (
    echo ❌ Docker扩展安装失败
)

echo.
echo ========================================
echo 配置项目环境
echo ========================================
echo.

REM 创建debug目录
if not exist "debug" (
    mkdir debug
    echo ✅ 创建debug目录
)

REM 检查apitest.http文件
if exist "debug\apitest.http" (
    echo ✅ API测试文件已存在
) else (
    echo ❌ API测试文件不存在，请检查项目结构
)

echo.
echo ========================================
echo 安装独立调试工具
echo ========================================
echo.

echo 推荐安装以下独立工具：
echo.
echo 1. Postman - 功能强大的API测试工具
echo    下载地址: https://www.postman.com/downloads/
echo.
echo 2. Insomnia - 轻量级API客户端
echo    下载地址: https://insomnia.rest/download
echo.
echo 3. curl - 命令行HTTP客户端（通常已预装）
echo.

set /p choice="是否要打开下载页面？(y/n): "
if /i "%choice%"=="y" (
    echo 正在打开下载页面...
    start https://www.postman.com/downloads/
    start https://insomnia.rest/download
)

echo.
echo ========================================
echo 配置说明
echo ========================================
echo.

echo 安装完成后，请按以下步骤使用：
echo.
echo 1. 打开VS Code
echo 2. 打开项目文件夹
echo 3. 打开 debug\apitest.http 文件
echo 4. 安装REST Client扩展（如果未自动安装）
echo 5. 点击请求上方的"Send Request"按钮
echo.

echo 环境变量配置：
echo - baseUrl: http://localhost:8080
echo - apiKey: 从创建租户的响应中获取
echo.

echo 更多配置信息请查看：
echo - docs\API_DEBUG_TOOLS.md
echo - debug\apitest.http
echo.

echo ========================================
echo 安装完成！
echo ========================================
echo.

pause 