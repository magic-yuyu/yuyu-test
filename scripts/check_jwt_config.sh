#!/bin/sh
# 检查JWT相关环境变量配置

warn() { echo "[WARN] $1"; }
fail() { echo "[ERROR] $1"; exit 1; }

ALG=${JWT_ALGORITHM:-HS256}
ENV=${GO_ENV:-development}

if [ "$ALG" = "HS256" ]; then
  [ -z "$JWT_USER_SECRET_KEY" ] && fail "JWT_USER_SECRET_KEY 未设置"
  [ -z "$JWT_SERVICE_SECRET_KEY" ] && fail "JWT_SERVICE_SECRET_KEY 未设置"
  if [ "$ENV" = "production" ]; then
    [ ${#JWT_USER_SECRET_KEY} -lt 32 ] && fail "生产环境 JWT_USER_SECRET_KEY 长度必须>=32"
    [ ${#JWT_SERVICE_SECRET_KEY} -lt 32 ] && fail "生产环境 JWT_SERVICE_SECRET_KEY 长度必须>=32"
  else
    [ ${#JWT_USER_SECRET_KEY} -lt 32 ] && warn "开发环境 JWT_USER_SECRET_KEY 长度<32, 不安全"
    [ ${#JWT_SERVICE_SECRET_KEY} -lt 32 ] && warn "开发环境 JWT_SERVICE_SECRET_KEY 长度<32, 不安全"
  fi
elif [ "$ALG" = "RS256" ] || [ "$ALG" = "ES256" ]; then
  [ -z "$JWT_USER_PRIVATE_KEY" ] && fail "JWT_USER_PRIVATE_KEY 未设置"
  [ -z "$JWT_USER_PUBLIC_KEY" ] && fail "JWT_USER_PUBLIC_KEY 未设置"
  [ -z "$JWT_SERVICE_PRIVATE_KEY" ] && fail "JWT_SERVICE_PRIVATE_KEY 未设置"
  [ -z "$JWT_SERVICE_PUBLIC_KEY" ] && fail "JWT_SERVICE_PUBLIC_KEY 未设置"
else
  fail "不支持的 JWT_ALGORITHM: $ALG"
fi

echo "JWT 配置检测通过 ($ALG)" 