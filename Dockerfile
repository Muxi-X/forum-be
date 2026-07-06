# 第一阶段：构建可执行文件
FROM golang:1.26.4 AS builder

# 设置工作目录
WORKDIR /app

# 将源代码添加到容器中
COPY . /app/

# 设置 Go 环境变量
ARG service_name
RUN go env -w GOPROXY="https://goproxy.cn,direct"

# 切换到具体服务的目录
WORKDIR /app/microservice/$service_name

# 执行 make 编译项目
RUN go mod tidy && go build -o main .

# 第二阶段：使用 debian 镜像作为基础
FROM debian:bookworm-slim

ARG service_name

# 安装必需的库
RUN apt-get update && apt-get install -y \
    libc6 \
    tzdata \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# 设置时区为 Asia/Shanghai
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

# 设置工作目录
WORKDIR /app

# 从构建阶段复制可执行文件到最终镜像
COPY --from=builder /app/microservice/$service_name/main /app/main

# 确保执行权限
RUN chmod +x /app/main

# 设置默认启动命令
CMD ["./main"]
