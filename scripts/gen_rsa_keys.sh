#!/bin/sh
# 用于生成RSA密钥对（适用于JWT RS256）
# 用法: ./gen_rsa_keys.sh <private_key_file> <public_key_file>

set -e

PRIV_FILE=${1:-jwt_rsa_private.pem}
PUB_FILE=${2:-jwt_rsa_public.pem}

openssl genpkey -algorithm RSA -out "$PRIV_FILE" -pkeyopt rsa_keygen_bits:2048
openssl rsa -in "$PRIV_FILE" -pubout -out "$PUB_FILE"

echo "私钥: $PRIV_FILE"
echo "公钥: $PUB_FILE" 