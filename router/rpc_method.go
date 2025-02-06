package router

import (
	dbhandler "second_hand_mall/rpc_method/db_demo"
	"second_hand_mall/rpc_method/hello"

	"google.golang.org/grpc"
)

// 注册rpc服务所有方法
func RegisterRpcMethod(s *grpc.Server) {
	hello.RegisterHelloServiceServer(s, &hello.RpcHello{})
	dbhandler.RegisterDBServiceServer(s, &dbhandler.Db{})
}
