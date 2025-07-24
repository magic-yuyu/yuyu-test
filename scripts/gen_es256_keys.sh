#!/bin/sh
# 用于生成ES256(ECDSA)密钥对（适用于JWT ES256）
# 用法: ./gen_es256_keys.sh <private_key_file> <public_key_file>

set -e

PRIV_FILE=${1:-jwt_es256_private.pem}
PUB_FILE=${2:-jwt_es256_public.pem}

openssl ecparam -name prime256v1 -genkey -noout -out "$PRIV_FILE"
openssl ec -in "$PRIV_FILE" -pubout -out "$PUB_FILE"

echo "私钥: $PRIV_FILE"
echo "公钥: $PUB_FILE" 