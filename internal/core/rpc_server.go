package core

import (
	"context"
	"fmt"
	"log"
	"net"
	"second_hand_mall/internal/global"
	"second_hand_mall/router"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

func RunRpcServer(wg *sync.WaitGroup) {
	defer wg.Done()
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", global.GVAL_CONFIG.RpcServer.Host, global.GVAL_CONFIG.RpcServer.Port))
	if err != nil {
		log.Fatalf("监听tcp失败: %v", err)
	}
	s := grpc.NewServer(grpc.KeepaliveParams(keepalive.ServerParameters{
		MaxConnectionIdle: 5 * time.Minute,  // 连接保持空闲的最大时间
		Time:              30 * time.Second, // 服务器发送保持活动（keepalive）探测的时间间隔。服务器会定期发送探测消息以检查连接是否仍然有效。
		Timeout:           10 * time.Second, // 服务器等待保持活动探测响应的超时时间。如果服务器在发送探测消息后，在 Timeout 时间内没有收到响应，它会认为连接已经断开，并关闭连接。
		//MaxConnectionAge:      30 * time.Minute, // 一个连接存在的最大时间
		//MaxConnectionAgeGrace: 5 * time.Minute, // 当超过最大时间时,给与正在进行请求的宽限时间
	}))

	// 注册健康检查服务
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(s, healthServer)

	// 设置服务状态为 SERVING
	healthServer.SetServingStatus("dbhandler.DBService", grpc_health_v1.HealthCheckResponse_SERVING)

	// 注册Rpc服务
	router.RegisterRpcMethod(s)
	//pb.RegisterGreeterServer(s, &rpc_server{})
	log.Printf("rpc 服务监听 %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("rpc服务启动失败: %v", err)
	}
}

// 定义白名单
var whitelist = map[string]bool{
	"127.0.0.1":     true,
	"192.168.1.100": true,
	// 添加更多允许的 IP 地址
}

// 自定义拦截器
func unaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	// 获取元数据
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "missing metadata")
	}

	// 获取客户端 IP 地址
	addresses := md.Get("x-forwarded-for")
	if len(addresses) == 0 {
		// 如果没有 x-forwarded-for 头，则尝试从上下文获取
		peer, ok := peer.FromContext(ctx)
		if ok {
			addresses = append(addresses, peer.Addr.String())
		}
	}

	if len(addresses) == 0 {
		return nil, status.Error(codes.Unauthenticated, "unable to determine client IP")
	}

	// 假设第一个地址是客户端的真实 IP
	clientIP := addresses[0]
	// 去除端口号
	clientIP, _, err := net.SplitHostPort(clientIP)
	if err != nil {
		clientIP = addresses[0]
	}

	if !whitelist[clientIP] {
		return nil, status.Error(codes.PermissionDenied, "your IP address is not allowed")
	}

	// 继续处理请求
	return handler(ctx, req)
}
