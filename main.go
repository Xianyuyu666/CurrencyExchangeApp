package main

import (
	"ExchangeApp/config"
	"ExchangeApp/router"
)

func main() {
	//初始化配置
	config.InitConfig()
	//获取路由
	r := router.SetupRouter()

	port := config.AppConfig.App.Port
	if port == "" {
		port = ":8080"
	}

	//启动服务器
	r.Run(config.AppConfig.App.Port)
}
