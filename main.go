package main

import (
	"ExchangeApp/config"
	"fmt"
)

func main() {
	config.InitConfig()

	fmt.Println(config.AppConfig.App.Port)
}
