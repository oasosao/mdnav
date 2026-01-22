package main

import (
	"log/slog"
	"os"

	"mdnav/internal/conf"
	"mdnav/internal/core"
	"mdnav/internal/pkg/logger"
	"mdnav/internal/router"
	"mdnav/internal/store"
)

func main() {

	// 初始化Logger
	logger := logger.NewLogger()
	defer logger.Sync()

	ctx := &core.Context{
		Logger: logger,
	}

	conf.InitConfig(ctx, ".", "config")

	if err := store.LoadAllDocuments(ctx); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	go store.WatcherFile(ctx, func() {
		store.LoadAllDocuments(ctx)
	})

	router.Run(ctx)
}
