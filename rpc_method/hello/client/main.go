package main

import (
	"context"
	"log"
	"second_hand_mall/rpc_method/hello"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("1.94.227.242:8012", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("连接失败:%v", err)
	}
	defer conn.Close()
	c := hello.NewHelloServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Hello(ctx, &hello.HelloRequest{Name: "测试名称"})
	if err != nil {
		log.Fatalf("请求错误:%v", err)
	}
	log.Printf("服务器返回信息:%s,%d,%s", r.HelloName, r.RespData.Code, r.RespData.Msg)
}
