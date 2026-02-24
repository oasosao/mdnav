# 第一阶段：构建环境
FROM golang:1.25.6-alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制 go.mod 和 go.sum 文件
COPY go.mod go.sum ./

ENV GOPROXY='https://goproxy.cn'

# 安装依赖
RUN go mod tidy

# 复制项目文件
COPY . .

# 编译项目
RUN go build -o mdnav main.go

# 第二阶段：运行环境
FROM alpine:latest

# 安装必要的依赖
RUN apk --no-cache add ca-certificates

# 设置工作目录
WORKDIR /app

# 从构建阶段复制编译好的二进制文件
COPY --from=builder /app/mdnav /app/

# 复制配置文件
COPY config.yaml /app/

# 复制模板目录
COPY tpl/ /app/tpl/

# 创建内容目录（后续通过卷挂载）
RUN mkdir -p /app/contents

# 暴露端口
EXPOSE 8081

# 设置默认命令
CMD ["./mdnav"]