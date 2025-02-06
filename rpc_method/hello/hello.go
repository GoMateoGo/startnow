package hello

import (
	context "context"
	"log"
)

type RpcHello struct {
	UnimplementedHelloServiceServer
}

// 路由组下的具体方法
func (s *RpcHello) Hello(ctx context.Context, in *HelloRequest) (*HelloResponse, error) {
	log.Printf("接收: %v\n", in.GetName())
	return &HelloResponse{
		RespData:  &RetData{Code: 1, Msg: "aaa"},
		HelloName: "名字",
	}, nil
}
