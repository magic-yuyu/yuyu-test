# SQLC 使用指南

## 📋 概述

SQLC 是一个代码生成工具，它从SQL查询文件自动生成类型安全的Go代码，避免了手写ORM代码的繁琐。

## 🏗️ 工作原理

### 1. 输入文件
- **SQL查询文件** (`internal/store/queries/`)
  - `tenant.sql` - 租户相关查询
  - `user.sql` - 用户相关查询
- **数据库架构** (`migrations/`)
  - `0001_initial_schema.up.sql` - 表结构定义

### 2. 输出文件
- **models.go** - 数据模型结构体
- **querier.go** - 查询接口定义
- **queries.go** - 查询方法实现

## 📁 配置文件

### sqlc.yaml 配置说明

```yaml
version: "2"                    # sqlc版本
sql:
  - engine: "postgresql"        # 数据库类型
    queries: "internal/store/queries/"    # SQL文件目录
    schema: "migrations/"       # 架构文件目录
    gen:
      go:
        package: "database"     # 生成的包名
        out: "internal/store/database"    # 输出目录
        emit_json_tags: true    # 生成JSON标签
        emit_prepared_queries: false      # 预处理查询
        emit_interface: true    # 生成接口
        emit_exact_table_names: false     # 精确表名
        emit_empty_slices: true # 空切片处理
```

## 🛠️ 使用方法

### 安装 sqlc

```bash
# 安装sqlc
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# 验证安装
sqlc version
```

### 基本命令

```bash
# 生成代码
sqlc generate

# 验证SQL语法
sqlc vet

# 查看配置
sqlc config
```

### 使用脚本

```cmd
# Windows
scripts\sqlc.bat

# 选择操作：
# 1. 生成代码
# 2. 验证SQL语法
# 3. 查看配置
# 4. 清理生成文件
```

## 📝 SQL 查询语法

### 查询类型

```sql
-- 返回单条记录
-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- 返回多条记录
-- name: ListUsers :many
SELECT * FROM users ORDER BY created_at;

-- 执行操作（不返回数据）
-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;
```

### 参数绑定

```sql
-- 使用 $1, $2, $3... 进行参数绑定
-- name: CreateUser :one
INSERT INTO users (id, tenant_id, email, hashed_password, profile)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;
```

### 生成的Go代码

```go
// 自动生成的接口
type Querier interface {
    CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
    GetUserByID(ctx context.Context, id string) (User, error)
    ListUsers(ctx context.Context) ([]User, error)
    DeleteUser(ctx context.Context, id string) error
}

// 自动生成的参数结构
type CreateUserParams struct {
    ID             string       `json:"id"`
    TenantID       string       `json:"tenant_id"`
    Email          string       `json:"email"`
    HashedPassword *string      `json:"hashed_password"`
    Profile        pgtype.JSONB `json:"profile"`
}
```

## 🔄 开发流程

### 1. 修改数据库架构
```sql
-- 在 migrations/ 中添加新的迁移文件
-- 例如：0002_add_user_profile.up.sql
ALTER TABLE users ADD COLUMN profile JSONB;
```

### 2. 添加SQL查询
```sql
-- 在 internal/store/queries/ 中添加查询
-- name: UpdateUserProfile :one
UPDATE users SET profile = $2 WHERE id = $1 RETURNING *;
```

### 3. 生成代码
```bash
sqlc generate
```

### 4. 使用生成的代码
```go
// 在业务逻辑中使用
user, err := queries.UpdateUserProfile(ctx, userID, profile)
if err != nil {
    return err
}
```

## 🎯 优势

### 类型安全
- 编译时检查SQL语法
- 自动生成类型安全的Go代码
- 避免运行时SQL错误

### 性能优化
- 生成高效的数据库访问代码
- 支持预处理语句
- 减少反射开销

### 开发效率
- 自动生成CRUD操作
- 减少样板代码
- 保持SQL和Go代码同步

## ⚠️ 注意事项

1. **SQL语法** - 使用PostgreSQL语法
2. **参数绑定** - 使用 `$1, $2, $3...` 格式
3. **查询命名** - 使用 `-- name: FunctionName :type` 格式
4. **文件同步** - 修改SQL后需要重新生成代码

## 🔧 故障排除

### 常见问题

1. **sqlc命令未找到**
   ```bash
   go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
   ```

2. **SQL语法错误**
   ```bash
   sqlc vet  # 验证SQL语法
   ```

3. **生成代码失败**
   - 检查sqlc.yaml配置
   - 确认SQL文件路径正确
   - 验证数据库架构文件

### 调试技巧

```bash
# 查看详细输出
sqlc generate --debug

# 验证配置
sqlc config

# 检查SQL文件
sqlc vet --strict
``` 