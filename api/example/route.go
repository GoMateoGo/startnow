package example

import (
	"second_hand_mall/api"

	"github.com/gin-gonic/gin"
)

type API struct{}

func (b *API) RegisterRoutes(rgPublic *gin.RouterGroup, rgPrivate *gin.RouterGroup) {
	auth := rgPrivate.Group("db")
	{
		auth.POST("insert", func(c *gin.Context) { api.JSON(c, &Insert{}) })
		auth.POST("del", func(c *gin.Context) { api.JSON(c, &Del{}) })
		auth.GET("get", func(c *gin.Context) { api.FORM(c, &Get{}) })
		auth.POST("update", func(c *gin.Context) { api.JSON(c, &Update{}) })
		auth.GET("list", func(c *gin.Context) { api.FORM(c, &List{}) })
	}
	public := rgPublic.Group("user")
	{
		public.POST("login", func(c *gin.Context) { api.JSON(c, &Login{}) })
	}
}
