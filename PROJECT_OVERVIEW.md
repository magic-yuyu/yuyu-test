# IDaaS项目概览

## 🎯 项目简介

IDaaS（身份认证即服务）是一个基于Go语言构建的SaaS平台，提供完整的用户认证和管理功能。采用模块化单体架构，支持多租户，部署在Railway平台。

## 📁 目录结构

```
yuyu-test/
├── 📂 cmd/server/                  # 应用入口
│   ├── main.go                     # 主程序
│   └── main_test.go                # 测试文件
├── 📂 internal/                    # 内部包
│   ├── 📂 api/                     # API层
│   │   ├── 📂 handlers/            # 请求处理器
│   │   ├── 📂 middleware/          # 中间件
│   │   └── router.go               # 路由配置
│   ├── 📂 auth/                    # 认证模块
│   ├── 📂 config/                  # 配置管理
│   ├── 📂 store/database/          # 数据库层
│   ├── 📂 tenant/                  # 租户模块
│   └── 📂 user/                    # 用户模块
├── 📂 scripts/                     # 开发脚本
│   ├── start-dev.bat               # Windows启动脚本
│   ├── start-dev.sh                # Linux/Mac启动脚本
│   ├── migrate.bat                 # 数据库迁移
│   ├── dev-tools.bat               # 开发工具菜单
│   └── README.md                   # 脚本说明
├── 📂 config/                      # 配置文件
│   ├── .env.example                # 环境变量示例
│   ├── docker-compose.dev.yml      # 开发环境配置
│   └── README.md                   # 配置说明
├── 📂 docs/                        # 文档
│   ├── API_USAGE.md                # API使用指南
│   ├── 开发文档.md                 # 技术架构文档
│   └── README.md                   # 文档说明
├── 📂 migrations/                  # 数据库迁移
├── 📂 doc/                         # 原始开发文档
├── 📄 go.mod                       # Go模块文件
├── 📄 go.sum                       # 依赖校验
├── 📄 Dockerfile                   # Docker构建
├── 📄 docker-compose.yml           # Docker配置
├── 📄 railway.toml                 # Railway配置
├── 📄 sqlc.yaml                    # SQL代码生成配置
├── 📄 .cursorrules                 # Cursor AI规则
├── 📄 .dockerignore                # Docker忽略文件
└── 📄 README.md                    # 项目说明
```

## 🚀 快速开始

### 1. 环境准备
```bash
# 安装依赖
go mod tidy

# 启动PostgreSQL（Docker）
docker run -d --name postgres-idaas \
  -e POSTGRES_DB=idaas_dev \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=password \
  -p 5432:5432 postgres:15-alpine
```

### 2. 数据库迁移
```sql
-- 执行 migrations/0001_initial_schema.up.sql
```

### 3. 启动应用
```bash
# Windows
scripts\start-dev.bat

# Linux/Mac
./scripts/start-dev.sh
```

### 4. 验证运行
```bash
curl http://localhost:8080/health
```

## 🔧 开发工具

### 脚本工具
- `scripts/start-dev.bat` - 启动开发环境
- `scripts/dev-tools.bat` - 开发工具菜单
- `scripts/migrate.bat` - 数据库迁移

### 配置管理
- `config/.env.example` - 环境变量模板
- `config/docker-compose.dev.yml` - 开发环境Docker配置

### 文档资源
- `docs/API_USAGE.md` - API使用指南
- `docs/开发文档.md` - 技术架构文档

## 🏗️ 技术架构

### 核心特性
- **模块化单体** - 清晰的模块划分，便于扩展
- **多租户支持** - 严格的数据隔离
- **双重认证** - API密钥 + JWT令牌
- **类型安全** - 强类型数据库访问

### 技术栈
- **语言**: Go 1.18+
- **框架**: Gin Web框架
- **数据库**: PostgreSQL + pgx驱动
- **认证**: JWT + bcrypt
- **部署**: Railway平台

### API端点
- `POST /v1/tenants` - 创建租户
- `POST /v1/auth/register` - 用户注册
- `POST /v1/auth/login` - 用户登录
- `GET /v1/users/me` - 获取当前用户
- `GET /v1/users` - 获取用户列表

## 📚 文档导航

- **入门指南**: [README.md](README.md)
- **API文档**: [docs/API_USAGE.md](docs/API_USAGE.md)
- **技术架构**: [docs/开发文档.md](docs/开发文档.md)
- **脚本说明**: [scripts/README.md](scripts/README.md)
- **配置说明**: [config/README.md](config/README.md)

## 🔄 开发流程

1. **环境设置** - 使用 `scripts/start-dev.bat`
2. **代码开发** - 遵循 `.cursorrules` 规范
3. **测试验证** - 使用 `scripts/dev-tools.bat`
4. **部署上线** - 推送到Railway平台

## 📝 维护说明

- 新增功能时更新相应文档
- 保持代码与文档同步
- 遵循Go语言最佳实践
- 定期更新依赖版本 