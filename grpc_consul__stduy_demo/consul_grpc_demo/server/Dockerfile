FROM golang:1.20

WORKDIR /app/server

# 设置环境变量
ENV BUILD_ENV local

# 暴露端口
EXPOSE 50051

# 复制应用程序到工作目录
ADD . .

RUN ls

# 编译应用程序
RUN go build -o server server.go

# RUN go build -ldflags "-w -X main.version=${VERSION}" -o /nats

# 运行应用程序
CMD ["./server"]