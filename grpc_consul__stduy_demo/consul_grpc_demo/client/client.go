package main

import (
	"context"
	"log"
	"time"

	cs "client/consul"
	pb "client/proto"

	"google.golang.org/grpc"
)

func main() {
	//使用consul获取指定服务的ip与端口
	//consul-1:连接至consul服务中心
	consulClient, err := cs.NewConsulClient()
	if err != nil {
		panic(err)
	}
	//consul-2:发现指定服务名称的服务地址
	serviceName := "demo"
	serviceAddr, err := consulClient.DiscoverService(serviceName)
	if err != nil {
		panic(err)
	}
	// grpc-1.连接到指定ip,端口的grpc服务器

	// serviceAddr := "localhost:50051"
	conn, err := grpc.Dial(serviceAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// grpc-2.初始化HelloService服务客户端
	client := pb.NewHelloServiceClient(conn)

	//grpc-3.初始化上下文，设置请求超时时间为1秒
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	//grpc-4.调用远程的grpc方法
	response, err := client.SayHello(ctx, &pb.HelloRequest{Name: "Yu Gambler"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", response.GetMessage())
}
