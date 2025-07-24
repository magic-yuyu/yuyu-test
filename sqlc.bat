@echo off
chcp 65001 >nul
echo SQLC 代码生成工具
echo =================

:menu
echo.
echo 请选择操作:
echo 1. 生成代码
echo 2. 验证SQL语法
echo 3. 查看配置
echo 4. 清理生成文件
echo 5. 返回主菜单
echo.

set /p choice=请输入选择 (1-5): 

if "%choice%"=="1" goto generate
if "%choice%"=="2" goto validate
if "%choice%"=="3" goto config
if "%choice%"=="4" goto clean
if "%choice%"=="5" goto exit
echo 无效选择，请重新输入
goto menu

:generate
echo 正在生成Go代码...
sqlc generate
if %ERRORLEVEL% EQU 0 (
    echo [成功] 代码生成完成!
    echo 生成的文件:
    if exist internal\store\database\models.go echo   - internal\store\database\models.go
    if exist internal\store\database\querier.go echo   - internal\store\database\querier.go
    if exist internal\store\database\queries.go echo   - internal\store\database\queries.go
) else (
    echo [错误] 代码生成失败!
    echo 请检查:
    echo   1. 是否已安装sqlc: go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
    echo   2. sqlc.yaml配置是否正确
    echo   3. SQL文件语法是否有误
)
pause
goto menu

:validate
echo 正在验证SQL语法...
sqlc vet
if %ERRORLEVEL% EQU 0 (
    echo [成功] SQL语法验证通过!
) else (
    echo [错误] SQL语法验证失败!
    echo 请检查SQL文件中的语法错误
)
pause
goto menu

:config
echo 当前SQLC配置:
echo ================
if exist sqlc.yaml (
    type sqlc.yaml
) else (
    echo [错误] 未找到sqlc.yaml配置文件
)
echo ================
pause
goto menu

:clean
echo 正在清理生成的文件...
set count=0
if exist internal\store\database\models.go (
    del internal\store\database\models.go
    set /a count+=1
    echo 已删除: internal\store\database\models.go
)
if exist internal\store\database\querier.go (
    del internal\store\database\querier.go
    set /a count+=1
    echo 已删除: internal\store\database\querier.go
)
if exist internal\store\database\queries.go (
    del internal\store\database\queries.go
    set /a count+=1
    echo 已删除: internal\store\database\queries.go
)
if %count% EQU 0 (
    echo 没有找到需要清理的文件
) else (
    echo 清理完成，共删除 %count% 个文件
)
pause
goto menu

:exit
echo 返回主菜单...
exit /b 0