# IDaaS API 使用指南

## 概述

IDaaS平台提供完整的身份认证即服务功能，支持多租户架构，每个租户都有独立的用户管理。

## 快速开始

### 1. 创建租户

首先需要创建一个租户来获得API密钥：

```bash
curl -X POST http://localhost:8080/v1/tenants \
  -H "Content-Type: application/json" \
  -d '{
    "name": "我的应用"
  }'
```

响应示例：
```json
{
  "id": "tnt_abc123...",
  "name": "我的应用",
  "public_key": "pub_xyz789...",
  "secret_key": "sec_def456...",
  "created_at": "2024-01-01T00:00:00Z"
}
```

### 2. 用户注册

使用租户的公钥或密钥进行用户注册：

```bash
curl -X POST http://localhost:8080/v1/auth/register \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer pub_xyz789..." \
  -d '{
    "email": "user@example.com",
    "password": "password123",
    "profile": {
      "name": "张三",
      "role": "user"
    }
  }'
```

### 3. 用户登录

用户登录获取JWT令牌：

```bash
curl -X POST http://localhost:8080/v1/auth/login \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer pub_xyz789..." \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

响应示例：
```json
{
  "user": {
    "id": "usr_abc123...",
    "email": "user@example.com",
    "profile": {},
    "created_at": "2024-01-01T00:00:00Z"
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### 4. 获取用户信息

使用JWT令牌获取当前用户信息：

```bash
curl -X GET http://localhost:8080/v1/users/me \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### 5. 管理用户

使用租户密钥管理用户：

```bash
# 获取所有用户
curl -X GET http://localhost:8080/v1/users \
  -H "Authorization: Bearer sec_def456..."

# 获取指定用户
curl -X GET http://localhost:8080/v1/users/usr_abc123... \
  -H "Authorization: Bearer sec_def456..."
```

## 认证方式

### API密钥认证
- 用于服务器间调用
- 使用租户的公钥或密钥
- 格式：`Authorization: Bearer <api_key>`

### JWT认证
- 用于用户会话
- 登录后获得JWT令牌
- 格式：`Authorization: Bearer <jwt_token>`

## 错误处理

所有API都返回标准的HTTP状态码和JSON错误信息：

```json
{
  "error": "错误描述"
}
```

常见状态码：
- `200` - 成功
- `201` - 创建成功
- `400` - 请求参数错误
- `401` - 认证失败
- `404` - 资源不存在
- `500` - 服务器内部错误

## 环境变量

部署时需要设置以下环境变量：

```bash
DATABASE_URL=postgresql://user:password@host:port/database
JWT_SECRET=your_jwt_secret_key
PORT=8080
``` 