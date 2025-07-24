# IDaaS API 接口文档

## 📋 目录
- [概述](#概述)
- [认证方式](#认证方式)
- [基础信息](#基础信息)
- [API端点](#api端点)
- [错误处理](#错误处理)
- [示例代码](#示例代码)
- [最佳实践](#最佳实践)

## 概述

IDaaS（身份认证即服务）平台提供完整的用户身份管理功能，支持多租户架构。每个租户都有独立的用户管理空间，确保数据隔离和安全。

### 主要特性
- 🔐 **多租户支持** - 每个客户应用独立管理
- 🛡️ **双重认证** - API密钥 + JWT令牌
- 👥 **用户管理** - 完整的用户生命周期管理
- 🔑 **安全存储** - 密码加密存储
- 📊 **灵活配置** - 支持自定义用户属性

## 认证方式

### 1. API密钥认证
用于服务器间调用，需要在请求头中提供：
```
Authorization: Bearer {api_key}
```

**支持的密钥类型：**
- **Public Key**: 用于客户端调用（用户注册、登录）
- **Secret Key**: 用于管理调用（用户管理、租户管理）

### 2. JWT令牌认证
用于用户会话认证，需要在请求头中提供：
```
Authorization: Bearer {jwt_token}
```

## 基础信息

### 基础URL
```
开发环境: http://localhost:8080
生产环境: https://your-domain.com
```

### 请求格式
- **Content-Type**: `application/json`
- **字符编码**: `UTF-8`

### 响应格式
所有响应均为JSON格式，包含以下字段：
- `data`: 响应数据（成功时）
- `error`: 错误信息（失败时）
- `message`: 状态消息

## API端点

### 健康检查

#### GET /health
检查服务健康状态

**请求参数**: 无

**响应示例**:
```json
{
  "status": "ok"
}
```

---

### 租户管理

#### POST /v1/tenants
创建新租户

**认证**: 无需认证

**请求参数**:
```json
{
  "name": "string"  // 租户名称，必填
}
```

**响应示例**:
```json
{
  "id": "tnt_abc123def456",
  "name": "我的应用",
  "public_key": "pub_xyz789abc123",
  "secret_key": "sec_def456ghi789",
  "created_at": "2024-01-01T00:00:00Z"
}
```

**字段说明**:
- `id`: 租户唯一标识符
- `name`: 租户名称
- `public_key`: 客户端API密钥
- `secret_key`: 管理API密钥（请妥善保管）
- `created_at`: 创建时间

#### GET /v1/tenants/:id
获取租户信息

**认证**: 无需认证

**路径参数**:
- `id`: 租户ID

**响应示例**:
```json
{
  "id": "tnt_abc123def456",
  "name": "我的应用",
  "api_secret_key_hash": "hashed_secret_key",
  "api_public_key": "pub_xyz789abc123",
  "created_at": "2024-01-01T00:00:00Z"
}
```

---

### 用户认证

#### POST /v1/auth/register
用户注册

**认证**: 需要API密钥（Public Key或Secret Key）

**请求头**:
```
Authorization: Bearer {api_key}
```

**请求参数**:
```json
{
  "email": "user@example.com",     // 邮箱地址，必填，格式验证
  "password": "password123",       // 密码，必填，最少6位
  "profile": {                     // 用户属性，可选
    "name": "张三",
    "role": "user",
    "department": "技术部"
  }
}
```

**响应示例**:
```json
{
  "id": "usr_def456ghi789",
  "email": "user@example.com",
  "profile": {
    "name": "张三",
    "role": "user",
    "department": "技术部"
  },
  "created_at": "2024-01-01T00:00:00Z"
}
```

#### POST /v1/auth/login
用户登录

**认证**: 需要API密钥（Public Key或Secret Key）

**请求头**:
```
Authorization: Bearer {api_key}
```

**请求参数**:
```json
{
  "email": "user@example.com",     // 邮箱地址，必填
  "password": "password123"        // 密码，必填
}
```

**响应示例**:
```json
{
  "user": {
    "id": "usr_def456ghi789",
    "email": "user@example.com",
    "profile": {
      "name": "张三",
      "role": "user",
      "department": "技术部"
    },
    "created_at": "2024-01-01T00:00:00Z"
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**字段说明**:
- `user`: 用户信息
- `token`: JWT访问令牌（有效期24小时）

---

### 用户管理

#### GET /v1/users/me
获取当前用户信息

**认证**: 需要JWT令牌

**请求头**:
```
Authorization: Bearer {jwt_token}
```

**响应示例**:
```json
{
  "id": "usr_def456ghi789",
  "email": "user@example.com",
  "profile": {
    "name": "张三",
    "role": "user",
    "department": "技术部"
  },
  "created_at": "2024-01-01T00:00:00Z"
}
```

#### GET /v1/users
获取租户下的所有用户

**认证**: 需要API密钥（Secret Key）

**请求头**:
```
Authorization: Bearer {secret_key}
```

**响应示例**:
```json
{
  "users": [
    {
      "id": "usr_def456ghi789",
      "email": "user1@example.com",
      "profile": {
        "name": "张三",
        "role": "user"
      },
      "created_at": "2024-01-01T00:00:00Z"
    },
    {
      "id": "usr_abc123def456",
      "email": "user2@example.com",
      "profile": {
        "name": "李四",
        "role": "admin"
      },
      "created_at": "2024-01-02T00:00:00Z"
    }
  ]
}
```

#### GET /v1/users/:id
获取指定用户信息

**认证**: 需要API密钥（Secret Key）

**请求头**:
```
Authorization: Bearer {secret_key}
```

**路径参数**:
- `id`: 用户ID

**响应示例**:
```json
{
  "id": "usr_def456ghi789",
  "email": "user@example.com",
  "profile": {
    "name": "张三",
    "role": "user",
    "department": "技术部"
  },
  "created_at": "2024-01-01T00:00:00Z"
}
```

## 错误处理

### HTTP状态码
- `200` - 请求成功
- `201` - 创建成功
- `400` - 请求参数错误
- `401` - 认证失败
- `404` - 资源不存在
- `500` - 服务器内部错误

### 错误响应格式
```json
{
  "error": "错误描述信息"
}
```

### 常见错误
| 错误码 | 错误信息                                   | 说明           |
| ------ | ------------------------------------------ | -------------- |
| 400    | `"email" is required`                      | 邮箱字段必填   |
| 400    | `"password" is required`                   | 密码字段必填   |
| 400    | `"password" must be at least 6 characters` | 密码长度不足   |
| 400    | `"email" is not a valid email`             | 邮箱格式错误   |
| 401    | `"tenant not found"`                       | 租户不存在     |
| 401    | `"invalid email or password"`              | 邮箱或密码错误 |
| 401    | `"user not authenticated"`                 | 用户未认证     |
| 404    | `"tenant not found"`                       | 租户不存在     |
| 404    | `"user not found"`                         | 用户不存在     |

## 示例代码

### JavaScript/Node.js

#### 创建租户
```javascript
const axios = require('axios');

async function createTenant() {
  try {
    const response = await axios.post('http://localhost:8080/v1/tenants', {
      name: '我的应用'
    });
    
    console.log('租户创建成功:', response.data);
    return response.data;
  } catch (error) {
    console.error('创建租户失败:', error.response.data);
  }
}
```

#### 用户注册
```javascript
async function registerUser(apiKey, userData) {
  try {
    const response = await axios.post('http://localhost:8080/v1/auth/register', {
      email: userData.email,
      password: userData.password,
      profile: {
        name: userData.name,
        role: 'user'
      }
    }, {
      headers: {
        'Authorization': `Bearer ${apiKey}`,
        'Content-Type': 'application/json'
      }
    });
    
    console.log('用户注册成功:', response.data);
    return response.data;
  } catch (error) {
    console.error('用户注册失败:', error.response.data);
  }
}
```

#### 用户登录
```javascript
async function loginUser(apiKey, credentials) {
  try {
    const response = await axios.post('http://localhost:8080/v1/auth/login', {
      email: credentials.email,
      password: credentials.password
    }, {
      headers: {
        'Authorization': `Bearer ${apiKey}`,
        'Content-Type': 'application/json'
      }
    });
    
    console.log('登录成功:', response.data);
    return response.data;
  } catch (error) {
    console.error('登录失败:', error.response.data);
  }
}
```

#### 获取用户信息
```javascript
async function getUserInfo(jwtToken) {
  try {
    const response = await axios.get('http://localhost:8080/v1/users/me', {
      headers: {
        'Authorization': `Bearer ${jwtToken}`
      }
    });
    
    console.log('用户信息:', response.data);
    return response.data;
  } catch (error) {
    console.error('获取用户信息失败:', error.response.data);
  }
}
```

### Python

#### 创建租户
```python
import requests

def create_tenant():
    try:
        response = requests.post('http://localhost:8080/v1/tenants', json={
            'name': '我的应用'
        })
        response.raise_for_status()
        
        tenant_data = response.json()
        print('租户创建成功:', tenant_data)
        return tenant_data
    except requests.exceptions.RequestException as e:
        print('创建租户失败:', e)
```

#### 用户注册
```python
def register_user(api_key, user_data):
    try:
        headers = {
            'Authorization': f'Bearer {api_key}',
            'Content-Type': 'application/json'
        }
        
        response = requests.post('http://localhost:8080/v1/auth/register', 
                               json={
                                   'email': user_data['email'],
                                   'password': user_data['password'],
                                   'profile': {
                                       'name': user_data['name'],
                                       'role': 'user'
                                   }
                               }, headers=headers)
        response.raise_for_status()
        
        user_data = response.json()
        print('用户注册成功:', user_data)
        return user_data
    except requests.exceptions.RequestException as e:
        print('用户注册失败:', e)
```

### cURL

#### 创建租户
```bash
curl -X POST http://localhost:8080/v1/tenants \
  -H "Content-Type: application/json" \
  -d '{
    "name": "我的应用"
  }'
```

#### 用户注册
```bash
curl -X POST http://localhost:8080/v1/auth/register \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer pub_xyz789abc123" \
  -d '{
    "email": "user@example.com",
    "password": "password123",
    "profile": {
      "name": "张三",
      "role": "user"
    }
  }'
```

#### 用户登录
```bash
curl -X POST http://localhost:8080/v1/auth/login \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer pub_xyz789abc123" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

#### 获取用户信息
```bash
curl -X GET http://localhost:8080/v1/users/me \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

## 最佳实践

### 1. 安全性
- 🔐 **妥善保管Secret Key** - 不要在前端代码中使用
- 🔄 **定期轮换密钥** - 建议定期更新API密钥
- 🛡️ **HTTPS传输** - 生产环境必须使用HTTPS
- ⏰ **令牌过期** - JWT令牌有效期为24小时

### 2. 错误处理
- ✅ **检查响应状态码** - 始终验证HTTP状态码
- 📝 **记录错误日志** - 记录详细的错误信息
- 🔄 **重试机制** - 对临时错误实现重试逻辑
- 🚫 **用户友好提示** - 向用户显示友好的错误信息

### 3. 性能优化
- 📦 **批量操作** - 避免频繁的单次请求
- 🗄️ **缓存策略** - 缓存用户信息和租户信息
- ⏱️ **超时设置** - 设置合理的请求超时时间
- 📊 **监控指标** - 监控API调用频率和响应时间

### 4. 开发建议
- 🧪 **测试环境** - 使用独立的测试环境
- 📚 **文档同步** - 保持代码和文档的一致性
- 🔄 **版本控制** - 使用API版本控制
- 📝 **日志记录** - 记录重要的操作日志

### 5. 数据管理
- 🗂️ **数据备份** - 定期备份重要数据
- 🔍 **数据验证** - 验证所有输入数据
- 📊 **数据统计** - 监控用户增长和活跃度
- 🗑️ **数据清理** - 定期清理过期数据

---

## 对内服务认证（M2M）

### POST /oauth/token

**用途**：内部微服务获取服务间通信的访问令牌（服务JWT）

**认证方式**：Basic Auth（client_id/client_secret）

**请求头**：
```
Authorization: Basic base64(client_id:client_secret)
Content-Type: application/x-www-form-urlencoded
```

**请求体**：
```
grant_type=client_credentials
```

**响应示例**：
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "token_type": "Bearer",
  "expires_in": 300
}
```

**错误响应**：
```json
{
  "error": "Client not found"
}
```

**说明**：
- 仅支持grant_type=client_credentials
- access_token为服务JWT，包含sub（client_id）、scope、exp、iss等字段
- 令牌有效期5分钟
- 需先在数据库注册internal_client并分配scope

## 服务Token获取接口说明

### 1. /oauth/token
- **用途**：标准OAuth2 Client Credentials授权，适合对外API、三方平台、API Gateway等标准OAuth2场景。
- **认证方式**：HTTP Basic Auth（Authorization: Basic base64(client_id:client_secret)）
- **请求体**：
  ```
  grant_type=client_credentials
  ```
- **响应**：
  ```json
  {
    "access_token": "...",
    "token_type": "Bearer",
    "expires_in": 300
  }
  ```
- **调用示例（curl）**：
  ```bash
  curl -X POST http://localhost:8080/oauth/token \
    -H "Authorization: Basic $(echo -n 'client_id:client_secret' | base64)" \
    -H "Content-Type: application/x-www-form-urlencoded" \
    -d 'grant_type=client_credentials'
  ```
- **适用场景**：
  - 标准OAuth2对接
  - 云平台API Gateway
  - 需要标准协议的三方集成

### 2. /v1/internal/services/authenticate
- **用途**：平台自定义服务认证，适合内部微服务、自动化脚本、非标准OAuth2客户端。
- **认证方式**：JSON体传递 client_id 和 client_secret
- **请求体**：
  ```json
  {
    "client_id": "your-client-id",
    "client_secret": "your-client-secret"
  }
  ```
- **响应**：
  ```json
  {
    "access_token": "...",
    "token_type": "Bearer",
    "expires_in": 300,
    "scopes": ["user:read", "user:write"]
  }
  ```
- **调用示例（curl）**：
  ```bash
  curl -X POST http://localhost:8080/v1/internal/services/authenticate \
    -H "Content-Type: application/json" \
    -d '{"client_id": "your-client-id", "client_secret": "your-client-secret"}'
  ```
- **适用场景**：
  - 平台内部微服务间认证
  - 自动化脚本、CI/CD工具
  - 需要灵活扩展的自定义集成

### 3. 选择建议
- **对外/三方/标准OAuth2场景**：优先使用 `/oauth/token`
- **平台内部/自用/脚本**：优先使用 `/v1/internal/services/authenticate`

## 支持与反馈

如果您在使用过程中遇到问题或有改进建议，请通过以下方式联系我们：

- 📧 **邮箱**: support@idaas.com
- 📖 **文档**: https://docs.idaas.com
- 🐛 **问题反馈**: https://github.com/idaas/issues

---

*最后更新时间: 2024年1月* 

---

## 内部服务管理API

### 1. 服务注册
#### POST /v1/internal/services/register
- **认证**：无需认证
- **请求参数**：
```json
{
  "service_name": "string", // 服务名称，必填
  "description": "string"    // 服务描述，可选
}
```
- **响应示例**：
```json
{
  "client_id": "svc_abc123def456",
  "client_secret": "secret_abcdefg123456", // 仅返回一次
  "service_name": "服务A",
  "description": "内部服务A",
  "created_at": "2024-01-01T00:00:00Z",
  "message": "Service registered successfully",
  "warning": "Please save the client_secret securely. It will not be shown again."
}
```

### 2. 服务认证（获取JWT）
#### POST /v1/internal/services/authenticate
- **认证**：无需认证
- **请求参数**：
```json
{
  "client_id": "svc_abc123def456",
  "client_secret": "secret_abcdefg123456"
}
```
- **响应示例**：
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "token_type": "Bearer",
  "expires_in": 300
}
```

### 3. 验证服务JWT
#### POST /v1/internal/services/validate-token
- **认证**：无需认证
- **请求参数**：
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```
- **响应示例**：
```json
{
  "valid": true,
  "client_id": "svc_abc123def456",
  "scopes": ["user:read", "user:write"],
  "exp": 1735680000
}
```

### 4. 授权Scope权限
#### POST /v1/internal/services/grant-scope
- **认证**：Bearer 服务JWT（需 internal:admin 权限）
- **请求参数**：
```json
{
  "client_id": "svc_abc123def456",
  "scope_name": "user:read"
}
```
- **响应示例**：
```json
{
  "message": "Scope granted successfully"
}
```

### 5. 撤销Scope权限
#### POST /v1/internal/services/revoke-scope
- **认证**：Bearer 服务JWT（需 internal:admin 权限）
- **请求参数**：
```json
{
  "client_id": "svc_abc123def456",
  "scope_name": "user:read"
}
```
- **响应示例**：
```json
{
  "message": "Scope revoked successfully"
}
```

### 6. 检查权限
#### POST /v1/internal/services/check-permission
- **认证**：Bearer 服务JWT
- **请求参数**：
```json
{
  "client_id": "svc_abc123def456",
  "scope_name": "user:read"
}
```
- **响应示例**：
```json
{
  "has_permission": true
}
```

### 7. 服务列表
#### GET /v1/internal/services
- **认证**：Bearer 服务JWT（需 internal:admin 权限）
- **响应示例**：
```json
{
  "services": [
    {
      "client_id": "svc_abc123def456",
      "service_name": "服务A",
      "description": "内部服务A",
      "created_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

### 8. 获取服务访问日志
#### GET /v1/internal/services/{client_id}/logs?limit=50&offset=0
- **认证**：Bearer 服务JWT（需 internal:admin 权限）
- **响应示例**：
```json
[
  {
    "id": 1,
    "client_id": "svc_abc123def456",
    "endpoint": "/api/internal/users",
    "method": "GET",
    "status_code": 200,
    "response_time_ms": 35,
    "ip_address": "127.0.0.1",
    "user_agent": "curl/7.68.0",
    "created_at": "2024-01-01T00:00:00Z"
  }
]
```

### 9. 获取服务统计信息
#### GET /v1/internal/services/{client_id}/statistics?since=24h
- **认证**：Bearer 服务JWT（需 internal:admin 权限）
- **响应示例**：
```json
{
  "client_id": "svc_abc123def456",
  "since": "2024-01-01T00:00:00Z",
  "total_requests": 100,
  "avg_response_time": 30.5,
  "error_count": 2,
  "success_rate": 98
}
```

### 10. 清理过期Token
#### POST /v1/internal/services/cleanup-tokens
- **认证**：Bearer 服务JWT（需 internal:admin 权限）
- **响应示例**：
```json
{
  "message": "Expired tokens cleaned up successfully"
}
```

---

## /api/internal/ 路由权限说明

- 该路由下所有接口均需 Bearer 服务JWT 认证。
- 权限（Scope）控制：
  - user:read / user:write / tenant:read / auth:token / internal:admin 等
- 示例：
  - GET /api/internal/users 需 user:read
  - POST /api/internal/users 需 user:write
  - GET /api/internal/tenants 需 tenant:read
  - POST /api/internal/auth/token 需 auth:token
  - GET /api/internal/admin/services 需 internal:admin
- 若权限不足，返回 403 Forbidden。 