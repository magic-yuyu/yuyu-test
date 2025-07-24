@echo off
echo IDaaS 开发工具
echo ================

:menu
echo.
echo 请选择操作:
echo 1. 启动应用
echo 2. 运行测试
echo 3. 构建应用
echo 4. 清理构建文件
echo 5. 检查依赖
echo 6. 退出
echo.

set /p choice=请输入选择 (1-6): 

if "%choice%"=="1" goto start
if "%choice%"=="2" goto test
if "%choice%"=="3" goto build
if "%choice%"=="4" goto clean
if "%choice%"=="5" goto deps
if "%choice%"=="6" goto exit
goto menu

:start
echo 启动应用...
call start-dev.bat
goto menu

:test
echo 运行测试...
go test ./...
goto menu

:build
echo 构建应用...
go build ./cmd/server
echo 构建完成: server.exe
goto menu

:clean
echo 清理构建文件...
del /q server.exe 2>nul
echo 清理完成
goto menu

:deps
echo 检查依赖...
go mod tidy
go mod download
echo 依赖检查完成
goto menu

:exit
echo 再见!
exit /b 0 