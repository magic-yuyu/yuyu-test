# 配置文件

这个文件夹包含了IDaaS平台的各种配置文件。

## 文件说明

### 环境配置

- **.env.example** - 环境变量示例文件
  - 复制为 `.env` 并根据实际情况修改
  - 包含数据库连接、JWT密钥等配置

### Docker配置

- **docker-compose.dev.yml** - 开发环境Docker Compose配置
  - 包含PostgreSQL数据库和应用服务
  - 适用于本地开发环境

## 使用方法

### 环境变量配置

1. 复制环境变量示例文件：
   ```bash
   cp config/.env.example config/.env
   ```

2. 编辑 `.env` 文件，设置实际的环境变量：
   ```bash
   DATABASE_URL=postgresql://username:password@localhost:5432/idaas_dev
   JWT_SECRET=your-secret-key
   PORT=8080
   GO_ENV=development
   ```

### Docker开发环境

使用Docker Compose启动完整的开发环境：

```bash
docker-compose -f config/docker-compose.dev.yml up --build
```

这将启动：
- PostgreSQL数据库（端口5432）
- IDaaS应用（端口8080）

## 配置说明

### 数据库配置

- **DATABASE_URL**: PostgreSQL连接字符串
- 格式：`postgresql://username:password@host:port/database`

### 安全配置

- **JWT_ALGORITHM**: JWT签名算法，可选`HS256`（对称，默认）或`RS256`（非对称）。
- **JWT_USER_SECRET_KEY**: 用户JWT令牌签名密钥（HS256时必填，优先于JWT_SECRET，强烈建议32字节以上，生产环境必须≥32字符，否则应用无法启动）
- **JWT_SERVICE_SECRET_KEY**: 服务JWT令牌签名密钥（HS256时必填，强烈建议32字节以上，生产环境必须≥32字符，否则应用无法启动）
- **JWT_USER_PRIVATE_KEY**/**JWT_USER_PUBLIC_KEY**: 用户JWT私钥/公钥（RS256时必填，PEM内容或文件路径）
- **JWT_SERVICE_PRIVATE_KEY**/**JWT_SERVICE_PUBLIC_KEY**: 服务JWT私钥/公钥（RS256时必填，PEM内容或文件路径）
- **JWT_SECRET**: 兼容老版本，若未设置上述密钥则作为默认密钥
- **USER_TOKEN_EXPIRATION**: 用户JWT有效期（单位：秒，默认3600=1小时）
- **SERVICE_TOKEN_EXPIRATION**: 服务JWT有效期（单位：秒，默认300=5分钟）

### 服务配置

- **PORT**: 应用监听端口（默认8080）
- **GO_ENV**: 运行环境（development/production）

## 注意事项

1. 生产环境请修改默认密码和密钥
2. 数据库连接信息请根据实际部署环境调整
3. 不要将包含敏感信息的配置文件提交到版本控制 