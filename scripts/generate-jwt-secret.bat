@echo off
chcp 65001 >nul
echo ========================================
echo JWT_SECRET 生成工具
echo ========================================
echo.

echo 正在生成JWT_SECRET...
echo.

REM 检查OpenSSL是否可用
openssl version >nul 2>&1
if %errorlevel% neq 0 (
    echo ❌ OpenSSL未安装，请先安装OpenSSL
    echo 下载地址: https://slproweb.com/products/Win32OpenSSL.html
    pause
    exit /b 1
)

echo ✅ OpenSSL已安装
echo.

REM 生成64字节密钥
echo �� 生成64字节JWT_SECRET:
for /f "delims=" %%i in ('openssl rand -base64 64') do set JWT_SECRET=%%i
echo %JWT_SECRET%
echo.

REM 生成32字节密钥
echo �� 生成32字节JWT_SECRET:
for /f "delims=" %%i in ('openssl rand -base64 32') do set JWT_SECRET_32=%%i
echo %JWT_SECRET_32%
echo.

REM 生成十六进制密钥
echo �� 生成十六进制JWT_SECRET:
for /f "delims=" %%i in ('openssl rand -hex 64') do set JWT_SECRET_HEX=%%i
echo %JWT_SECRET_HEX%
echo.

echo ========================================
echo 环境变量配置示例
echo ========================================
echo.
echo 开发环境 (.env.development):
echo JWT_SECRET=%JWT_SECRET%
echo.
echo 测试环境 (.env.test):
echo JWT_SECRET=%JWT_SECRET_32%
echo.
echo 生产环境 (.env.production):
echo JWT_SECRET=%JWT_SECRET_HEX%
echo.

echo ⚠️  安全提醒:
echo - 不要将生产环境的密钥提交到版本控制
echo - 定期更换密钥
echo - 使用环境变量存储密钥
echo.

pause