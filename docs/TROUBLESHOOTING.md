# 故障排除指南

## 🚨 常见问题及解决方案

### 1. 启动错误：CreateFile cmd/server/main.go: The system cannot find the path specified

**问题原因：**
- 脚本编码问题导致路径解析错误
- 环境变量设置不正确

**解决方案：**
```cmd
# 使用修复后的启动脚本
scripts\start-dev.bat

# 或者手动设置环境变量
set DATABASE_URL=postgresql://postgres:password@localhost:5432/idaas_dev
set JWT_SECRET=your-development-secret-key
set PORT=8080
set GO_ENV=development
go run cmd/server/main.go
```

### 2. 配置错误：DATABASE_URL environment variable is required

**问题原因：**
- 环境变量未正确设置
- 环境变量包含空格

**解决方案：**
```cmd
# 使用测试脚本验证环境变量
scripts\test-env.bat

# 检查环境变量是否正确设置
echo %DATABASE_URL%
echo %JWT_SECRET%
echo %PORT%
echo %GO_ENV%
```

### 3. 数据库连接错误：dial tcp 127.0.0.1:5432: connectex: No connection could be made

**问题原因：**
- PostgreSQL数据库未启动
- 数据库端口被占用
- 数据库配置错误

**解决方案：**

#### 使用Docker启动数据库
```cmd
# 启动PostgreSQL容器
docker run -d --name postgres-idaas \
  -e POSTGRES_DB=idaas_dev \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=password \
  -p 5432:5432 postgres:15-alpine

# 等待数据库启动
timeout /t 5
```

#### 使用完整启动脚本
```cmd
# 自动启动数据库和应用
scripts\start-complete.bat
```

#### 检查数据库状态
```cmd
# 检查容器是否运行
docker ps | findstr postgres

# 检查端口是否被占用
netstat -an | findstr :5432
```

### 4. 数据库迁移错误：relation "tenants" does not exist

**问题原因：**
- 数据库表未创建
- 迁移脚本未执行

**解决方案：**

#### 手动执行迁移
```sql
-- 连接到数据库后执行
-- 创建租户表
CREATE TABLE IF NOT EXISTS tenants (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    api_secret_key_hash VARCHAR(255) UNIQUE NOT NULL,
    api_public_key VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- 创建用户表
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(255) PRIMARY KEY,
    tenant_id VARCHAR(255) NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    email VARCHAR(255) NOT NULL,
    hashed_password VARCHAR(255),
    profile JSONB,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    UNIQUE(tenant_id, email)
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_users_tenant_id ON users(tenant_id);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_tenants_public_key ON tenants(api_public_key);
CREATE INDEX IF NOT EXISTS idx_tenants_secret_key_hash ON tenants(api_secret_key_hash);
```

#### 使用数据库客户端
- **pgAdmin**: 图形化界面
- **DBeaver**: 跨平台数据库工具
- **命令行**: `psql -h localhost -p 5432 -U postgres -d idaas_dev`

### 5. 编译错误：undefined: database.Querier

**问题原因：**
- SQLC生成的代码不存在
- 导入路径错误

**解决方案：**
```cmd
# 生成数据库代码
scripts\sqlc.bat

# 或者手动生成
sqlc generate

# 检查生成的文件
dir internal\store\database\
```

### 6. 依赖错误：could not import github.com/gin-gonic/gin

**问题原因：**
- Go模块依赖未下载
- 网络连接问题

**解决方案：**
```cmd
# 下载依赖
go mod tidy

# 清理模块缓存
go clean -modcache

# 重新下载
go mod download
```

### 7. 端口占用错误：address already in use

**问题原因：**
- 8080端口被其他程序占用
- 之前的实例未正确关闭

**解决方案：**
```cmd
# 查找占用端口的进程
netstat -ano | findstr :8080

# 结束进程（替换PID为实际进程ID）
taskkill /f /pid <PID>

# 或者使用不同端口
set PORT=8081
go run cmd/server/main.go
```

## 🔧 调试技巧

### 1. 启用详细日志
```go
// 在main.go中添加
slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
    Level: slog.LevelDebug,
})))
```

### 2. 测试数据库连接
```cmd
# 使用psql测试连接
psql -h localhost -p 5432 -U postgres -d idaas_dev

# 或者使用Docker
docker exec -it postgres-idaas psql -U postgres -d idaas_dev
```

### 3. 验证环境变量
```cmd
# 创建测试脚本
echo %DATABASE_URL%
echo %JWT_SECRET%
echo %PORT%
echo %GO_ENV%
```

### 4. 检查Go版本
```cmd
go version
go env
```

## 📞 获取帮助

如果以上解决方案无法解决问题，请：

1. **检查日志** - 查看详细的错误信息
2. **验证环境** - 确保所有依赖都已正确安装
3. **查看文档** - 参考 `docs/` 目录下的相关文档
4. **使用测试脚本** - 运行 `scripts\test-env.bat` 验证环境

## 🎯 快速诊断

运行以下命令进行快速诊断：

```cmd
# 1. 检查Go环境
go version
go mod tidy

# 2. 检查数据库
docker ps | findstr postgres

# 3. 测试环境变量
scripts\test-env.bat

# 4. 生成代码
scripts\sqlc.bat

# 5. 启动应用
scripts\start-complete.bat
``` 