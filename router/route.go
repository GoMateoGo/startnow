package router

import (
	"net/http"
	"second_hand_mall/api/example"
	"second_hand_mall/internal/global"
	"second_hand_mall/middleware"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// IRouteRegister 定义路由注册的接口
type IRouteRegister interface {
	RegisterRoutes(rgPublic *gin.RouterGroup, rgPrivate *gin.RouterGroup)
}

// 初始化总路由
func InitRouter() *gin.Engine {
	Router := gin.New()
	Router.Use(gin.Recovery(), gin.Logger())
	//if gin.Mode() == gin.DebugMode {
	//	Router.Use(gin.Logger())
	//}

	Router.Use(cors.New(cors.Config{
		AllowOriginFunc:  func(origin string) bool { return true },
		AllowMethods:     []string{"OPTIONS", "GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	CreateRouteGroup(Router)
	return Router
}

// 创建路由组权限
func CreateRouteGroup(r *gin.Engine) {

	// 初始化路由模块
	routeRegisters := initBaseRoutes()

	// 路由前缀
	PublicGroup := r.Group(global.GVAL_CONFIG.Server.RoutePrefix)
	PrivateGroup := r.Group(global.GVAL_CONFIG.Server.RoutePrefix)

	PrivateGroup.Use(middleware.JwtSaleLogin()) // 鉴权中间件

	// 注册路由
	for _, routeRegister := range routeRegisters {
		routeRegister.RegisterRoutes(PublicGroup, PrivateGroup)
	}

	// 健康监测
	PublicGroup.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, "ok")
	})
}

// 初始化路由模块
func initBaseRoutes() []IRouteRegister {
	var r []IRouteRegister
	// 初始化路由注册器
	r = append(r, &example.API{}) // 示例

	return r
}
