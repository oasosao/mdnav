package router

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"mdnav/internal/conf"
	"mdnav/internal/core"
	"mdnav/internal/handler"
	"mdnav/internal/middleware"
	"mdnav/internal/pkg/zap"
	"mdnav/internal/service"

	"github.com/gin-gonic/gin"
)

func Run(ctx *core.Context) {

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.RequestError(ctx))
	router.Use(middleware.Logger(ctx))
	router.Use(middleware.Options(ctx))

	h := &handler.Handler{
		Ctx:    ctx,
		TplDir: ctx.Conf.GetString("template.dir"),
	}

	router.Static("/static", ctx.Conf.GetString("template.static_dir"))

	authorized := router.Group("/system", gin.BasicAuth(gin.Accounts{
		"admin-manger": "admin-oaeoe-password",
	})).Use(middleware.IpRateLimiter(ctx))

	authorized.GET("/update", func(c *gin.Context) {
		if err := conf.InitConfig(".", "config", "false"); err != nil {
			h.Ctx.Log.Error("加载配置出错", zap.Error(err))
		}
		if err := service.LoadAllData(h.Ctx); err != nil {
			h.Ctx.Log.Error("加载数据出错", zap.Error(err))
		}
		c.AbortWithStatus(200)
	})

	r := router.Group("").Use(middleware.IpRateLimiter(ctx))
	r.GET("/", h.Index)

	r.GET("/:slug", h.Category)
	r.GET("/tag/:tagName", h.Tag)
	r.GET("/article/*slug", h.Article)

	serverPort := ctx.Conf.GetString("server.port")
	srv := &http.Server{
		Addr:           serverPort,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		ctx.Log.Info("服务启动", zap.String("host", serverPort))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			ctx.Log.Fatal("服务启动失败", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx.Log.Info("停止服务中...")

	ctxx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctxx); err != nil {
		ctx.Log.Error("停止服务出错", zap.Error(err))
		return
	}

	ctx.Log.Info("服务退出")
}
