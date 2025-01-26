package main

import (
	"encoding/gob"
	"github.com/axlle-com/blog/pkg/app"
	"github.com/axlle-com/blog/pkg/app/config"
	"github.com/axlle-com/blog/pkg/app/db"
	"github.com/axlle-com/blog/pkg/app/models"
	"github.com/axlle-com/blog/pkg/app/models/contracts"
	"github.com/axlle-com/blog/pkg/app/routes"
	"github.com/axlle-com/blog/pkg/app/web"
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

	container := app.New()

	//web.InitMinify()
	web.InitTemplate(router)
	routes.InitializeApiRoutes(router, container)
	routes.InitializeWebRoutes(router, container)
	return router
}
