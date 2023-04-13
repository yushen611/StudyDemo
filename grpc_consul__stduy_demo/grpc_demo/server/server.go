package main

import (
	"context"
	"log"
	"net"

	pb "test_grpc/test_grpc/helloworld" // 导入生成的 helloworld.pb.go 文件

	"google.golang.org/grpc"
)

const (
	port = ":50051" // 端口号
)

type server struct{} // 定义 gRPC 服务的 server

// 实现 SayHello 方法，接受客户端的请求，返回 HelloResponse 响应
func (s *server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{Message: "Hello " + req.Name}, nil
}

func main() {
	lis, err := net.Listen("tcp", port) // 监听 TCP 端口
	if err != nil {
		log.Fatalf("failed to listen: %v", err) // 如果监听失败，则退出程序并打印错误信息
	}
	s := grpc.NewServer()                  // 创建一个 gRPC 服务实例
	pb.RegisterGreeterServer(s, &server{}) // 注册服务，把实现服务的 server 结构体绑定到 gRPC 服务实例中
	log.Printf("start gRPC server on port %s\n", port)
	if err := s.Serve(lis); err != nil { // 开始监听
		log.Fatalf("failed to serve: %v", err) // 开始监听，如果监听失败，则退出程序并打印错误信息
	}
}
