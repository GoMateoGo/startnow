package main

import (
	"context"
	"fmt"
	"log"
	dbhandler "second_hand_mall/rpc_method/db_demo"
	"sync"

	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
)

var (
	connMutex      sync.RWMutex
	conn           *grpc.ClientConn            // rpc连接
	client         dbhandler.DBServiceClient   // rpc服务客户端连接
	dbHealthClient grpc_health_v1.HealthClient // 健康检测客户端连接
	healthOnce     sync.Once                   // 创建健康检测实例为单例
)

func main() {
	// 关闭链接
	defer func() {
		connMutex.Lock()
		conn.Close()
		connMutex.Unlock()
	}()

	// 初始化链接
	initConnection()

	// 启动健康检查
	go healthCheckLoop(30 * time.Second) // 调整为30秒

	// 业务模拟
	for {
		time.Sleep(10 * time.Second)
		Get()
		List()
	}
}

func initConnection() {
	connMutex.Lock()
	defer connMutex.Unlock()

	var err error
	conn, err = grpc.NewClient("localhost:5001",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                20 * time.Second,
			Timeout:             10 * time.Second,
			PermitWithoutStream: true,
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
        }`),
	)
	if err != nil {
		log.Fatalf("连接失败: %v", err)
	}

	client = dbhandler.NewDBServiceClient(conn)
}

func resetConnection() {
	connMutex.Lock()
	defer connMutex.Unlock()

	log.Println("尝试重置连接...")
	conn.Close()
	initConnection() // 重新初始化连接
}

// 健康检查
func checkHealth() error {
	healthOnce.Do(func() {
		dbHealthClient = grpc_health_v1.NewHealthClient(conn)
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	resp, err := dbHealthClient.Check(ctx, &grpc_health_v1.HealthCheckRequest{
		Service: "dbhandler.DBService", // 与服务端一致
	})

	if err != nil {
		return fmt.Errorf("健康检查失败: %v", err)
	}

	if resp.Status != grpc_health_v1.HealthCheckResponse_SERVING {
		return fmt.Errorf("服务状态异常: %s", resp.Status)
	}

	log.Println("服务状态正常")
	return nil
}

// 健康检查循环
func healthCheckLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		// 1. 官方健康协议
		if err := checkHealth(); err != nil {
			log.Printf("健康检查异常: %v", err)
			resetConnection()
		}

		// 2. 检查连接状态
		connMutex.RLock()
		state := conn.GetState()
		connMutex.RUnlock()

		if !isConnectionHealthy(state) {
			log.Printf("连接状态异常: %v", state)
			resetConnection()
		}
	}
}

func isConnectionHealthy(state connectivity.State) bool {
	return state == connectivity.Ready || state == connectivity.Idle
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
