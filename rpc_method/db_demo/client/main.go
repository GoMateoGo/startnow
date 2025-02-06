package main

import (
	"context"
	"log"
	dbhandler "second_hand_mall/rpc_method/db_demo"
	"sync"

	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

var (
	conn   *grpc.ClientConn
	client dbhandler.DBServiceClient
)

func main() {
	defer conn.Close()
	var err error
	conn, err = grpc.NewClient("localhost:5001",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                20 * time.Second, // 每隔 30 秒发送心跳
			Timeout:             10 * time.Second, // 等待确认的超时时间
			PermitWithoutStream: true,             // 无活跃流时仍发送心跳
		}),
		grpc.WithDefaultServiceConfig(`{
        "methodConfig": [{
            "name": [{"service": "dbhandler.DBService"}],
            "retryPolicy": {
                "MaxAttempts": 3,
                "InitialBackoff": "0.1s",
                "MaxBackoff": "1s",
                "BackoffMultiplier": 2.0,
                "RetryableStatusCodes": ["UNAVAILABLE"]
            }
        }]
    }`))
	if err != nil {
		log.Fatalf("连接失败:%v", err)
	}

	client = dbhandler.NewDBServiceClient(conn)

	var wg sync.WaitGroup

	// 初始化 ticker，每1秒触发一次
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	// 启动健康检查
	wg.Add(1)
	go checkHealth(&wg, ticker)
	go func() {
		for {
			time.Sleep(30 * time.Second)
			Get()
			List()
		}
	}()
	wg.Wait()
}
func checkHealth(wg *sync.WaitGroup, ticker *time.Ticker) {
	defer wg.Done()

	for {
		select {
		case <-ticker.C:
			state := conn.GetState()
			if state != connectivity.Ready && state != connectivity.Idle {
				log.Println("连接不健康，尝试重置...")
				conn.ResetConnectBackoff()
			} else {
				log.Println("rpc client 连接健康", state)
			}
		case <-time.After(5 * time.Second):
			// 设置超时，防止 ticker 阻塞
			log.Println("检测超时")
			return
		}
	}
}

func Get() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := client.Get(ctx, &dbhandler.GetRequest{Cid: 0, Id: 3})
	if err != nil {
		log.Printf("请求错误:%v", err)
		return
	}

	log.Println(r.Id, r.Name)
}

func List() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := client.List(ctx, &dbhandler.ListRequest{Cid: 2008})
	if err != nil {
		log.Printf("请求错误:%v", err)
		return
	}

	log.Println(r.GetData())
}
