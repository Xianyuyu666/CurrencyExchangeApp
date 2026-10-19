package main

import (
	"ExchangeApp/config"
	"ExchangeApp/router"
)

func main() {
	// 初始化配置
	config.InitConfig()
	// 初始化路由，获取引擎
	r := router.SetupRouter()

	// 设置port默认值，防止为空
	port := config.AppConfig.App.Port
	if port == "" {
		port = ":8080"
	}

	// 启动服务器
	r.Run(config.AppConfig.App.Port)
}
