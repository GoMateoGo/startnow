package core

import (
	"fmt"
	"log"
	"net"
	"second_hand_mall/internal/global"
	"second_hand_mall/router"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
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
	// 注册Rpc服务
	router.RegisterRpcMethod(s)
	//pb.RegisterGreeterServer(s, &rpc_server{})
	log.Printf("rpc 服务监听 %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("rpc服务启动失败: %v", err)
	}
}
