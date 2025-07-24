# API调试工具安装指南

## 📋 目录
- [工具推荐](#工具推荐)
- [VS Code扩展](#vs-code扩展)
- [独立工具](#独立工具)
- [配置说明](#配置说明)
- [使用技巧](#使用技巧)

## 工具推荐

### 🥇 推荐工具

#### 1. VS Code REST Client (最推荐)
- **优点**: 集成在VS Code中，支持`.http`文件，语法高亮，变量支持
- **安装**: VS Code扩展市场搜索"REST Client"
- **文件**: `debug/apitest.http`

#### 2. Postman
- **优点**: 功能强大，界面友好，支持团队协作
- **缺点**: 需要单独安装，免费版有限制
- **下载**: https://www.postman.com/downloads/

#### 3. Insomnia
- **优点**: 轻量级，界面简洁，开源
- **缺点**: 功能相对简单
- **下载**: https://insomnia.rest/download

#### 4. curl (命令行)
- **优点**: 系统自带，脚本友好
- **缺点**: 命令行操作，不够直观

## VS Code扩展

### 安装REST Client扩展

1. **打开VS Code**
2. **按 `Ctrl+Shift+X` 打开扩展面板**
3. **搜索 "REST Client"**
4. **点击安装 "REST Client" 扩展**

### 配置VS Code设置

在VS Code的`settings.json`中添加以下配置：

```json
{
  "rest-client.environmentVariables": {
    "$shared": {
      "version": "1.0.0"
    },
    "local": {
      "baseUrl": "http://localhost:8080",
      "apiKey": "your_local_api_key"
    },
    "dev": {
      "baseUrl": "https://dev-api.idaas.com",
      "apiKey": "your_dev_api_key"
    },
    "prod": {
      "baseUrl": "https://api.idaas.com",
      "apiKey": "your_prod_api_key"
    }
  },
  "rest-client.defaultHeaders": {
    "Content-Type": "application/json",
    "User-Agent": "IDaaS-API-Client/1.0"
  }
}
```

### 使用REST Client

1. **打开文件**: `debug/apitest.http`
2. **选择环境**: 右下角选择环境（local/dev/prod）
3. **发送请求**: 点击请求上方的"Send Request"
4. **查看响应**: 右侧会显示响应结果

## 独立工具

### Postman安装配置

#### 1. 下载安装
```bash
# Windows
# 访问 https://www.postman.com/downloads/ 下载安装包

# macOS
brew install --cask postman

# Linux
# 下载AppImage或使用包管理器
```

#### 2. 创建集合
1. 打开Postman
2. 点击"New" → "Collection"
3. 命名为"IDaaS API"
4. 设置环境变量

#### 3. 环境变量配置
```json
{
  "baseUrl": "http://localhost:8080",
  "apiKey": "your_api_key",
  "jwtToken": "your_jwt_token"
}
```

#### 4. 导入API文档
```bash
# 如果有OpenAPI/Swagger文档
# 可以导入到Postman中自动生成请求
```

### Insomnia安装配置

#### 1. 下载安装
```bash
# Windows/macOS/Linux
# 访问 https://insomnia.rest/download 下载
```

#### 2. 创建项目
1. 打开Insomnia
2. 创建新项目"IDaaS"
3. 设置环境变量

#### 3. 环境配置
```json
{
  "baseUrl": "http://localhost:8080",
  "apiKey": "your_api_key"
}
```

## 配置说明

### 环境变量

在`debug/apitest.http`文件中，我们定义了以下变量：

```http
@baseUrl = http://localhost:8080
@apiKey = r8_OjJpSnPlnFs06Gjyhr8vLI1pHuLdKTY0DGv6d
```

### 认证方式

#### API密钥认证
```http
Authorization: Bearer {{apiKey}}
```

#### JWT认证
```http
Authorization: Bearer {{jwtToken}}
```

### 请求格式

#### JSON请求体
```json
{
  "email": "test@example.com",
  "password": "password123",
  "profile": {
    "name": "测试用户",
    "role": "user"
  }
}
```

#### 查询参数
```http
GET {{baseUrl}}/v1/users?page=1&limit=10
```

## 使用技巧

### 1. 变量使用

#### 动态变量
```http
### 使用响应中的值作为下一个请求的变量
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

### 使用登录获得的JWT令牌
GET {{baseUrl}}/v1/users/me
Authorization: Bearer {{jwtToken}}
```

#### 环境切换
```http
### 开发环境
@baseUrl = http://localhost:8080

### 生产环境
@baseUrl = https://api.idaas.com
```

### 2. 批量测试

#### 使用分隔符
```http
### 请求1
GET {{baseUrl}}/health

###

### 请求2
POST {{baseUrl}}/v1/tenants
Content-Type: application/json

{
  "name": "测试租户"
}
```

### 3. 错误处理测试

#### 测试各种错误情况
```http
### 测试无效认证
GET {{baseUrl}}/v1/users/me
Authorization: Bearer invalid_token

### 测试无效请求体
POST {{baseUrl}}/v1/auth/register
Content-Type: application/json
Authorization: Bearer {{apiKey}}

{
  "email": "invalid-email"
}
```

### 4. 性能测试

#### 并发请求
```http
### 并发健康检查
GET {{baseUrl}}/health

###

GET {{baseUrl}}/health

###

GET {{baseUrl}}/health
```

## 调试流程

### 1. 基础测试
1. 健康检查 → 确认服务运行
2. 创建租户 → 获得API密钥
3. 用户注册 → 测试认证
4. 用户登录 → 获得JWT令牌

### 2. 功能测试
1. 获取用户信息 → 测试JWT认证
2. 获取用户列表 → 测试API密钥认证
3. 错误处理 → 测试各种错误情况

### 3. 集成测试
1. 完整流程测试
2. 边界条件测试
3. 性能压力测试

## 常见问题

### 1. 连接失败
```bash
# 检查服务是否运行
curl http://localhost:8080/health

# 检查端口是否被占用
netstat -an | findstr :8080
```

### 2. 认证失败
```bash
# 检查API密钥格式
# 确保使用正确的认证头格式
Authorization: Bearer your_api_key
```

### 3. 请求格式错误
```bash
# 检查Content-Type
Content-Type: application/json

# 检查JSON格式
# 使用在线JSON验证工具
```

## 最佳实践

### 1. 组织测试用例
- 按功能模块分组
- 使用清晰的注释
- 保持测试用例的独立性

### 2. 环境管理
- 使用环境变量
- 区分开发/测试/生产环境
- 定期更新API密钥

### 3. 文档同步
- 保持测试用例与API文档同步
- 记录测试结果
- 更新错误处理用例

### 4. 自动化测试
- 集成到CI/CD流程
- 使用脚本批量执行
- 生成测试报告

---

## 快速开始

1. **安装VS Code REST Client扩展**
2. **打开 `debug/apitest.http` 文件**
3. **选择环境（local/dev/prod）**
4. **开始测试API接口**

更多信息请参考：
- [REST Client文档](https://marketplace.visualstudio.com/items?itemName=humao.rest-client)
- [Postman文档](https://learning.postman.com/)
- [Insomnia文档](https://docs.insomnia.rest/) 