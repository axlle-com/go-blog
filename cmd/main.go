package main

import (
	"encoding/gob"
	"errors"
	"github.com/axlle-com/blog/pkg/common/config"
	"github.com/axlle-com/blog/pkg/common/db"
	"github.com/axlle-com/blog/pkg/common/logger"
	"github.com/axlle-com/blog/pkg/common/models"
	"github.com/axlle-com/blog/pkg/common/models/contracts"
	"github.com/axlle-com/blog/pkg/common/routes"
	"github.com/axlle-com/blog/pkg/common/web"
	user "github.com/axlle-com/blog/pkg/user/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Config()
	router := Init(cfg)
	err := router.Run(cfg.Port())
	if err != nil {
		panic("Error run")
	}
}

func Init(cfg contracts.Config) *gin.Engine {
	gob.Register(user.User{})

	router := gin.Default()
	err := router.SetTrustedProxies(nil)
	if err != nil {
		panic(err.Error())
	}

	store := models.Store(cfg.RedisHost(), "", cfg.KeyCookie())
	router.Use(sessions.Sessions(cfg.SessionsName(), store))

	db.Init(cfg.DBUrl())

	//web.InitMinify()
	web.InitTemplate(router)
	routes.InitializeApiRoutes(router)
	routes.InitializeWebRoutes(router)
	logger.Debug(555)
	logger.Error(errors.New("55"))
	logger.Info(errors.New("55"))
	logger.Info(55)
	logger.Debug(55)
	return router
}
