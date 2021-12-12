package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"kratos-demo/proto"
	"log"
	"sync"
)

type Client struct {
}

func main() {

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		// 客户端，连接服务端
		cc, err := grpc.Dial("127.0.0.1:8888", grpc.WithInsecure())
		if err != nil {
			log.Fatal("dail server error ..........", err)
		}
		client := proto.NewGreeterClient(cc)

		req := proto.HelloRequest{
			Name: "meng",
		}

		resp, err := client.SayHello(context.Background(), &req)
		if err != nil {
			log.Fatal("request server error.......", err)
		}
		fmt.Println(resp)
	}()

	wg.Wait()
}
