package conf

import (
	"github.com/spf13/viper"
)

type Config = *viper.Viper

var config Config

func InitConfig(configPath, configName string, isDebug string) error {

	config = viper.New()

	// 设置配置文件路径
	config.AddConfigPath(configPath)
	config.SetConfigName(configName)
	config.SetConfigType("yaml")

	// 读取配置文件
	if err := config.ReadInConfig(); err != nil {
		return err
	}

	// 监听配置文件变化
	if isDebug == "true" {
		config.WatchConfig()
	}

	return nil
}

func Get() Config {
	return config
}
