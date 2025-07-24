# 内部服务管理 (Internal Service Management)

## 概述

内部服务管理是IDaaS平台的核心功能之一，用于管理内部服务之间的认证和授权。它提供了完整的M2M（Machine-to-Machine）认证解决方案，包括服务注册、权限管理、访问控制等功能。

## 核心功能

### 1. 服务注册与管理
- **服务注册**: 注册新的内部服务，自动生成客户端ID和密钥
- **服务列表**: 查看所有已注册的内部服务
- **服务状态管理**: 激活/停用内部服务
- **服务信息更新**: 更新服务名称和描述

### 2. 权限管理
- **权限定义**: 预定义系统权限（如user:read, user:write等）
- **权限授权**: 为内部服务授权特定权限
- **权限撤销**: 撤销内部服务的特定权限
- **权限检查**: 检查内部服务是否拥有特定权限

### 3. 认证与令牌管理
- **JWT令牌颁发**: 为内部服务颁发JWT访问令牌
- **令牌验证**: 验证JWT令牌的有效性
- **令牌撤销**: 撤销已颁发的令牌
- **令牌清理**: 自动清理过期的令牌

### 4. 访问控制
- **基于权限的访问控制**: 根据服务权限控制API访问
- **中间件支持**: 提供多种认证中间件
- **复合权限**: 支持任意权限、所有权限的组合
- **可选认证**: 支持可选的认证机制

### 5. 监控与日志
- **访问日志**: 记录所有内部服务的API访问
- **统计信息**: 提供访问统计和性能指标
- **审计追踪**: 记录权限变更和操作历史

## 数据库设计

### 核心表结构

#### 1. internal_clients (内部服务表)
```sql
CREATE TABLE internal_clients (
    client_id VARCHAR(255) PRIMARY KEY,
    client_secret_hash VARCHAR(255) NOT NULL,
    service_name VARCHAR(255) NOT NULL,
    description TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);
```

#### 2. scopes (权限表)
```sql
CREATE TABLE scopes (
    id SERIAL PRIMARY KEY,
    scope_name VARCHAR(255) UNIQUE NOT NULL,
    description TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);
```

#### 3. client_scopes (服务权限关联表)
```sql
CREATE TABLE client_scopes (
    client_id VARCHAR(255) NOT NULL REFERENCES internal_clients(client_id) ON DELETE CASCADE,
    scope_id INT NOT NULL REFERENCES scopes(id) ON DELETE CASCADE,
    granted_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    granted_by VARCHAR(255),
    PRIMARY KEY (client_id, scope_id)
);
```

#### 4. service_access_logs (访问日志表)
```sql
CREATE TABLE service_access_logs (
    id SERIAL PRIMARY KEY,
    client_id VARCHAR(255) NOT NULL REFERENCES internal_clients(client_id) ON DELETE CASCADE,
    endpoint VARCHAR(255) NOT NULL,
    method VARCHAR(10) NOT NULL,
    status_code INT NOT NULL,
    response_time_ms INT,
    ip_address INET,
    user_agent TEXT,
    request_body TEXT,
    response_body TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);
```

#### 5. service_tokens (服务令牌表)
```sql
CREATE TABLE service_tokens (
    id SERIAL PRIMARY KEY,
    client_id VARCHAR(255) NOT NULL REFERENCES internal_clients(client_id) ON DELETE CASCADE,
    token_hash VARCHAR(255) NOT NULL,
    scopes TEXT[] NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    is_revoked BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);
```

## API接口

### 1. 服务注册

#### 注册新服务
```http
POST /v1/internal/services/register
Content-Type: application/json

{
    "client_id": "user-service",
    "service_name": "用户服务",
    "description": "负责用户管理的内部服务"
}
```

**响应示例:**
```json
{
    "client_id": "user-service",
    "client_secret": "generated-secret-key",
    "service_name": "用户服务",
    "description": "负责用户管理的内部服务",
    "created_at": "2024-01-01T00:00:00Z",
    "message": "Service registered successfully",
    "warning": "Please save the client_secret securely. It will not be shown again."
}
```

### 2. 服务认证

#### 获取访问令牌
```http
POST /v1/internal/services/authenticate
Content-Type: application/json

{
    "client_id": "user-service",
    "client_secret": "your-client-secret"
}
```

**响应示例:**
```json
{
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "token_type": "Bearer",
    "expires_in": 86400,
    "scopes": ["user:read", "user:write"]
}
```

### 3. 权限管理

#### 授权权限
```http
POST /v1/internal/services/grant-scope
Content-Type: application/json
Authorization: Bearer <admin-token>

{
    "client_id": "user-service",
    "scope_name": "user:read",
    "granted_by": "admin"
}
```

#### 撤销权限
```http
POST /v1/internal/services/revoke-scope
Content-Type: application/json
Authorization: Bearer <admin-token>

{
    "client_id": "user-service",
    "scope_name": "user:delete"
}
```

#### 检查权限
```http
POST /v1/internal/services/check-permission
Content-Type: application/json
Authorization: Bearer <access-token>

{
    "client_id": "user-service",
    "scope_name": "user:read"
}
```

### 4. 服务管理

#### 列出所有服务
```http
GET /v1/internal/services
Authorization: Bearer <access-token>
```

#### 获取访问日志
```http
GET /v1/internal/services/{client_id}/logs?limit=50&offset=0
Authorization: Bearer <access-token>
```

#### 获取统计信息
```http
GET /v1/internal/services/{client_id}/statistics?since=24h
Authorization: Bearer <access-token>
```

## 中间件使用

### 1. 基础认证中间件

```go
// 要求认证
router.Use(internalAuthMiddleware.RequireAuth())

// 要求特定权限
router.Use(internalAuthMiddleware.RequireScope("user:read"))

// 要求任意一个权限
router.Use(internalAuthMiddleware.RequireAnyScope("user:read", "tenant:read"))

// 要求所有权限
router.Use(internalAuthMiddleware.RequireAllScopes("user:read", "user:write"))

// 可选认证
router.Use(internalAuthMiddleware.OptionalAuth())
```

### 2. 在路由中使用

```go
// 用户管理API（需要user:read权限）
internalUsers := router.Group("/api/internal/users")
internalUsers.Use(internalAuthMiddleware.RequireScope("user:read"))
{
    internalUsers.GET("", userHandler.GetUsers)
    internalUsers.GET("/:id", userHandler.GetUser)
}

// 用户写入API（需要user:write权限）
internalUserWrite := router.Group("/api/internal/users")
internalUserWrite.Use(internalAuthMiddleware.RequireScope("user:write"))
{
    internalUserWrite.POST("", userHandler.CreateUser)
    internalUserWrite.PUT("/:id", userHandler.UpdateUser)
}
```

## 预定义权限

系统预定义了以下权限：

| 权限名称         | 描述               |
| ---------------- | ------------------ |
| `user:read`      | 读取用户信息       |
| `user:write`     | 创建和更新用户信息 |
| `user:delete`    | 删除用户           |
| `tenant:read`    | 读取租户信息       |
| `tenant:write`   | 创建和更新租户信息 |
| `tenant:delete`  | 删除租户           |
| `auth:token`     | 生成认证令牌       |
| `auth:validate`  | 验证令牌           |
| `internal:admin` | 内部服务管理权限   |

## 使用示例

### 1. 完整的服务注册和认证流程

```bash
# 1. 注册新服务
curl -X POST http://localhost:8080/v1/internal/services/register \
  -H "Content-Type: application/json" \
  -d '{
    "client_id": "my-service",
    "service_name": "我的服务",
    "description": "示例内部服务"
  }'

# 2. 使用返回的client_secret进行认证
curl -X POST http://localhost:8080/v1/internal/services/authenticate \
  -H "Content-Type: application/json" \
  -d '{
    "client_id": "my-service",
    "client_secret": "returned-secret"
  }'

# 3. 使用access_token访问API
curl -X GET http://localhost:8080/api/internal/users \
  -H "Authorization: Bearer your-access-token"
```

### 2. 权限管理示例

```bash
# 为服务授权权限
curl -X POST http://localhost:8080/v1/internal/services/grant-scope \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer admin-token" \
  -d '{
    "client_id": "my-service",
    "scope_name": "user:read",
    "granted_by": "admin"
  }'

# 检查权限
curl -X POST http://localhost:8080/v1/internal/services/check-permission \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer access-token" \
  -d '{
    "client_id": "my-service",
    "scope_name": "user:read"
  }'
```

### 3. 监控和日志

```bash
# 获取访问日志
curl -X GET "http://localhost:8080/v1/internal/services/my-service/logs?limit=10&offset=0" \
  -H "Authorization: Bearer access-token"

# 获取统计信息
curl -X GET "http://localhost:8080/v1/internal/services/my-service/statistics?since=24h" \
  -H "Authorization: Bearer access-token"
```

## 服务Token获取接口对比

### 1. /oauth/token （推荐对接标准OAuth2客户端/三方平台）
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

### 2. /v1/internal/services/authenticate （推荐平台内部服务、自动化脚本）
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

---

## 最佳实践
- 客户端密钥（client_secret）务必妥善保存，仅注册时返回一次。
- access_token（JWT）有效期短，建议定期刷新。
- 权限（scope）需在服务注册后由管理员分配。
- 推荐使用HTTPS保障传输安全。

## 安全考虑

### 1. 密钥管理
- 客户端密钥使用bcrypt进行哈希存储
- 密钥仅在注册时返回一次，后续无法查看
- 建议定期轮换客户端密钥

### 2. 令牌安全
- JWT令牌有效期为24小时
- 支持令牌撤销机制
- 自动清理过期令牌

### 3. 权限控制
- 基于最小权限原则设计
- 支持细粒度的权限控制
- 提供权限审计日志

### 4. 访问监控
- 记录所有API访问日志
- 提供访问统计和异常检测
- 支持IP地址和用户代理记录

## 最佳实践

### 1. 服务命名
- 使用有意义的服务名称
- 遵循统一的命名规范
- 避免使用敏感信息

### 2. 权限管理
- 遵循最小权限原则
- 定期审查和清理权限
- 记录权限变更历史

### 3. 监控告警
- 设置访问频率限制
- 监控异常访问模式
- 配置安全事件告警

### 4. 文档维护
- 及时更新服务文档
- 记录权限使用说明
- 提供故障排查指南

## 故障排查

### 常见问题

1. **认证失败**
   - 检查客户端ID和密钥是否正确
   - 确认服务是否已激活
   - 验证令牌是否过期

2. **权限不足**
   - 检查服务是否拥有所需权限
   - 确认权限是否已正确授权
   - 验证权限是否已激活

3. **令牌无效**
   - 检查令牌格式是否正确
   - 确认令牌是否已撤销
   - 验证令牌签名是否有效

### 调试工具

1. **日志查看**
   ```bash
   # 查看服务访问日志
   curl -X GET "http://localhost:8080/v1/internal/services/{client_id}/logs"
   ```

2. **权限检查**
   ```bash
   # 检查服务权限
   curl -X POST http://localhost:8080/v1/internal/services/check-permission
   ```

3. **令牌验证**
   ```bash
   # 验证令牌
   curl -X POST http://localhost:8080/v1/internal/services/validate-token
   ```

## 扩展功能

### 1. 自定义权限
- 支持添加自定义权限
- 提供权限模板管理
- 支持权限继承机制

### 2. 高级监控
- 实时访问监控
- 性能指标分析
- 安全事件检测

### 3. 集成支持
- 支持OAuth2集成
- 提供SDK和客户端库
- 支持多种编程语言

### 4. 自动化管理
- 自动权限分配
- 智能访问控制
- 自动化运维工具 