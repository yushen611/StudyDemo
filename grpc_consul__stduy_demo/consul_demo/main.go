package main

import (
	"fmt"
	"net/http"
)

func main() {
	// 创建Consul客户端
	consulClient, err := NewConsulClient()
	if err != nil {
		panic(err)
	}
	// 注册服务
	serviceID := "demo-service"
	serviceName := "demo"
	serviceHost := "localhost"
	servicePort := 8080
	err = consulClient.RegisterService(serviceID, serviceName, serviceHost, servicePort)
	if err != nil {
		panic(err)
	}

	// 发现服务
	serviceAddr, err := consulClient.DiscoverService(serviceName)
	if err != nil {
		panic(err)
	}

	// 处理HTTP请求
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World!")
	})
	fmt.Printf("Starting server at %s...\n", serviceAddr)
	err = http.ListenAndServe(serviceAddr, nil)
	if err != nil {
		panic(err)
	}
}
