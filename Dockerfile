# --- 第一阶段：构建阶段 (Builder Stage) ---
# 使用官方的Go语言镜像作为构建环境
FROM golang:1.22-alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制 Go 模块文件
COPY go.mod go.sum ./

# 下载依赖项
RUN go mod download

# 复制所有源代码
COPY . .

# 构建Go应用。使用CGO_ENABLED=0和-ldflags来创建一个静态的、无依赖的二进制文件
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o /main .


# --- 第二阶段：运行阶段 (Final Stage) ---
# 使用一个极小的基础镜像 (scratch)，它几乎是空的，只包含我们的应用，极致安全
FROM scratch

# 从构建阶段复制编译好的二进制文件
COPY --from=builder /main /main

# （可选）如果您的应用需要处理HTTPS或访问其他Google服务，需要复制根证书
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# 暴露应用监听的端口
EXPOSE 8080

# 定义容器启动时运行的命令
ENTRYPOINT ["/main"]

# CMD ["./main"]
