package config

import (
	"log"

	"github.com/spf13/viper"
)

// Config 配置模型
type Config struct {
	App struct {
		Name string
		Port string
	}
	Database struct {
		Dsn          string
		MaxIdleConns int
		MaxOpenConns int
	}
	Jwt struct {
		Secret      string
		ExpireHours int
	}
}

// AppConfig 全局配置单例
var AppConfig *Config

// InitConfig 配置初始化函数:利用viper获取配置映射到全局配置单例中
func InitConfig() {
	// 设置配置文件路径和名称
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath("./config")

	// 配置读取到viper内部
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("error reading config file: %v", err)
	}

	// 从viper内部反序列化映射到全局配置单例中
	AppConfig = &Config{}
	if err := viper.Unmarshal(AppConfig); err != nil {
		log.Fatalf("Unable to decode into struct: %v", err)
	}

	// 初始化数据库
	initDB()
}
