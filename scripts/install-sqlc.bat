@echo off
echo 正在安装sqlc...

REM 检查是否已安装Go
go version >nul 2>&1
if errorlevel 1 (
    echo 错误: 未安装Go，请先安装Go 1.22或更高版本
    echo 下载地址: https://golang.org/dl/
    pause
    exit /b 1
)

REM 安装sqlc
echo 使用go install安装sqlc...
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

REM 检查安装是否成功
sqlc version >nul 2>&1
if errorlevel 1 (
    echo 错误: sqlc安装失败
    echo 请确保Go的bin目录在PATH环境变量中
    echo 通常路径为: %USERPROFILE%\go\bin
    pause
    exit /b 1
)

echo sqlc安装成功！
echo 版本信息:
sqlc version

echo.
echo 现在可以运行以下命令生成Go代码:
echo sqlc generate

pause 