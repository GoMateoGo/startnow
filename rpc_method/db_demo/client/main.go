package main

import (
	"context"
	"log"
	dbhandler "second_hand_mall/rpc_method/db_demo"

	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:5001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("连接失败:%v", err)
	}
	defer conn.Close()
	c := dbhandler.NewDBServiceClient(conn)
	//Get(c)
	List(c)
}

func Get(c dbhandler.DBServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Get(ctx, &dbhandler.GetRequest{Cid: 0, Id: 3})
	if err != nil {
		log.Fatalf("请求错误:%v", err)
	}

	log.Println(r.Id, r.Name)
}

func List(c dbhandler.DBServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.List(ctx, &dbhandler.ListRequest{Cid: 2008})
	if err != nil {
		log.Fatalf("请求错误:%v", err)
	}

	log.Println(r.GetData())
}
