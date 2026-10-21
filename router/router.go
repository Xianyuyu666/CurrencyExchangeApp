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
	// 登录注册相关
	auth := r.Group("/api/auth")
	{
		auth.POST("/login", controllers.Login)
		auth.POST("/register", controllers.Register)
	}
	// 用户相关
	user := r.Group("/api/user", middleware.JWTAuth())
	{
		user.GET("/info", controllers.Info)
	}
	// 汇率相关
	exchangeRate := r.Group("/api/exchangerate")
	{
		exchangeRate.POST("/add", controllers.CreateExchangeRate)
		exchangeRate.GET("/query", controllers.Query)
	}
	// 文章相关
	article := r.Group("api/article")
	{
		article.POST("/create", middleware.JWTAuth(), controllers.CreateArticle)
		article.GET("", controllers.ListArticles)
	}

	return r
}
