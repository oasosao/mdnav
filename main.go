package main

import (
	"os"

	"mdnav/internal/conf"
	"mdnav/internal/core"
	"mdnav/internal/pkg/wacher"
	"mdnav/internal/pkg/zap"
	"mdnav/internal/router"
	"mdnav/internal/service"
)

func main() {

	// 初始化Logger
	logger := zap.NewLogger()
	defer logger.Sync()

	logger.Info("应用启动")

	if err := conf.InitConfig(".", "config"); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	ctx := &core.Context{
		Log:  logger,
		Conf: conf.Get(),
	}

	if err := service.LoadAllData(ctx); err != nil {
		logger.Error("加载文档失败 " + err.Error())
		os.Exit(1)
	}

	go wacher.WatcherFile(ctx, func() {
		logger.Info("文件变化，重新加载文档")
		if err := service.LoadAllData(ctx); err != nil {
			logger.Error("加载文档失败 " + err.Error())
			os.Exit(1)
		}
	})

	router.Run(ctx)
}
