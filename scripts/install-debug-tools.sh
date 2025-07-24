#!/bin/bash

# 设置颜色
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}========================================"
echo "IDaaS API调试工具安装脚本"
echo -e "========================================${NC}"
echo

echo -e "${YELLOW}正在检查系统环境...${NC}"
echo

# 检查Go是否安装
if command -v go &> /dev/null; then
    echo -e "${GREEN}✅ Go已安装${NC}"
    go version
else
    echo -e "${RED}❌ Go未安装，请先安装Go${NC}"
    echo "下载地址: https://golang.org/dl/"
    exit 1
fi

# 检查VS Code是否安装
if command -v code &> /dev/null; then
    echo -e "${GREEN}✅ VS Code已安装${NC}"
    code --version
else
    echo -e "${RED}❌ VS Code未安装，请先安装VS Code${NC}"
    echo "下载地址: https://code.visualstudio.com/"
    exit 1
fi

echo
echo -e "${BLUE}========================================"
echo "安装VS Code扩展"
echo -e "========================================${NC}"
echo

# 安装REST Client扩展
echo -e "${YELLOW}正在安装REST Client扩展...${NC}"
if code --install-extension humao.rest-client; then
    echo -e "${GREEN}✅ REST Client扩展安装成功${NC}"
else
    echo -e "${RED}❌ REST Client扩展安装失败${NC}"
fi

# 安装Go扩展
echo -e "${YELLOW}正在安装Go扩展...${NC}"
if code --install-extension golang.go; then
    echo -e "${GREEN}✅ Go扩展安装成功${NC}"
else
    echo -e "${RED}❌ Go扩展安装失败${NC}"
fi

# 安装JSON扩展
echo -e "${YELLOW}正在安装JSON扩展...${NC}"
if code --install-extension ms-vscode.vscode-json; then
    echo -e "${GREEN}✅ JSON扩展安装成功${NC}"
else
    echo -e "${RED}❌ JSON扩展安装失败${NC}"
fi

# 安装Docker扩展
echo -e "${YELLOW}正在安装Docker扩展...${NC}"
if code --install-extension ms-vscode.vscode-docker; then
    echo -e "${GREEN}✅ Docker扩展安装成功${NC}"
else
    echo -e "${RED}❌ Docker扩展安装失败${NC}"
fi

echo
echo -e "${BLUE}========================================"
echo "配置项目环境"
echo -e "========================================${NC}"
echo

# 创建debug目录
if [ ! -d "debug" ]; then
    mkdir debug
    echo -e "${GREEN}✅ 创建debug目录${NC}"
fi

# 检查apitest.http文件
if [ -f "debug/apitest.http" ]; then
    echo -e "${GREEN}✅ API测试文件已存在${NC}"
else
    echo -e "${RED}❌ API测试文件不存在，请检查项目结构${NC}"
fi

echo
echo -e "${BLUE}========================================"
echo "安装独立调试工具"
echo -e "========================================${NC}"
echo

echo -e "${YELLOW}推荐安装以下独立工具：${NC}"
echo
echo "1. Postman - 功能强大的API测试工具"
echo "   下载地址: https://www.postman.com/downloads/"
echo
echo "2. Insomnia - 轻量级API客户端"
echo "   下载地址: https://insomnia.rest/download"
echo
echo "3. curl - 命令行HTTP客户端（通常已预装）"

# 检查curl
if command -v curl &> /dev/null; then
    echo -e "${GREEN}✅ curl已安装${NC}"
else
    echo -e "${RED}❌ curl未安装${NC}"
fi

echo
read -p "是否要打开下载页面？(y/n): " choice
if [[ $choice == "y" || $choice == "Y" ]]; then
    echo -e "${YELLOW}正在打开下载页面...${NC}"
    
    # 检测操作系统并打开浏览器
    if [[ "$OSTYPE" == "darwin"* ]]; then
        # macOS
        open https://www.postman.com/downloads/
        open https://insomnia.rest/download
    elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
        # Linux
        if command -v xdg-open &> /dev/null; then
            xdg-open https://www.postman.com/downloads/
            xdg-open https://insomnia.rest/download
        elif command -v gnome-open &> /dev/null; then
            gnome-open https://www.postman.com/downloads/
            gnome-open https://insomnia.rest/download
        else
            echo "无法自动打开浏览器，请手动访问："
            echo "https://www.postman.com/downloads/"
            echo "https://insomnia.rest/download"
        fi
    fi
fi

echo
echo -e "${BLUE}========================================"
echo "配置说明"
echo -e "========================================${NC}"
echo

echo -e "${YELLOW}安装完成后，请按以下步骤使用：${NC}"
echo
echo "1. 打开VS Code"
echo "2. 打开项目文件夹"
echo "3. 打开 debug/apitest.http 文件"
echo "4. 安装REST Client扩展（如果未自动安装）"
echo "5. 点击请求上方的'Send Request'按钮"
echo

echo -e "${YELLOW}环境变量配置：${NC}"
echo "- baseUrl: http://localhost:8080"
echo "- apiKey: 从创建租户的响应中获取"
echo

echo -e "${YELLOW}更多配置信息请查看：${NC}"
echo "- docs/API_DEBUG_TOOLS.md"
echo "- debug/apitest.http"
echo

echo -e "${BLUE}========================================"
echo "安装完成！"
echo -e "========================================${NC}"
echo

# 设置脚本可执行权限
chmod +x scripts/install-debug-tools.sh

echo -e "${GREEN}脚本已设置为可执行权限${NC}"
echo -e "${YELLOW}下次可以直接运行: ./scripts/install-debug-tools.sh${NC}" 