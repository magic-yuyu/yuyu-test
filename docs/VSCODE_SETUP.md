# VS Code 开发环境配置

## 📋 配置文件说明

本项目包含了完整的VS Code开发环境配置，位于 `.vscode/` 目录下：

### 🔧 配置文件

- **`launch.json`** - 调试配置
- **`settings.json`** - 编辑器设置
- **`tasks.json`** - 任务定义
- **`extensions.json`** - 扩展推荐

## 🚀 调试配置

### 可用的调试配置

1. **🚀 启动IDaaS (开发模式)**
   - 标准开发模式启动
   - 自动设置环境变量
   - 集成终端输出

2. **🐛 调试IDaaS (带断点)**
   - 支持断点调试
   - 变量查看
   - 调用栈跟踪

3. **🧪 测试模式**
   - 运行测试用例
   - 测试环境变量
   - 详细测试输出

4. **🔧 构建并运行**
   - 优化构建
   - 生产环境配置

5. **📊 性能分析模式**
   - 调试符号保留
   - 性能分析支持

6. **🔍 远程调试**
   - 远程连接调试
   - 端口2345

## 🎯 使用方法

### 启动调试

1. **按F5** 或点击调试按钮
2. **选择调试配置**：
   - 开发模式：`🚀 启动IDaaS (开发模式)`
   - 调试模式：`🐛 调试IDaaS (带断点)`

### 设置断点

1. 在代码行号左侧点击设置断点
2. 启动调试模式
3. 程序会在断点处暂停
4. 查看变量、调用栈等信息

### 环境变量

所有调试配置都预设置了必要的环境变量：

```json
{
    "DATABASE_URL": "postgresql://postgres:password@localhost:5432/idaas_dev",
    "JWT_SECRET": "your-development-secret-key",
    "PORT": "8080",
    "GO_ENV": "development"
}
```

## 📝 任务配置

### 可用的任务

1. **🚀 启动开发服务器** (默认构建任务)
2. **🐳 启动完整环境**
3. **🔨 构建应用**
4. **🧪 运行测试**
5. **📊 运行测试覆盖率**
6. **🔍 代码检查**
7. **📝 格式化代码**
8. **📦 整理依赖**
9. **🗄️ 生成数据库代码**
10. **🔄 数据库迁移**
11. **🧹 清理构建文件**
12. **📋 显示项目信息**

### 使用方法

1. **Ctrl+Shift+P** 打开命令面板
2. 输入 "Tasks: Run Task"
3. 选择要执行的任务

### 快捷键

- **Ctrl+Shift+B** - 运行默认构建任务
- **Ctrl+Shift+P** - 打开命令面板

## ⚙️ 编辑器设置

### 自动格式化

- 保存时自动格式化代码
- 自动整理导入语句
- 使用 `goimports` 格式化工具

### 代码检查

- 保存时自动运行 `go vet`
- 保存时自动运行 `golangci-lint`
- 实时显示错误和警告

### 文件排除

自动排除以下文件：
- `vendor/` 目录
- `node_modules/` 目录
- `.git/` 目录
- 可执行文件 (`.exe`)
- 测试文件 (`.test`)

## 🔌 推荐扩展

### 必需扩展

- **Go** - Go语言支持
- **JSON** - JSON文件支持
- **YAML** - YAML文件支持
- **PowerShell** - Windows脚本支持

### 推荐扩展

- **Docker** - Docker容器支持
- **GitLens** - Git增强功能
- **Code Spell Checker** - 拼写检查
- **MarkdownLint** - Markdown检查
- **VSCode Icons** - 文件图标

## 🎨 主题设置

- **颜色主题**: Default Dark+
- **图标主题**: VSCode Icons
- **字体大小**: 12px
- **行高**: 18px

## 🔧 自定义配置

### 修改环境变量

编辑 `.vscode/launch.json` 中的 `env` 部分：

```json
"env": {
    "DATABASE_URL": "your-database-url",
    "JWT_SECRET": "your-jwt-secret",
    "PORT": "8080",
    "GO_ENV": "development"
}
```

### 添加新的调试配置

在 `launch.json` 的 `configurations` 数组中添加：

```json
{
    "name": "自定义配置",
    "type": "go",
    "request": "launch",
    "mode": "auto",
    "program": "${workspaceFolder}/cmd/server/main.go",
    "env": {
        "CUSTOM_VAR": "value"
    }
}
```

### 添加新的任务

在 `tasks.json` 的 `tasks` 数组中添加：

```json
{
    "label": "自定义任务",
    "type": "shell",
    "command": "your-command",
    "group": "build",
    "presentation": {
        "echo": true,
        "reveal": "always",
        "focus": false,
        "panel": "shared"
    }
}
```

## 🚨 故障排除

### 调试不工作

1. 确保Go扩展已安装
2. 检查Go工具是否正确安装
3. 验证环境变量设置

### 任务执行失败

1. 检查脚本文件是否存在
2. 验证文件路径是否正确
3. 确认依赖工具已安装

### 扩展不工作

1. 重新加载VS Code窗口
2. 检查扩展是否启用
3. 查看扩展输出日志

## 📚 相关文档

- [Go扩展文档](https://github.com/golang/vscode-go)
- [VS Code调试文档](https://code.visualstudio.com/docs/editor/debugging)
- [VS Code任务文档](https://code.visualstudio.com/docs/editor/tasks) 