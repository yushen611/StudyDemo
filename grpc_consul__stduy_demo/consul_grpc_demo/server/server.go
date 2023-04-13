package main

import (
	"context"
	"log"
	"net"

	cs "server/consul"
	pb "server/proto/server"

	"google.golang.org/grpc"
)

// 定义 实现 gRPC 服务的 server结构体

type serverImp struct {
	pb.UnimplementedHelloServiceServer
}

// 实现 SayHello 方法，接受客户端的请求，返回 HelloResponse 响应
func (s *serverImp) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{Message: "Hello " + req.Name}, nil
}

func main() {

	//Consul-1 创建Consul客户端
	consulClient, err := cs.NewConsulClient()
	if err != nil {
		log.Fatalf("failed to Connect consul: %v", err)
		panic(err)
	}
	//Consul-2 注册服务
	serviceID := "demo-service"
	serviceName := "demo"
	serviceHost := "127.0.0.1"
	servicePort := 50051
	err = consulClient.RegisterService(serviceID, serviceName, serviceHost, servicePort)
	if err != nil {
		log.Fatalf("failed to RegisterService: %v", err)
		panic(err)
	}

	//grpc-1.定义服务启动端口
	port := ":50051"
	lis, err := net.Listen("tcp", port) // 监听 TCP 端口
	if err != nil {
		log.Fatalf("failed to listen: %v", err) // 如果监听失败，则退出程序并打印错误信息
	}

	//grpc-2. 创建一个 gRPC 服务实例
	s := grpc.NewServer()

	//grpc-3.把实现的服务结构体 与 gRPC 服务实例绑定
	pb.RegisterHelloServiceServer(s, &serverImp{}) // 注册服务，把实现服务的 server 结构体绑定到 gRPC 服务实例中

	log.Printf("start gRPC server on port %s\n", port)

	//grpc-4:启动gRPC服务
	err = s.Serve(lis)

	if err != nil { // 开始监听
		log.Fatalf("failed to serve: %v", err) // 开始监听，如果监听失败，则退出程序并打印错误信息
	}
}
