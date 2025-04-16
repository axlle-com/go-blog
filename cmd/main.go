package main

import (
	"context"
	"encoding/gob"
	"errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/axlle-com/blog/app"
	"github.com/axlle-com/blog/app/config"
	"github.com/axlle-com/blog/app/db"
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/app/models/contracts"
	"github.com/axlle-com/blog/app/routes"
	"github.com/axlle-com/blog/app/web"
	user "github.com/axlle-com/blog/pkg/user/models"
)

func main() {
	appCtx, stop := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
	defer stop()

	container := app.New(appCtx)

	cfg := config.Config()
	router := Init(cfg, container)

	srv := &http.Server{
		Addr:    cfg.Port(),
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatalf("server listen: %v", err)
		}
	}()

	// 5) ждём сигнал остановки
	<-appCtx.Done()
	logger.Infof("shutdown signal caught")

	// даём 5 секунд на корректное завершение активных запросов
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Errorf("HTTP shutdown: %v", err)
	}

	container.Queue.Close()

	// закрываем соединения БД, если у вас есть db.Close()
	//if err := db.Close(); err != nil {
	//	logger.Errorf("DB close: %v", err)
	//}

	logger.Info("graceful shutdown complete")
}

func Init(cfg contracts.Config, container *app.Container) *gin.Engine {
	gob.Register(user.User{})

	router := gin.Default()
	err := router.SetTrustedProxies(nil)
	if err != nil {
		panic(err.Error())
	}

	store := models.Store(cfg.RedisHost(), "", cfg.KeyCookie())
	router.Use(sessions.Sessions(cfg.SessionsName(), store))

	db.Init(cfg.DBUrl())

	web.InitMinify()
	web.NewTemplate(router)
	routes.InitializeApiRoutes(router, container)
	routes.InitializeWebRoutes(router, container)
	return router
}
