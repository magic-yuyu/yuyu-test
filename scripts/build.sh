#!/bin/bash

echo "构建IDaaS应用..."

# 设置环境变量
export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=0

echo "编译平台: $GOOS/$GOARCH"
echo "CGO: $CGO_ENABLED"

# 清理旧的构建文件
if [ -f "server" ]; then
    rm server
fi

# 构建应用
echo "开始构建..."
go build -ldflags="-s -w" -o server ./cmd/server

if [ $? -eq 0 ]; then
    echo "构建成功!"
    echo "可执行文件: server"
    echo "文件大小:"
    ls -lh server
else
    echo "构建失败!"
    exit 1
fi

echo ""
echo "运行方式:"
echo "1. 直接运行: ./server"
echo "2. 使用脚本: ./scripts/start-dev.sh" 