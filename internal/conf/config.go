package conf

import (
	"mdnav/internal/core"

	"github.com/spf13/viper"
)

var config *viper.Viper

func InitConfig(ctx *core.Context, configPath, configName string) {

	config = viper.New()

	// 设置配置文件路径
	config.AddConfigPath(configPath)
	config.SetConfigName(configName)
	config.SetConfigType("yaml")

	// 读取配置文件
	if err := config.ReadInConfig(); err != nil {
		ctx.Logger.Fatal("Error reading config file" + err.Error())
	}

	// 监听配置文件变化
	config.WatchConfig()
}

func Config() *viper.Viper {
	return config
}
