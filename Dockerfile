# 构建阶段
FROM golang:1.21-alpine AS builder

# 设置工作目录
WORKDIR /app

# 安装构建依赖
RUN apk add --no-cache git make

# 复制go mod文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN make build

# 运行阶段
FROM alpine:latest

# 安装运行时依赖
RUN apk add --no-cache ca-certificates tzdata

# 设置时区
ENV TZ=Asia/Shanghai

# 创建非root用户
RUN addgroup -g 1000 yggdrasil && \
    adduser -D -u 1000 -G yggdrasil yggdrasil

# 设置工作目录
WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/yggdrasil-api-server /app/

# 创建必要的目录
RUN mkdir -p /app/conf /app/keys /app/data /app/logs && \
    chown -R yggdrasil:yggdrasil /app

# 切换到非root用户
USER yggdrasil

# 暴露端口
EXPOSE 8080

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/ || exit 1

# 启动命令
CMD ["./yggdrasil-api-server", "-config", "/app/conf/config.yml"]