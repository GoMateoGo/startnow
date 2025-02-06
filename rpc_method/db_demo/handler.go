package dbhandler

import (
	"context"
	"fmt"
	"second_hand_mall/internal/core/initialize/db"
	"second_hand_mall/internal/model/example"
)

type Db struct {
	UnimplementedDBServiceServer
}

func (d *Db) Get(ctx context.Context, in *GetRequest) (*GetResponse, error) {
	var t = example.ExampleModel{Id: int64(in.Id)}

	// 自行封装db操作
	db, err := db.NewEngine(in.Cid, &t)
	if err != nil {
		return nil, err
	}
	session := db.NewSession()
	defer session.Close()
	has, err := session.Get(&t)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, fmt.Errorf("没有找到cid:%d, id为:%d的数据", in.Cid, in.Id)
	}
	return &GetResponse{Id: int32(t.Id), Name: t.Name}, nil
}

func (d *Db) List(ctx context.Context, in *ListRequest) (*ListResponse, error) {
	var t example.ExampleModel
	var resp []example.ExampleModel
	// 自行封装db操作
	db, err := db.NewEngine(in.Cid, &t)
	if err != nil {
		return nil, err
	}
	session := db.NewSession()
	defer session.Close()
	err = session.Table(t.TableName()).Find(&resp)
	if err != nil {
		return nil, err
	}

	var list []*ListData
	for _, v := range resp {
		list = append(list, &ListData{Id: v.Id, Name: v.Name})
	}

	return &ListResponse{Data: list}, nil
}
