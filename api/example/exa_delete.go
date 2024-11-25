package example

import (
	"second_hand_mall/api"
	"second_hand_mall/internal/core/initialize/db"
	"second_hand_mall/internal/model/example"

	"github.com/gin-gonic/gin"
)

type Del struct {
	Id int64 `json:"id"`
}

func (d *Del) Logic(c *gin.Context) api.Result {

	var t = example.ExampleModel{Id: d.Id}

	// 自行封装db操作
	db, err := db.NewEngine(0, &t)
	if err != nil {
		return api.Result{Code: 1, Msg: err.Error()}
	}

	okCount, err := db.Delete(&t)
	if err != nil {
		return api.Result{Code: 2, Msg: err.Error()}
	}

	if okCount == 0 {
		return api.Result{Code: 3, Msg: "没有找到该数据"}
	}

	return api.Result{Code: 0}
}
