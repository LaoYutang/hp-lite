# server构建阶段
FROM --platform=$TARGETPLATFORM golang:1.22-alpine AS builder

WORKDIR /app

# 复制项目源码
COPY ./hp-server .

# 复制前端构建结果到Go项目
COPY ./hp-web/dist ./internal/web/static

# 下载依赖
RUN go mod download

ARG TARGETOS
ARG TARGETARCH
ARG TARGETVARIANT

# 设置 Go 编译参数
ENV CGO_ENABLED=0
ENV GOOS=$TARGETOS
ENV GOARCH=$TARGETARCH
ENV GOARM=${TARGETVARIANT#v}

# 静态编译
RUN go build -o hp-server .

# 运行阶段
FROM alpine:3.19

WORKDIR /app

COPY --from=builder /app/hp-server .
RUN apk add --no-cache tzdata ca-certificates
ENV TZ=Asia/Shanghai

# 服务入口
CMD ["./hp-server"]