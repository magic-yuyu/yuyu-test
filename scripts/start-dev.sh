#!/bin/bash

echo "启动IDaaS开发环境..."

# 设置环境变量
export DATABASE_URL="postgresql://postgres:password@localhost:5432/idaas_dev"
export JWT_SECRET="your-development-secret-key"
export PORT=8080
export GO_ENV=development
export JWT_USER_SECRET_KEY="your-user-secret-key"
export JWT_SERVICE_SECRET_KEY="your-service-secret-key"
export USER_TOKEN_EXPIRATION=3600
export SERVICE_TOKEN_EXPIRATION=300

echo "环境变量已设置:"
echo "DATABASE_URL=$DATABASE_URL"
echo "JWT_SECRET=$JWT_SECRET"
echo "PORT=$PORT"
echo "GO_ENV=$GO_ENV"
echo "JWT_USER_SECRET_KEY=$JWT_USER_SECRET_KEY"
echo "JWT_SERVICE_SECRET_KEY=$JWT_SERVICE_SECRET_KEY"
echo "USER_TOKEN_EXPIRATION=$USER_TOKEN_EXPIRATION"
echo "SERVICE_TOKEN_EXPIRATION=$SERVICE_TOKEN_EXPIRATION"

echo ""
echo "启动应用..."
go run cmd/server/main.go 