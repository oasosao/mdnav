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

	"github.com/gin-gonic/gin"
)

func Run(ctx *core.Context) {

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(middleware.ZapLoggerWithConfig(ctx.Logger))
	router.Use(middleware.Options())
	router.Use(middleware.IpRateLimiter())
	router.Use(gin.Recovery())

	h := &handler.Handler{
		Ctx: ctx,
	}

	router.StaticFile("/favicon.ico", conf.Config().GetString("site.favicon"))
	router.Static("/static", conf.Config().GetString("template.static_dir"))

	router.GET("/", h.Index)
	router.GET("/:slug", h.Category)
	router.GET("/tag/:slug", h.Tag)
	router.GET("/article/*slug", h.Article)

	serverPort := conf.Config().GetString("server.port")
	srv := &http.Server{
		Addr:           serverPort,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		ctx.Logger.Info("Server started at" + serverPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			ctx.Logger.Fatal("listen: " + err.Error() + "\n")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx.Logger.Info("Shutdown Server ...")

	ctxx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctxx); err != nil {
		ctx.Logger.Info("Server Shutdown: " + err.Error())
	}

	ctx.Logger.Info("Server exiting")

}
