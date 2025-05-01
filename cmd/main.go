package main

import (
	"context"
	"encoding/gob"
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
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
	appCtx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT,
	)
	defer stop()

	cfg := config.Config()
	db.InitDB(cfg)

	container := app.NewContainer(cfg, appCtx)
	router := Init(cfg, container)

	srv := &http.Server{
		Addr:    "0.0.0.0" + cfg.Port(),
		Handler: router,
	}

	go func() {
		logger.Infof("[Main] Listening on %s and on [::]%s", "0.0.0.0"+cfg.Port(), cfg.Port())
		if err := srv.ListenAndServe(); err != nil {
			logger.Debugf("[Main] server error: %v", err)
		}
	}()

	// ждём сигнал остановки
	<-appCtx.Done()
	logger.Info("[Main] Shutdown signal caught")

	// даём 5 секунд на корректное завершение активных запросов
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Errorf("[Main] HTTP shutdown: %v", err)
	}

	container.Queue.Close()

	//if err := db.Close(); err != nil {
	//	logger.Errorf("DB close: %v", err)
	//}

	logger.Info("[Main] Graceful shutdown complete")
}

func Init(config contracts.Config, container *app.Container) *gin.Engine {
	gob.Register(user.User{})

	router := gin.Default()
	err := router.SetTrustedProxies(nil)
	if err != nil {
		panic(err.Error())
	}

	store := models.Store(config)
	router.Use(sessions.Sessions(config.SessionsName(), store))
	router.Use(gzip.Gzip(gzip.BestSpeed))
	router.Use(csrf.Middleware(csrf.Options{
		Secret: string(config.KeyCookie()),
		ErrorFunc: func(c *gin.Context) {
			c.String(http.StatusForbidden, "CSRF token mismatch")
			c.Abort()
		},
	}))

	err = container.Migrator.Migrate()
	if err != nil {
		logger.Errorf("[Main][Init] Migrate error: %v", err)
	}

	err = container.Seeder.Seed()
	if err != nil {
		logger.Errorf("[Main][Init] Seed error: %v", err)
	}

	web.Minify(config)
	web.NewTemplate(router)
	routes.InitializeApiRoutes(router, container)
	routes.InitializeWebRoutes(router, container)
	return router
}
