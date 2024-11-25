package example

import (
	"second_hand_mall/api"
	"second_hand_mall/internal/core/initialize/db"
	"second_hand_mall/internal/model/example"

	"github.com/gin-gonic/gin"
)

type Get struct {
	Id int64 `form:"id"`
}

func (g *Get) Logic(c *gin.Context) api.Result {

	var t = example.ExampleModel{Id: g.Id}

	// 自行封装db操作
	db, err := db.NewEngine(0, &t)
	if err != nil {
		return api.Result{Code: 1, Msg: err.Error()}
	}

	has, err := db.Get(&t)
	if err != nil {
		return api.Result{Code: 2, Msg: err.Error()}
	}

	if !has {
		return api.Result{Code: 3, Msg: "没有找到该数据"}
	}

	return api.Result{Data: t}
}
