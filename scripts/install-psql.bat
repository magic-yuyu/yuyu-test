@echo off
echo 正在安装PostgreSQL客户端工具...

REM 检查是否已安装psql
psql --version >nul 2>&1
if not errorlevel 1 (
    echo PostgreSQL客户端已安装
    psql --version
    pause
    exit /b 0
)

echo 未检测到PostgreSQL客户端工具
echo.
echo 请选择安装方式:
echo 1. 下载完整PostgreSQL安装包（推荐）
echo 2. 下载仅客户端工具
echo 3. 使用Chocolatey安装（如果已安装Chocolatey）
echo 4. 手动安装
echo.

set /p choice="请选择 (1-4): "

if "%choice%"=="1" (
    echo 正在打开PostgreSQL下载页面...
    start https://www.postgresql.org/download/windows/
    echo 请下载并安装PostgreSQL，安装时选择"Command Line Tools"
) else if "%choice%"=="2" (
    echo 正在打开EnterpriseDB下载页面...
    start https://www.enterprisedb.com/download-postgresql-binaries
    echo 请下载对应版本的"Command Line Tools"
) else if "%choice%"=="3" (
    echo 使用Chocolatey安装...
    choco install postgresql --params '/Password:password' --yes
    if errorlevel 1 (
        echo Chocolatey安装失败，请手动安装
    ) else (
        echo PostgreSQL安装成功
    )
) else if "%choice%"=="4" (
    echo 手动安装说明:
    echo 1. 访问 https://www.postgresql.org/download/windows/
    echo 2. 下载PostgreSQL安装包
    echo 3. 运行安装程序，选择"Command Line Tools"
    echo 4. 将安装目录的bin文件夹添加到PATH环境变量
    echo 5. 重新打开命令行窗口
) else (
    echo 无效选择
    pause
    exit /b 1
)

echo.
echo 安装完成后，请重新运行此脚本验证安装
pause 