package example

import (
	"log"
	"second_hand_mall/api"
	"second_hand_mall/internal/global"
	"second_hand_mall/middleware"

	"github.com/gin-gonic/gin"
)

type Login struct {
	UserId   int64  `json:"user_id" binding:"required"`   // 用户ID 创建时返回,在赋值
	UserName string `json:"user_name" binding:"required"` // 用户名字
	Nick     string `json:"nick" binding:"required"`      // 用户昵称
	Phone    string `json:"phone" binding:"required"`     // 手机号
	Email    string `json:"email" binding:"required"`     // 邮箱
}

func (l *Login) Logic(c *gin.Context) api.Result {
	var jwt = middleware.SaleLogin{
		JwtBase: middleware.JwtBase{
			UserId:   l.UserId,
			UserName: l.UserName,
			Nick:     l.Nick,
			Phone:    l.Phone,
			Email:    l.Email,
		},
	}
	key := global.GVAL_CONFIG.JWT.SigningKey
	log.Println("token 加密key:", key)
	t := jwt.CreateToken()
	token, err := t.SignedString([]byte(key))
	if err != nil {
		return api.Result{Code: 1, Msg: "生成Token失败:" + err.Error()}
	}
	return api.Result{Data: token}
}
