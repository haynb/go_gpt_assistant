# Dockerfile
# 从官方Go镜像开始构建
FROM golang:1.20 as builder

# 在容器内部设置工作目录
WORKDIR /src

# 将gomod和gosum文件复制到工作目录
COPY go.mod go.sum ./

# 下载所有依赖项
RUN go mod tidy

# 将源代码复制到工作目录
COPY . .

# 编译应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

# 使用 scratch 作为基础镜像创建一个新的阶段
FROM scratch

# 从 builder 阶段复制编译好的二进制文件
COPY --from=builder /src/app .

# 暴露容器端口
EXPOSE 9964

# 容器启动时执行的命令
CMD ["./app"]