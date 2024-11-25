package example

import (
	"second_hand_mall/api"
	"second_hand_mall/internal/core/initialize/db"
	"second_hand_mall/internal/model/example"

	"github.com/gin-gonic/gin"
)

type Update struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

func (u *Update) Logic(c *gin.Context) api.Result {

	var t = example.ExampleModel{Name: u.Name}

	// 自行封装db操作
	db, err := db.NewEngine(0, &t)
	if err != nil {
		return api.Result{Code: 1, Msg: err.Error()}
	}

	okCount, err := db.Where("id=?", u.Id).Update(&t)
	if err != nil {
		return api.Result{Code: 2, Msg: err.Error()}
	}

	if okCount == 0 {
		return api.Result{Code: 3, Msg: "更新了0条数据"}
	}

	return api.Result{Code: 0}
}
