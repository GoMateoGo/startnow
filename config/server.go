package config

type HttpServer struct {
	Host        string
	Port        string
	Name        string
	RoutePrefix string
}

type RpcServer struct {
	Host string
	Port string
}
