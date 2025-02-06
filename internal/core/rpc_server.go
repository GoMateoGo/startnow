package core

import (
	"fmt"
	"log"
	"net"
	"second_hand_mall/internal/global"
	"second_hand_mall/router"
	"sync"

	"google.golang.org/grpc"
)

func RunRpcServer(wg *sync.WaitGroup) {
	defer wg.Done()
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", global.GVAL_CONFIG.RpcServer.Host, global.GVAL_CONFIG.RpcServer.Port))
	if err != nil {
		log.Fatalf("监听tcp失败: %v", err)
	}
	s := grpc.NewServer()
	// 注册Rpc服务
	router.RegisterRpcMethod(s)
	//pb.RegisterGreeterServer(s, &rpc_server{})
	log.Printf("rpc 服务监听 %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("rpc服务启动失败: %v", err)
	}
}
