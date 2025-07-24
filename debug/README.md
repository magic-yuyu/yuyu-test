# API调试工具使用指南

## 📁 目录说明

这个目录包含了IDaaS平台的API调试工具和测试用例。

## 🛠️ 文件说明

### `apitest.http`
- **用途**: 包含所有API端点的测试用例
- **格式**: REST Client格式（VS Code扩展）
- **内容**: 
  - 健康检查测试
  - 租户管理测试
  - 用户认证测试
  - 用户管理测试
  - 错误处理测试
  - 性能测试

## 🚀 快速开始

### 1. 安装VS Code扩展

```bash
# Windows
scripts\install-debug-tools.bat

# Linux/Mac
chmod +x scripts/install-debug-tools.sh
./scripts/install-debug-tools.sh
```

### 2. 启动IDaaS服务

```bash
# 启动开发环境
scripts\start-dev.bat

# 或手动启动
go run cmd/server/main.go
```

### 3. 开始调试

1. **打开VS Code**
2. **打开项目文件夹**
3. **打开 `debug/apitest.http` 文件**
4. **点击请求上方的"Send Request"按钮**

## 📋 测试用例说明

### 基础测试流程

1. **健康检查** → 确认服务运行
   ```http
   GET http://localhost:8080/health
   ```

2. **创建租户** → 获得API密钥
   ```http
   POST http://localhost:8080/v1/tenants
   Content-Type: application/json
   
   {
     "name": "测试应用"
   }
   ```

3. **用户注册** → 测试认证
   ```http
   POST http://localhost:8080/v1/auth/register
   Content-Type: application/json
   Authorization: Bearer {{apiKey}}
   
   {
     "email": "test@example.com",
     "password": "password123"
   }
   ```

4. **用户登录** → 获得JWT令牌
   ```http
   POST http://localhost:8080/v1/auth/login
   Content-Type: application/json
   Authorization: Bearer {{apiKey}}
   
   {
     "email": "test@example.com",
     "password": "password123"
   }
   ```

### 环境变量

在 `apitest.http` 文件中定义了以下变量：


**注意**: `apiKey` 需要从创建租户的响应中获取。

## 🔧 配置说明

### VS Code设置

在VS Code的 `settings.json` 中添加：

```json
{
  "rest-client.environmentVariables": {
    "local": {
      "baseUrl": "http://localhost:8080",
      "apiKey": "your_local_api_key"
    },
    "dev": {
      "baseUrl": "https://dev-api.idaas.com",
      "apiKey": "your_dev_api_key"
    }
  }
}
```

### 动态变量

可以使用响应中的值作为下一个请求的变量：

```http
### 登录并保存JWT令牌
POST {{baseUrl}}/v1/auth/login
Content-Type: application/json
Authorization: Bearer {{apiKey}}

{
  "email": "test@example.com",
  "password": "password123"
}

> {%
client.global.set("jwtToken", response.body.token);
%}

### 使用JWT令牌获取用户信息
GET {{baseUrl}}/v1/users/me
Authorization: Bearer {{jwtToken}}
```

## 🧪 测试类型

### 1. 功能测试
- ✅ 正常流程测试
- ✅ 参数验证测试
- ✅ 认证授权测试

### 2. 错误处理测试
- ❌ 无效认证测试
- ❌ 无效参数测试
- ❌ 资源不存在测试

### 3. 性能测试
- ⚡ 并发请求测试
- ⚡ 响应时间测试
- ⚡ 负载测试

## 📊 响应格式

### 成功响应
```json
{
  "id": "usr_abc123",
  "email": "test@example.com",
  "profile": {
    "name": "测试用户"
  },
  "created_at": "2024-01-01T00:00:00Z"
}
```

### 错误响应
```json
{
  "error": "错误描述"
}
```

## 🔍 调试技巧

### 1. 查看请求详情
- 在VS Code中，点击请求上方的"Send Request"
- 右侧会显示完整的请求和响应信息

### 2. 查看响应头
- 响应信息包含状态码、响应头、响应体
- 可以检查认证、缓存、内容类型等信息

### 3. 批量测试
- 使用 `###` 分隔符分隔多个请求
- 可以一次性执行多个测试用例

### 4. 环境切换
- 在VS Code右下角选择环境（local/dev/prod）
- 不同环境使用不同的配置

## 🚨 常见问题

### 1. 连接失败
```bash
# 检查服务是否运行
curl http://localhost:8080/health

# 检查端口是否被占用
netstat -an | findstr :8080  # Windows
lsof -i :8080                # Linux/Mac
```

### 2. 认证失败
- 检查API密钥是否正确
- 确认认证头格式：`Authorization: Bearer <api_key>`
- 验证JWT令牌是否有效

### 3. 请求格式错误
- 检查Content-Type：`Content-Type: application/json`
- 验证JSON格式是否正确
- 确认请求体字段是否完整

## 📚 相关文档

- [API文档](../docs/API_DOCUMENTATION.md)
- [调试工具安装指南](../docs/API_DEBUG_TOOLS.md)
- [项目概览](../PROJECT_OVERVIEW.md)

## 🤝 贡献指南

1. 添加新的测试用例到 `apitest.http`
2. 更新相关文档
3. 确保测试用例覆盖所有API端点
4. 包含正常和异常情况的测试

---

**提示**: 建议在开发过程中经常运行这些测试用例，确保API功能正常。 