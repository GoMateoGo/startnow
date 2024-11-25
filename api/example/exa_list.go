package example

import (
	"second_hand_mall/api"
	"second_hand_mall/internal/core/initialize/db"
	"second_hand_mall/internal/model/example"

	"github.com/gin-gonic/gin"
)

type List struct {
	Page  int `form:"page" binding:"required"`
	Limit int `form:"limit" binding:"required"`
}

func (l *List) Logic(c *gin.Context) api.Result {

	var list []example.ExampleModel
	var t = &example.ExampleModel{}

	// 自行封装db操作
	db, err := db.NewEngine(0, t)
	if err != nil {
		return api.Result{Code: 1, Msg: err.Error()}
	}

	if err := db.Table(t).
		Limit(l.Limit, (l.Page-1)*l.Limit).
		Find(&list); err != nil {
		return api.Result{Code: 2, Msg: err.Error()}
	}

	return api.Result{Data: list}
}
