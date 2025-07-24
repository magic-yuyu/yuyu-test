# 开发脚本

这个文件夹包含了IDaaS平台的开发和部署脚本。

## 脚本说明

### Windows 脚本

- **start-dev.bat** - 启动开发环境（需要已有数据库）
- **start-complete.bat** - 完整启动（包含数据库）
- **test-env.bat** - 环境变量测试
- **build.bat** - 构建可执行文件
- **migrate.bat** - 数据库迁移脚本
- **dev-tools.bat** - 开发工具菜单
- **sqlc.bat** - SQLC代码生成工具

### Linux/Mac 脚本

- **start-dev.sh** - 启动开发环境
- **build.sh** - 构建可执行文件

## 使用方法

### 启动开发环境

**完整启动（推荐新手）：**
```cmd
scripts\start-complete.bat
```

**仅启动应用（需要已有数据库）：**
```cmd
scripts\start-dev.bat
```

**Linux/Mac:**
```bash
chmod +x scripts/start-dev.sh
./scripts/start-dev.sh
```

### 构建应用

**Windows:**
```cmd
scripts\build.bat
```

**Linux/Mac:**
```bash
chmod +x scripts/build.sh
./scripts/build.sh
```

### 开发工具菜单

**Windows:**
```cmd
scripts\dev-tools.bat
```

### 数据库迁移

**Windows:**
```cmd
scripts\migrate.bat
```

## 环境变量

这些脚本会自动设置以下环境变量：

- `DATABASE_URL` - PostgreSQL连接字符串
- `JWT_SECRET` - JWT签名密钥
- `PORT` - 服务端口（默认8080）
- `GO_ENV` - 环境标识（development）

## 注意事项

1. 运行脚本前请确保PostgreSQL数据库已启动
2. 首次运行需要执行数据库迁移
3. 脚本中的数据库连接信息请根据实际情况调整 