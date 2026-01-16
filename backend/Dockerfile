# ============================
# 阶段 1: 编译后端 (Go)
# ============================
FROM golang:1.22-alpine AS builder

# 设置工作目录
WORKDIR /app

# 1. 设置国内代理 (服务器构建时这步很重要，除非服务器在海外)
ENV GOPROXY=https://goproxy.cn,direct

# 2. 复制依赖文件并下载
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# 3. 复制源码
COPY backend/ .

# 4. 编译 (禁用 CGO，输出为 server)
RUN CGO_ENABLED=0 GOOS=linux go build -o server main.go

# ============================
# 阶段 2: 运行时 (Alpine)
# ============================
FROM alpine:latest

WORKDIR /root/

# 从编译阶段复制二进制文件
COPY --from=builder /app/server .

# 暴露端口 (假设你的 Go 代码监听 8080)
EXPOSE 8080

# 运行
CMD ["./server"]