# 构建阶段
FROM --platform=$BUILDPLATFORM golang:alpine AS builder

WORKDIR /app

# 复制项目源码
COPY . .

# 下载依赖
RUN go mod tidy

ARG TARGETOS
ARG TARGETARCH
ARG TARGETVARIANT

# 设置 Go 编译参数
ENV CGO_ENABLED=0
ENV GOOS=$TARGETOS
ENV GOARCH=$TARGETARCH

# 针对 ARM 架构设置 GOARM
RUN if [ "$TARGETARCH" = "arm" ]; then \
  case "${TARGETVARIANT}" in \
  "v7")  export GOARM=7 ;; \
  "v6")  export GOARM=6 ;; \
  *)     export GOARM=7 ;; \
  esac; \
  fi

# 静态编译
RUN go build -o hp-server .

# 运行阶段
FROM alpine:latest

WORKDIR /app

# 复制二进制文件和时区数据
COPY --from=builder /app/hp-server .
RUN apk add --no-cache tzdata

# 服务入口
CMD ["./hp-server"]