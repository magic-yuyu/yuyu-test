#!/bin/bash

# 设置颜色
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}========================================"
echo "JWT_SECRET 生成工具"
echo -e "========================================${NC}"
echo

echo -e "${YELLOW}正在生成JWT_SECRET...${NC}"
echo

# 检查OpenSSL是否可用
if ! command -v openssl &> /dev/null; then
    echo -e "${RED}❌ OpenSSL未安装${NC}"
    echo "请安装OpenSSL:"
    echo "Ubuntu/Debian: sudo apt-get install openssl"
    echo "CentOS/RHEL: sudo yum install openssl"
    echo "macOS: brew install openssl"
    exit 1
fi

echo -e "${GREEN}✅ OpenSSL已安装${NC}"
echo

# 生成64字节密钥
echo -e "${YELLOW}📋 生成64字节JWT_SECRET:${NC}"
JWT_SECRET_64=$(openssl rand -base64 64)
echo "$JWT_SECRET_64"
echo

# 生成32字节密钥
echo -e "${YELLOW}📋 生成32字节JWT_SECRET:${NC}"
JWT_SECRET_32=$(openssl rand -base64 32)
echo "$JWT_SECRET_32"
echo

# 生成十六进制密钥
echo -e "${YELLOW}📋 生成十六进制JWT_SECRET:${NC}"
JWT_SECRET_HEX=$(openssl rand -hex 64)
echo "$JWT_SECRET_HEX"
echo

echo -e "${BLUE}========================================"
echo "环境变量配置示例"
echo -e "========================================${NC}"
echo
echo -e "${YELLOW}开发环境 (.env.development):${NC}"
echo "JWT_SECRET=$JWT_SECRET_64"
echo
echo -e "${YELLOW}测试环境 (.env.test):${NC}"
echo "JWT_SECRET=$JWT_SECRET_32"
echo
echo -e "${YELLOW}生产环境 (.env.production):${NC}"
echo "JWT_SECRET=$JWT_SECRET_HEX"
echo

echo -e "${RED}⚠️  安全提醒:${NC}"
echo "- 不要将生产环境的密钥提交到版本控制"
echo "- 定期更换密钥"
echo "- 使用环境变量存储密钥"
echo "- 不同环境使用不同的密钥"
echo

# 保存到文件（可选）
read -p "是否保存密钥到文件？(y/n): " choice
if [[ $choice == "y" || $choice == "Y" ]]; then
    echo "JWT_SECRET_64=$JWT_SECRET_64" > jwt_secrets.txt
    echo "JWT_SECRET_32=$JWT_SECRET_32" >> jwt_secrets.txt
    echo "JWT_SECRET_HEX=$JWT_SECRET_HEX" >> jwt_secrets.txt
    echo -e "${GREEN}✅ 密钥已保存到 jwt_secrets.txt${NC}"
    echo -e "${RED}⚠️  请妥善保管此文件，不要提交到版本控制${NC}"
fi
