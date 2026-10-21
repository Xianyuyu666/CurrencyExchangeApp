package router

import (
	"ExchangeApp/controllers"
	"ExchangeApp/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRouter 初始化路由函数:负责获取Gin引擎和路由注册
func SetupRouter() *gin.Engine {
	// 获取引擎
	r := gin.Default()

	// 注册路由
	auth := r.Group("/api/auth")
	{
		auth.POST("/login", controllers.Login)
		auth.POST("/register", controllers.Register)
	}
	user := r.Group("/api/user", middleware.JWTAuth())
	{
		user.GET("/info", controllers.Info)
	}
	exchangeRate := r.Group("/api/exchangerate")
	{
		exchangeRate.POST("/add", controllers.CreateExchangeRate)
		exchangeRate.GET("/query", controllers.Query)
	}

	return r
}
