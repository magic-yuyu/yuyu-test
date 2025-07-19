# Go开发环境

这是一个基本的Go开发环境设置。

## 项目结构

```
yuyu-test/
├── main.go          # 主程序文件
├── go.mod           # Go模块文件
├── Dockerfile       # Docker构建文件
├── docker-compose.yml # Docker Compose配置
├── .dockerignore    # Docker忽略文件
└── README.md        # 项目说明
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

1. 确保已安装Go 1.21或更高版本

2. 运行应用：
```bash
go run main.go
```

3. 访问应用：
   - 打开浏览器访问 http://localhost:8080

## 开发命令

- `go run main.go` - 运行应用
- `go build` - 构建应用
- `go test` - 运行测试
- `go mod tidy` - 整理依赖

## Docker命令

- `docker-compose up` - 启动服务
- `docker-compose down` - 停止服务
- `docker-compose up --build` - 重新构建并启动
- `docker build -t yuyu-test .` - 构建镜像 