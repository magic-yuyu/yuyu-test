# IDaaS身份认证即服务平台

这是一个使用Go语言构建的IDaaS（身份认证即服务）SaaS平台，采用模块化单体架构，部署在Railway上。

## 项目结构

```
yuyu-test/
├── cmd/server/main.go              # 应用主入口
├── internal/
│   ├── api/                        # API层
│   │   ├── handlers/               # 请求处理器
│   │   ├── middleware/             # 中间件
│   │   └── router.go               # 路由配置
│   ├── auth/                       # 认证模块
│   ├── config/                     # 配置管理
│   ├── store/database/             # 数据库层
│   ├── tenant/                     # 租户模块
│   └── user/                       # 用户模块
├── scripts/                        # 开发脚本
│   ├── start-dev.bat               # Windows启动脚本
│   ├── start-dev.sh                # Linux/Mac启动脚本
│   ├── migrate.bat                 # 数据库迁移脚本
│   ├── dev-tools.bat               # 开发工具菜单
│   └── README.md                   # 脚本说明
├── config/                         # 配置文件
│   ├── .env.example                # 环境变量示例
│   ├── docker-compose.dev.yml      # 开发环境Docker配置
│   └── README.md                   # 配置说明
├── docs/                           # 文档
│   ├── API_USAGE.md                # API使用指南
│   ├── 开发文档.md                 # 技术架构文档
│   └── README.md                   # 文档说明
├── migrations/                     # 数据库迁移
├── go.mod                          # Go模块文件
├── Dockerfile                      # Docker构建文件
├── docker-compose.yml              # Docker Compose配置
├── railway.toml                    # Railway配置
├── .dockerignore                   # Docker忽略文件
└── README.md                       # 项目说明
```

## 快速开始

### 使用Docker运行

1. 构建并运行容器：
```bash
docker-compose up --build
```

2. 访问应用：
   - 打开浏览器访问 http://localhost:8080

### 本地开发

#### 快速开始

1. **安装依赖**
   ```bash
   go mod tidy
   ```

2. **设置数据库**
   
   **选项A：使用Docker（推荐）**
   ```bash
   docker run -d --name postgres-idaas \
     -e POSTGRES_DB=idaas_dev \
     -e POSTGRES_USER=postgres \
     -e POSTGRES_PASSWORD=password \
     -p 5432:5432 postgres:15-alpine
   ```
   
   **选项B：本地PostgreSQL**
   - 安装PostgreSQL并创建数据库 `idaas_dev`

3. **执行数据库迁移**
   ```sql
   -- 连接到idaas_dev数据库，执行migrations/0001_initial_schema.up.sql
   ```

4. **启动应用**
   
   **Windows用户：**
   ```cmd
   scripts\start-dev.bat
   ```
   
   **Linux/Mac用户：**
   ```bash
   chmod +x scripts/start-dev.sh
   ./scripts/start-dev.sh
   ```
   
   **手动启动：**
   ```bash
   set DATABASE_URL=postgresql://postgres:password@localhost:5432/idaas_dev
   set JWT_SECRET=your-development-secret-key
   set PORT=8080
   go run cmd/server/main.go
   ```

5. **验证运行**
   - 健康检查：http://localhost:8080/health
   - API文档：http://localhost:8080/v1

#### 开发工具

运行 `scripts\dev-tools.bat` 获取交互式开发菜单。

### 主要环境变量

- `DATABASE_URL`：PostgreSQL连接字符串
- `JWT_ALGORITHM`：JWT签名算法，可选`HS256`（对称，默认）或`RS256`（非对称）。
- `JWT_USER_SECRET_KEY`：用户JWT签名密钥（HS256时必填，强烈建议32字节以上，生产环境必须≥32字符，否则应用无法启动）
- `JWT_SERVICE_SECRET_KEY`：服务JWT签名密钥（HS256时必填，强烈建议32字节以上，生产环境必须≥32字符，否则应用无法启动）
- `JWT_USER_PRIVATE_KEY`/`JWT_USER_PUBLIC_KEY`：用户JWT私钥/公钥（RS256时必填，PEM内容或文件路径）
- `JWT_SERVICE_PRIVATE_KEY`/`JWT_SERVICE_PUBLIC_KEY`：服务JWT私钥/公钥（RS256时必填，PEM内容或文件路径）
- `USER_TOKEN_EXPIRATION`：用户JWT有效期（单位：秒，默认3600=1小时）
- `SERVICE_TOKEN_EXPIRATION`：服务JWT有效期（单位：秒，默认300=5分钟）
- `PORT`：服务监听端口
- `GO_ENV`：运行环境

### JWT 密钥生成与配置检测

#### 生成RSA密钥对（适用于RS256）
```sh
./scripts/gen_rsa_keys.sh jwt_rsa_private.pem jwt_rsa_public.pem
```
生成后，将密钥内容填入环境变量：
- `JWT_USER_PRIVATE_KEY`、`JWT_USER_PUBLIC_KEY`
- `JWT_SERVICE_PRIVATE_KEY`、`JWT_SERVICE_PUBLIC_KEY`

#### 生成ES256密钥对（适用于ES256）
```sh
./scripts/gen_es256_keys.sh jwt_es256_private.pem jwt_es256_public.pem
```
生成后，将密钥内容填入环境变量：
- `JWT_USER_PRIVATE_KEY`、`JWT_USER_PUBLIC_KEY`
- `JWT_SERVICE_PRIVATE_KEY`、`JWT_SERVICE_PUBLIC_KEY`

#### 检查JWT相关环境变量配置
```sh
source .env
./scripts/check_jwt_config.sh
```

### 环境变量配置示例

```env
# 数据库
DATABASE_URL=postgres://user:pass@localhost:5432/dbname

# JWT算法
JWT_ALGORITHM=RS256 # 可选: HS256/RS256/ES256

# HS256密钥
JWT_USER_SECRET_KEY=your_user_secret_key_32bytes_min
JWT_SERVICE_SECRET_KEY=your_service_secret_key_32bytes_min

# RS256/ES256密钥内容（PEM）
JWT_USER_PRIVATE_KEY="$(cat jwt_rsa_private.pem)"
JWT_USER_PUBLIC_KEY="$(cat jwt_rsa_public.pem)"
JWT_SERVICE_PRIVATE_KEY="$(cat jwt_rsa_private.pem)"
JWT_SERVICE_PUBLIC_KEY="$(cat jwt_rsa_public.pem)"

# 令牌有效期（秒）
USER_TOKEN_EXPIRATION=3600
SERVICE_TOKEN_EXPIRATION=3600

# 端口与环境
PORT=8080
GO_ENV=development
```

### 常见问题 FAQ

- Q: 如何生成和配置JWT密钥？
  A: 见上方密钥脚本说明，生成后将内容填入对应环境变量。
- Q: 配置检测脚本报错怎么办？
  A: 检查环境变量是否设置、密钥长度是否符合要求，生产环境HS256密钥需32字节以上。
- Q: 如何切换JWT算法？
  A: 修改`JWT_ALGORITHM`为HS256/RS256/ES256，并配置对应密钥。
- Q: Railway等云平台如何安全注入密钥？
  A: 推荐将PEM内容直接粘贴到环境变量，或使用密钥文件路径（需支持挂载）。

### 测试用例运行

本项目所有测试用例均位于各模块`*_test.go`文件，使用Go官方`testing`框架。

#### 运行全部测试
```sh
go test ./...
```

#### 数据库相关测试
- 推荐使用`testcontainers-go`自动拉起PostgreSQL容器，确保测试环境隔离。
- 可通过设置`TEST_DATABASE_URL`环境变量指定测试数据库。

```sh
export TEST_DATABASE_URL=postgres://user:pass@localhost:5433/testdb
go test ./internal/store/...
```

## API端点

### 租户管理
- `POST /v1/tenants` - 创建新租户
- `GET /v1/tenants/:id` - 获取租户信息

### 用户认证
- `POST /v1/auth/register` - 用户注册（需要API密钥）
- `POST /v1/auth/login` - 用户登录（需要API密钥）

### 用户管理
- `GET /v1/users/me` - 获取当前用户信息（需要JWT）
- `GET /v1/users` - 获取租户下所有用户（需要API密钥）
- `GET /v1/users/:id` - 获取指定用户信息（需要API密钥）

## 开发命令

- `go run cmd/server/main.go` - 运行应用
- `go build ./cmd/server` - 构建应用
- `go test ./...` - 运行测试
- `go mod tidy` - 整理依赖

## Docker命令

- `docker-compose up` - 启动服务
- `docker-compose down` - 停止服务
- `docker-compose up --build` - 重新构建并启动
- `docker build -t yuyu-test .` - 构建镜像 