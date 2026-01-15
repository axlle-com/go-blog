package main

import (
	"context"
	"encoding/gob"
	"errors"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/axlle-com/blog/app/config"
	"github.com/axlle-com/blog/app/db"
	"github.com/axlle-com/blog/app/di"
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contract"
	"github.com/axlle-com/blog/app/routes"
	user "github.com/axlle-com/blog/pkg/user/models"
	"github.com/gin-gonic/gin"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT,
	)
	defer stop()

	cfg := config.Config()

	newDB, err := db.SetupDB(cfg)
	if err != nil {
		panic("db not initialized")
	}

	container := di.NewContainer(cfg, newDB)
	router := Init(cfg, container)
	container.Queue.Start(ctx, 5)
	container.Scheduler.Start()

	srv := &http.Server{
		Addr:    "0.0.0.0" + cfg.Port(),
		Handler: router,
	}

	go func() {
		logger.Infof("[main][main] listening on %s and on [::]%s", "0.0.0.0"+cfg.Port(), cfg.Port())
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Errorf("[main][main] server error: %v", err)
		}
	}()

	// ждём сигнал остановки
	<-ctx.Done()
	logger.Info("[main][main] shutdown signal caught")

	// даём 5 секунд на корректное завершение активных запросов
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Errorf("[main][main] HTTP shutdown error: %v", err)
	}

	container.Queue.Close()
	container.Scheduler.Stop()

	if err = newDB.Close(); err != nil {
		logger.Errorf("[main][main] DB close error: %v", err)
	}

	logger.Info("[main][main] graceful shutdown complete")
}

func Init(config contract.Config, container *di.Container) *gin.Engine {
	gob.Register(user.User{})

	router := gin.Default()
	err := router.SetTrustedProxies(nil)
	if err != nil {
		panic(err.Error())
	}

	err = container.Migrator.Migrate()
	if err != nil {
		logger.Errorf("[main][Init][Migrator] migrate error: %v", err)
	}

	err = container.Seeder.Seed()
	if err != nil {
		logger.Errorf("[main][Init][Seeder] seed error: %v", err)
	}

	err = container.Minifier.Run()
	if err != nil {
		logger.Errorf("[main][Init][WebMinifier] error: %v", err)
	}

	err = container.Disk.SetupStaticFiles(router)
	if err != nil {
		logger.Errorf("[main][Init][SetupStaticFiles] error: %v", err)
	}

	container.View.SetRouter(router)
	container.View.Load()

	routes.InitApiRoutes(router, container)
	routes.InitWebRoutes(router, config, container)

	return router
}
