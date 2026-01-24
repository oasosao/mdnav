package router

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"mdnav/internal/core"
	"mdnav/internal/handler"
	"mdnav/internal/middleware"

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

	// router.StaticFile("/favicon.ico", ctx.Conf.GetString("site.favicon"))
	router.Static("/static", ctx.Conf.GetString("template.static_dir"))

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
		ctx.Log.Info("Server started at" + serverPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			ctx.Log.Fatal("listen: " + err.Error() + "\n")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx.Log.Info("Shutdown Server ...")

	ctxx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctxx); err != nil {
		ctx.Log.Error("Server Shutdown: " + err.Error())
		return
	}

	ctx.Log.Info("Server exiting")

}
