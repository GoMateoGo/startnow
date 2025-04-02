package main

import (
	"context"
	"fmt"
	"log"
	student "second_hand_mall/rpc_method/ca/proto"
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
	client         student.DBServiceClient     // rpc服务客户端连接
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
	go healthCheckLoop(15 * time.Second) // 调整为30秒

	// 业务模拟
	for {
		time.Sleep(10 * time.Second)
		ExampleUsage()
	}
}

func initConnection() {
	connMutex.Lock()
	defer connMutex.Unlock()

	var err error
	conn, err = grpc.NewClient("localhost:8054",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                20 * time.Second,
			Timeout:             10 * time.Second,
			PermitWithoutStream: true,
		}),
		//grpc.WithDefaultServiceConfig(`{
		//    "methodConfig": [{
		//        "name": [{"service": "student.DBService"}],
		//        "retryPolicy": {
		//            "MaxAttempts": 3,
		//            "InitialBackoff": "0.1s",
		//            "MaxBackoff": "1s",
		//            "BackoffMultiplier": 2.0,
		//            "RetryableStatusCodes": ["UNAVAILABLE"]
		//        }
		//    }]
		//}`),
	)
	if err != nil {
		log.Fatalf("连接失败: %v", err)
	}

	client = student.NewDBServiceClient(conn)
}

func resetConnection() {
	connMutex.Lock()
	defer connMutex.Unlock()

	log.Println("尝试重置连接...")
	conn.Close()
	initConnection() // 重新初始化连接
}

// 定义需要检查的服务列表
var servicesToCheck = []string{
	"student.DBService",
	"user.UserService",
	"order.OrderService",
}

// 健康检查单个服务
func checkSingleService(ctx context.Context, service string) error {
	resp, err := dbHealthClient.Check(ctx, &grpc_health_v1.HealthCheckRequest{
		Service: service,
	})

	if err != nil {
		return fmt.Errorf("服务 %s 健康检查失败: %v", service, err)
	}

	if resp.Status != grpc_health_v1.HealthCheckResponse_SERVING {
		return fmt.Errorf("服务 %s 状态异常: %s", service, resp.Status)
	}

	log.Printf("服务 %s 状态正常", service)
	return nil
}

// 健康检查
func checkHealth() error {
	healthOnce.Do(func() {
		dbHealthClient = grpc_health_v1.NewHealthClient(conn)
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var errors []error
	for _, service := range servicesToCheck {
		if err := checkSingleService(ctx, service); err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("健康检查发现异常: %v", errors)
	}

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

// 添加调用 Add 方法的函数
func Add(data []*student.AddListData) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &student.AddRequest{
		Data: data,
	}

	resp, err := client.Add(ctx, req)
	if err != nil {
		return fmt.Errorf("调用Add失败: %v", err)
	}

	if resp.Code != 0 { // 假设 0 表示成功
		return fmt.Errorf("添加失败，错误码: %d", resp.Code)
	}

	return nil
}

// 使用示例
func ExampleUsage() {
	data := []*student.AddListData{
		{
			StudentId:  10001,
			CId:        2008,
			CName:      "上海橙舟",
			PId:        301,
			PName:      "法律资格考试",
			UserName:   "周作人",
			ClassId:    2353,
			ClassName:  "2025年法考【VIP旗舰班】",
			Phone:      "13800138000",
			BillNumber: "BILL001",
			Amount:     100.0,
			Type:       "PUBLIC",
			PayCid:     16,
			UserId:     1,
		},
		{
			StudentId:  10002,
			CId:        2008,
			CName:      "上海橙舟",
			PId:        301,
			PName:      "法律资格考试",
			UserName:   "周树人",
			ClassId:    22501,
			ClassName:  "2025年法考【重点协议班】",
			Phone:      "13800138001",
			BillNumber: "BILL001",
			Amount:     100.0,
			Type:       "PUBLIC",
			PayCid:     16,
			UserId:     2,
		},
	}

	if err := Add(data); err != nil {
		log.Printf("添加失败: %v", err)
		return
	}
	log.Println("添加成功")
}
